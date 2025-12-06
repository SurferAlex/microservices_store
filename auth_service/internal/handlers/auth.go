package handlers

import (
	"auth_service/internal/entity"
	"auth_service/internal/repository/psql"
	"auth_service/internal/security"
	"auth_service/internal/service"
	"auth_service/internal/tokens"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	var registerData struct {
		Username string `json:"username" validate:"required,min=3,max=20"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверный JSON",
		})
		return
	}

	if err := validate.Struct(registerData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверка сложности пароля
	if !isPasswordStrong(registerData.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пароль должен содержать минимум 8 символов, заглавную букву, цифру и спецсимвол",
		})
		return
	}

	// Проверка на существование пользователя
	existingUser, err := psql.GetUserByUsername(registerData.Username)
	if err == nil && existingUser != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь уже существует.",
		})
		return
	}

	// Проверка на существование email
	existingByEmail, err := psql.GetUserByEmail(registerData.Email)
	if err == nil && existingByEmail != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Email уже используется.",
		})
		return
	}

	// Добавление нового пользователя
	newUser := entity.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	// Хеш пароля
	hashed, err := security.HashPassword(newUser.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Ошибка хеширования пароля",
		})
		return
	}
	newUser.Password = hashed

	if err := psql.InsertUser(newUser); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Ошибка при сохранении пользователя.",
		})
		return
	}

	// Получение пользователя и назначение роли
	createdUser, err := psql.GetUserByUsername(newUser.Username)
	if err == nil && createdUser != nil {
		_ = psql.AssignRoleToUser(createdUser.ID, "user")
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	})
}

func isPasswordStrong(password string) bool {
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString
	)
	return len(password) >= 8 && hasUpper(password) && hasNumber(password) && hasSpecial(password)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Путь к HTML

	tmpl, err := template.ParseFiles("./frontend/templates/register.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text.html")

	// Выполнение шаблона (отправка в браузер)
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		return
	}

}
func Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверный формат запроса.",
		})
	}

	user, err := psql.GetUserByUsername(loginData.Username)
	if err != nil || user == nil || !security.CheckPasswordHash(loginData.Password, user.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверные данные для входа.",
		})
		return
	}

	// Создание JWT - токена
	token, err := tokens.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Ошибка создания токена: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Внутренняя ошибка сервера",
		})
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"token":    token,
		"redirect": "/profile",
	})

	// Сгенерировать и сохранить refresh
	refresh, err := tokens.GenerateRefreshOpaque(32)
	if err == nil {
		hash := tokens.HashRefreshToken(refresh)
		exp := time.Now().Add(tokens.RefreshTTL())
		ua := r.UserAgent()
		ip := r.RemoteAddr
		_ = psql.SaveRefreshToken(user.ID, hash, ua, ip, exp)
		service.SetRefreshCookie(w, refresh, tokens.RefreshTTL())
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./frontend/templates/login.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text.html")

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		return
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, "no refresh", http.StatusUnauthorized)
		return
	}
	access, newRefresh, err := service.RotateRefresh(c.Value, r.UserAgent(), r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid refresh", http.StatusUnauthorized)
		return
	}
	service.SetRefreshCookie(w, newRefresh, tokens.RefreshTTL())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": access,
	})
}

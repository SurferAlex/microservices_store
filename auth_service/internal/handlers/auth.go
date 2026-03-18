package handlers

import (
	"auth_service/internal/entity"
	"auth_service/internal/repository/psql"
	"auth_service/internal/security"
	"auth_service/internal/service"
	"auth_service/internal/tokens"
	"encoding/json"
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
		if err2 := json.NewEncoder(w).Encode(map[string]any{
			"message": "Неверный JSON",
			"error":   err.Error(),
		}); err2 != nil {
			log.Printf("не удалось отправить JSON-ответ об ошибке: %v", err2)
		}
		return
	}

	if err := validate.Struct(registerData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err2 := json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверный JSON",
			"error":   err.Error(),
		}); err2 != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err2)
		}
		return
	}

	// Проверка сложности пароля
	if !isPasswordStrong(registerData.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Пароль должен содержать минимум 8 символов, заглавную букву, цифру и спецсимвол",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
		return
	}

	// Проверка на существование пользователя
	existingUser, err := psql.GetUserByUsername(registerData.Username)
	if err == nil && existingUser != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь уже существует.",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
		return
	}

	// Проверка на существование email
	existingByEmail, err := psql.GetUserByEmail(registerData.Email)
	if err == nil && existingByEmail != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Email уже используется.",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
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
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Ошибка хеширования пароля",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
		return
	}
	newUser.Password = hashed

	if err := psql.InsertUser(newUser); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Ошибка при сохранении пользователя.",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
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
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	}); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}
}

func isPasswordStrong(password string) bool {
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString
	)
	return len(password) >= 8 && hasUpper(password) && hasNumber(password) && hasSpecial(password)
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
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверный формат запроса.",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
	}

	user, err := psql.GetUserByUsername(loginData.Username)
	if err != nil || user == nil || !security.CheckPasswordHash(loginData.Password, user.Password) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Неверные данные для входа.",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
		return
	}

	// Создание JWT - токена
	token, err := tokens.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Ошибка создания токена: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Внутренняя ошибка сервера",
		}); err != nil {
			log.Printf("не удалось отправить JSON-ответ: %v", err)
		}
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"status":   "success",
		"token":    token,
		"redirect": "http://localhost:8081/profile",
	}); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}

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

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Refresh token не найден",
	}); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}

	hash := tokens.HashRefreshToken(c.Value)

	// Отзываем токен
	_ = psql.RevokeRefreshToken(hash)

	// Удаляем cookie
	service.DeleteCookie(w, c.Value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err2 := json.NewEncoder(w).Encode(map[string]string{
		"message": "Выход выполнен",
	}); err2 != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
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
	if err := json.NewEncoder(w).Encode(map[string]string{
		"access_token": access,
	}); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}
}

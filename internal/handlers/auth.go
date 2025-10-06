package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"vscode_test/internal/entity"
	"vscode_test/internal/repository/psql"
	"vscode_test/internal/tokens"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var registerData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Проверка на существование пользователя
	existingUser, err := psql.GetUserByUsername(registerData.Username)
	if err == nil && existingUser != nil {
		http.Error(w, "Пользователь уже существует.", http.StatusConflict)
		return
	}

	// Добавление нового пользователя

	newUser := entity.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	if err := psql.InsertUser(newUser); err != nil {
		http.Error(w, "Ошибка при сохранении пользователя.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	})

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
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Неверный формат запроса.", http.StatusBadRequest)
		return
	}

	user, err := psql.GetUserByUsername(loginData.Username)
	if err != nil || user == nil || user.Password != loginData.Password {
		http.Error(w, "Неверные данные для входа.", http.StatusUnauthorized)
		return
	}

	// Создание JWT - токена
	token, err := tokens.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("Ошибка создания токена: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"token":    token,
		"redirect": "/profile",
	})
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

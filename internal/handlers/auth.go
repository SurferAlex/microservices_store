package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
)

var usersReg []Creds

type Creds struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {

	var newUser Creds

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, что такой пользователь не существует
	for _, user := range usersReg {
		if user.Username == newUser.Username || user.Email == newUser.Email {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}

	}

	usersReg = append(usersReg, newUser)

	response := map[string]interface{}{
		"username": newUser.Username,
		"email":    newUser.Email,
		"message":  "Пользователь успешно зарегистрирован",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Путь к HTML

	tmpl, err := template.ParseFiles("../frontend/templates/register.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Countent-type", "text.html")

	// Выполнение шаблона (отправка в браузер)
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Countent-type", "application/json")

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный JSON",
		})
		return
	}

	var foundUser *Creds
	for i := range usersReg {
		if usersReg[i].Username == loginData.Username && usersReg[i].Password == loginData.Password {
			foundUser = &usersReg[i]
			break
		}
	}

	if foundUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный логин или пароль",
		})
		return
	}

	response := map[string]interface{}{
		"message": "Успешный вход в систему",
		"type":    "Bearer",
		"user": map[string]interface{}{
			"username": foundUser.Username,
		},
	}

	json.NewEncoder(w).Encode(response)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("../frontend/templates/login.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Countent-type", "text.html")

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка выполнения шаблона", http.StatusInternalServerError)
		return
	}
}

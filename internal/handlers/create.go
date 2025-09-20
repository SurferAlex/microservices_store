package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

var users []User
var nextID = 1

func init() {
	users = []User{
		{ID: 1, Name: "Artur", Email: "artielucifer@gmail.com", Age: 27},
		{ID: 2, Name: "Tanya", Email: "tanya_ustinova_98@list.ru", Age: 26},
	}
	nextID = 3
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок
	w.Header().Set("Countetn-type", "application/json")

	//Создаем нового пользователя
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	//Присваиваем ID и добавляем в список
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	// Возвращаем созданного пользователя
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)

}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Countetn-type", "application/json")

	// Получение ID из URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Конвертируем в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Поиск пользователя
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "Пользователь не найден.", http.StatusNotFound)

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updateUser User
	err = json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if user.ID == id {
			updateUser.ID = id
			users[i] = updateUser
			json.NewEncoder(w).Encode(updateUser)
			return
		}
	}
	http.Error(w, "Пользователь не найден.", http.StatusNotFound)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contetnt-type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.Error(w, "Пользователь не найден.", http.StatusNotFound)
	}
}

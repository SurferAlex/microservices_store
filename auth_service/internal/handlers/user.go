package handlers

import (
	"auth_service/internal/repository/psql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Некорректный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := psql.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}
}

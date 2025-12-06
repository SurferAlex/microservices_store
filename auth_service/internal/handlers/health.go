package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок
	w.Header().Set("Countent-type", "application/json")

	// Создаем ответ
	response := HealthResponse{
		Status:  "OK",
		Message: "Сервер работает без перебоев",
	}

	json.NewEncoder(w).Encode(response)
}

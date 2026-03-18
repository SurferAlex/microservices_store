package handlers

import (
	"encoding/json"
	"log"
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("не удалось отправить JSON-ответ: %v", err)
	}
}

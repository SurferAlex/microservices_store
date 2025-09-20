package main

import (
	"fmt"
	"net/http"
	"vscode_test/internal/api"

	"github.com/gorilla/mux"
)

func main() {

	// Регистрация маршрутов
	r := mux.NewRouter()
	api.SetupRoutes(r)

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("🔍 Проверь: http://localhost:8080/health")
	fmt.Println("⏹️  Для остановки нажми Ctrl+C")

	http.ListenAndServe(":8080", r) // Просто запускаем
}

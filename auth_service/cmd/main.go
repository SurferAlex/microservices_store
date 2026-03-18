package main

import (
	"auth_service/internal/api"
	"auth_service/internal/config"
	"auth_service/internal/repository/psql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// Загрузка .env
	_ = godotenv.Load()

	// Загружаем конфиг с дефолтным значением
	cfg := config.LoadConfig()

	// Используем готовую строку подключения
	connectionString := cfg.GetDBConnectionString()

	// Подключение к БД
	_, err := psql.InitDB(connectionString)
	if err != nil {
		log.Fatalf("init DB: %v", err)
	}

	// Регистрация маршрутов
	r := mux.NewRouter()
	api.SetupRoutes(r)
	port := ":8080"

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")

	err = http.ListenAndServe(port, r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("ошибка запуска сервера: %v", err)
	} // Просто запускаем
}

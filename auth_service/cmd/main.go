package main

import (
	"auth_service/internal/api"
	"auth_service/internal/config"
	"auth_service/internal/repository/psql"
	"fmt"
	"log"
	"net/http"
	"os"

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
	psql.InitDB(connectionString)

	// Проверка миграций
	dbURL := os.Getenv("DB_URL")
	err := psql.RunMigrations(dbURL)
	if err != nil {
		log.Fatalf("Миграции не прошли: %v\n", err)
	}
	log.Println("Миграции успешно применены!")

	// Регистрация маршрутов
	r := mux.NewRouter()
	api.SetupRoutes(r)
	port := ":8080"

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("🔍 Проверь: http://localhost:8080/health")
	fmt.Println("⏹️  Для остановки нажми Ctrl+C")

	http.ListenAndServe(port, r) // Просто запускаем
}

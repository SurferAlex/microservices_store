package main

import (
	"fmt"
	"log"
	"net/http"
	"vscode_test/internal/api"
	"vscode_test/internal/config"
	"vscode_test/internal/repository/psql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	// Загружаем конфиг с дефолтным значением
	cfg := config.LoadConfig()

	// Используем готовую строку подключения
	connectionString := cfg.GetDBConnectionString()

	// Подключение к БД
	psql.InitDB(connectionString)

	// Создание таблиц в БД
	if err := psql.CreateUsersTable(); err != nil {
		log.Fatal("Ошибка создания таблицы users", err)
	}

	// Регистрация маршрутов
	r := mux.NewRouter()
	api.SetupRoutes(r)

	fmt.Println("🚀 Сервер запущен на http://localhost:8080")
	fmt.Println("🔍 Проверь: http://localhost:8080/health")
	fmt.Println("⏹️  Для остановки нажми Ctrl+C")

	http.ListenAndServe(":8080", r) // Просто запускаем
}

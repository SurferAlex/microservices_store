package main

import (
	"fmt"
	"log"
	"os"
	"profile_service/internal/api"
	"profile_service/internal/config"
	"profile_service/internal/repository/psql"
	"profile_service/internal/service"

	"github.com/gin-gonic/gin"
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

	//Создание HTTP клиента для auth_service
	authClient := service.NewAuthClient(cfg.AuthServiceURL)

	// Создание Gin роутера
	r := gin.Default()
	api.SetupRoutes(r, authClient)

	// Указываем доверенные прокси
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	fmt.Println("🚀 Сервер запущен на http://localhost:8081")

	// Запуск сервера
	r.Run(":8081")

}

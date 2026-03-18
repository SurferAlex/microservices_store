package main

import (
	"fmt"
	"log"
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
	_, err := psql.InitDB(connectionString)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}

	//Создание HTTP клиента для auth_service
	authClient := service.NewAuthClient(cfg.AuthServiceURL)

	// Создание Gin роутера
	r := gin.Default()
	api.SetupRoutes(r, authClient)

	// Указываем доверенные прокси
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Fatalf("set trusted proxies: %v", err)
	}

	fmt.Println("🚀 Сервер запущен на http://localhost:8081")

	// Запуск сервера
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Ошибка запуска севвера %v", err)
	}

}

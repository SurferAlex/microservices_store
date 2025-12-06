package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Продолжаем выполнение
		c.Next()

		// Логируем после выполнения
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - %v\n",
			start.Format("15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Устанавливаем заголовки CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обработка preflight запросы
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Генерируем ID запроса
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		// Сохраняем в контексте
		c.Set("request_id", requestID)

		// Добавляем в заголовок ответа
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

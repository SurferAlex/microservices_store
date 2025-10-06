package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"vscode_test/internal/tokens"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Запоминаем время начала
		start := time.Now()

		// Логируем начало запроса
		fmt.Printf("→ [%s] %s %s\n",
			start.Format("15:04:05"), // время в формате ЧЧ:ММ:СС
			r.Method,                 // GET, POST, etc.
			r.URL.Path,               // /users, /users/123
		)

		// Выполняем основной handler
		next.ServeHTTP(w, r)

		// Логируем завершение с временем выполнения
		duration := time.Since(start)
		fmt.Printf("← Выполнено за %v\n\n", duration)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков для CORS
		w.Header().Set("Acces-Control-Allow-Origin", "*")
		w.Header().Set("Acces-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Acces-Control-Allow-Headers", "Countent-type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsValidToken(token string) bool {
	tokenString := strings.TrimPrefix(token, "Bearer")
	_, err := tokens.ValidateJWT(tokenString)
	return err == nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	// Публичные ендпоинты не требующие авторизации

	publicPath := map[string]bool{
		"/health":   true,
		"/login":    true,
		"/register": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Пропускаем все статические файлы без авторизации
		if strings.HasPrefix(r.URL.Path, "/frontend/") {
			next.ServeHTTP(w, r)
			return
		}

		// Проверка публичных путей
		if publicPath[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Логика проверки токена
		token := r.Header.Get("Authorization")

		if token == "" {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"Ошибка": "Токен отсутвует."}`))
			return
		}

		if !IsValidToken(token) {
			w.Header().Set("Counter-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"Ошибка": "Неверный токен."}`))
			return
		}

		next.ServeHTTP(w, r)

	})

}

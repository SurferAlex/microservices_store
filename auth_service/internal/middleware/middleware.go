package middleware

import (
	"auth_service/internal/tokens"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("→ [%s] %s %s\n", start.Format("15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		fmt.Printf("← Выполнено за %v\n\n", duration)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

type ctxKey string

const ctxUserID ctxKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	publicPath := map[string]bool{
		"/health":       true,
		"/login":        true,
		"/register":     true,
		"/logout":       true,
		"/auth/refresh": true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/users/") {
			next.ServeHTTP(w, r)
			return
		}

		if publicPath[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(`{"error": "Токен отсутвует."}`)); err != nil {
				log.Printf("не удалось отправить JSON-ответ: %v", err)
			}
			return
		}
		tokenString := strings.TrimPrefix(token, "Bearer ")
		claims, err := tokens.ValidateJWT(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(`{"error": "Неверный токен."}`)); err != nil {
				log.Printf("не удалось отправить JSON-ответ: %v", err)
			}
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) (int, bool) {
	v := r.Context().Value(ctxUserID)
	id, ok := v.(int)
	return id, ok
}

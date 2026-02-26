package api

import (
	"auth_service/internal/handlers"
	"auth_service/internal/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	// Middleware
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	// Публичные endpoints
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/api/v1/users/{id}", handlers.GetUserByID).Methods("GET")

	r.Handle("/register", middleware.RateLimitIP(3, time.Minute)(http.HandlerFunc(handlers.Register))).Methods("POST")
	r.Handle("/login", middleware.RateLimitIP(5, time.Minute)(http.HandlerFunc(handlers.Login))).Methods("POST")
	r.HandleFunc("/auth/refresh", handlers.Refresh).Methods("POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")

	// Защищенные endpoints
	r.HandleFunc("/admin/users/{id}/roles", handlers.GetUserRoles).Methods("GET")
	r.HandleFunc("/admin/users/{id}/roles", handlers.AssignRole).Methods("PUT")

}

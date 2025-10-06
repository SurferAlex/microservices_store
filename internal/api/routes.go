package api

import (
	"net/http"
	"vscode_test/internal/handlers"
	"vscode_test/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {

	// Middleware
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	// Публичные endpoints
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Защищенные endpoints
	r.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	// Статические файлы

	r.PathPrefix("/frontend/").Handler(http.StripPrefix("/frontend/", http.FileServer(http.Dir("./frontend"))))

}

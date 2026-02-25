package api

import (
	"net/http"
	"profile_service/internal/handlers"
	"profile_service/internal/middleware"
	"profile_service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authClient *service.AuthClient) {
	// Глобальный middleware
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoginMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Публичные маршруты
	r.GET("/health", handlers.HealthHandler)

	// HTML шаблоны
	r.LoadHTMLGlob("frontend/templates/*")
	r.Static("/frontend", "./frontend")

	r.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	r.GET("/profile/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create-profile.html", nil)
	})

	api := r.Group("/api/v1")
	{
		// Профили
		profiles := api.Group("/profiles")
		{
			profiles.GET("/:id", handlers.GetProfile)
			profiles.GET("/user/:user_id", handlers.GetProfileByUserID)
			profiles.POST("", handlers.CreateProfile(authClient))
			profiles.PUT("/:id", handlers.UpdateProfile)
			profiles.DELETE("/:id", handlers.DeleteProfile)
		}

		// Адреса (отдельная группа, без конфликта)
		addresses := api.Group("/addresses")
		{
			addresses.GET("/profile/:profile_id", handlers.GetAddresses)
			addresses.POST("/profile/:profile_id", handlers.CreateAddress)
			addresses.PUT("/:id", handlers.UpdateAddress)
			addresses.DELETE("/:id", handlers.DeleteAddress)
			addresses.PUT("/profile/:profile_id/:address_id/primary", handlers.SetPrimaryAddress)
		}

		// Контакты (отдельная группа)
		contacts := api.Group("/contacts")
		{
			contacts.GET("/profile/:profile_id", handlers.GetContacts)
			contacts.POST("/profile/:profile_id", handlers.CreateContact)
			contacts.PUT("/:id", handlers.UpdateContact)
			contacts.DELETE("/:id", handlers.DeleteContact)
			contacts.PUT("/:id/verify", handlers.VerifyContact)
		}
	}
}

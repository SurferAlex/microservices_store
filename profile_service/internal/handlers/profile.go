package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"profile_service/internal/entity"
	"profile_service/internal/repository/psql"
	"profile_service/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func GetProfile(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	profile, err := psql.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при получении профиля",
			"details": err.Error(),
		})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Профиль найден",
		"data":       profile,
		"request_id": requestID,
	})
}

func GetProfileByUserID(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный user_id"})
		return
	}

	profile, err := psql.GetProfileByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Профиль найден",
		"data":       profile,
		"request_id": requestID,
	})
}

func CreateProfile(authClient *service.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("request_id")

		// Читаем сырое тело запроса для отладки
		bodyBytes, _ := c.GetRawData()
		fmt.Printf("🔍 Получен запрос: %s\n", string(bodyBytes))

		// Восстанавливаем тело для биндинга
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var payload entity.Profile
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверный формат данных",
			})
			return
		}

		fmt.Printf("✅ Данные распарсены: %+v\n", payload)

		if payload.UserID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id обязателен"})
			return
		}

		exists, err := authClient.CheckUserExists(payload.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Не удалось проверить пользователя в auth_service",
			})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Пользователь не найден в auth_service",
			})
			return
		}

		created, err := psql.CreateProfile(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать профиль"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "Профиль создан",
			"data":       created,
			"request_id": requestID,
		})
	}
}

func UpdateProfile(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	profileID, err := strconv.Atoi(idParam)
	if err != nil || profileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	// Получаем существующий профиль для user_id
	existingProfile, err := psql.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при получении профиля",
			"details": err.Error(),
		})
		return
	}

	// Проверяем, что профиль существует
	if existingProfile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль не найден"})
		return
	}

	var payload entity.Profile
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	// Устанавливаем ID и UserID из существующего профиля
	payload.ID = profileID
	payload.UserID = existingProfile.UserID // Берем из БД, не из запроса

	if err := psql.UpdateProfile(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить профиль"})
		return
	}

	updated, err := psql.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Профиль обновлён, но не удалось получить данные"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Профиль обновлён",
		"data":       updated,
		"request_id": requestID,
	})
}

func DeleteProfile(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	if err := psql.DeleteProfile(id); err != nil {
		// Проверяем на ошибки
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить профиль"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Профиль удален",
		"request_id": requestID,
	})
}

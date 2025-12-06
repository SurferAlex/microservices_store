package handlers

import (
	"net/http"
	"profile_service/internal/entity"
	"profile_service/internal/repository/psql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateContact(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	profileIDParam := c.Param("profile_id")
	profileID, err := strconv.Atoi(profileIDParam)
	if err != nil || profileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный profile_id"})
		return
	}

	// Проверяем, что профиль существует
	existingProfile, err := psql.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при получении профиля",
			"details": err.Error(),
		})
		return
	}
	if existingProfile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль не найден"})
		return
	}

	var payload entity.ProfileContact
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	payload.ProfileID = profileID

	created, err := psql.CreateContact(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Не удалось создать контакт",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Контакт создан",
		"data":       created,
		"request_id": requestID,
	})
}

func GetContacts(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	profileIDParam := c.Param("profile_id")
	profileID, err := strconv.Atoi(profileIDParam)
	if err != nil || profileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный profile_id"})
		return
	}

	contacts, err := psql.GetContactsByProfileID(profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Контакты не найдены"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Контакты получены",
		"data":       contacts,
		"request_id": requestID,
	})
}

func UpdateContact(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	contactID, err := strconv.Atoi(idParam)
	if err != nil || contactID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID контакта"})
		return
	}

	var payload entity.ProfileContact
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	payload.ID = contactID

	if err := psql.UpdateContact(payload); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить контакт"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Контакт обновлён",
		"data":       payload,
		"request_id": requestID,
	})
}

func DeleteContact(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	contactID, err := strconv.Atoi(idParam)
	if err != nil || contactID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID контакта"})
		return
	}

	if err := psql.DeleteContact(contactID); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить контакт"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Контакт удалён",
		"request_id": requestID,
	})
}

func VerifyContact(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	contactID, err := strconv.Atoi(idParam)
	if err != nil || contactID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID контакта"})
		return
	}

	if err := psql.VerifyContact(contactID); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось подтвердить контакт"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Контакт подтверждён",
		"request_id": requestID,
	})
}

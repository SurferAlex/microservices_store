package handlers

import (
	"net/http"
	"profile_service/internal/entity"
	"profile_service/internal/repository/psql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateAddress(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	// Получаем profile_id из URL
	profileIDParam := c.Param("profile_id")
	profileID, err := strconv.Atoi(profileIDParam)
	if err != nil || profileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный profile_id"})
		return
	}

	// Проверяем, что профиль существует
	_, err = psql.GetProfileByID(profileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Профиль не найден"})
		return
	}

	var payload entity.ProfileAddress
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неправильный формат данных",
			"details": err.Error(),
		})
		return
	}

	// Используем profile_id из URL
	payload.ProfileID = profileID

	created, err := psql.CreateAddress(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Не удалось создать адрес",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Адрес создан",
		"data":       created,
		"request_id": requestID,
	})
}

func GetAddresses(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("profile_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный profile_id."})
		return
	}

	addresses, err := psql.GetAddressesByProfileID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Адреса не найдены"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Адреса найдены",
		"data":       addresses,
		"request_id": requestID,
	})
}

func UpdateAddress(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	addressID, err := strconv.Atoi(idParam)
	if err != nil || addressID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	// Получаем существующий адрес для profile_id
	existingAddress, err := psql.GetAddressByID(addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при получении адреса",
			"details": err.Error(),
		})
		return
	}

	if existingAddress == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Адрес не найден"})
		return
	}

	var payload entity.ProfileAddress
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	// Устанавливаем ID и ProfileID из существующего адреса
	payload.ID = addressID
	payload.ProfileID = existingAddress.ProfileID

	if err := psql.UpdateAddress(payload); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить адрес"})
		}
		return
	}

	// Получаем обновленный адрес напрямую
	updated, err := psql.GetAddressByID(addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Адрес обновлен, но не удалось получить данные"})
		return
	}

	if updated == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Адрес обновлен, но не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Адрес обновлен",
		"data":       updated,
		"request_id": requestID,
	})
}

func DeleteAddress(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	if err := psql.DeleteAddress(id); err != nil {
		if strings.Contains(err.Error(), "не найден") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить адрес"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Адрес удален",
		"request_id": requestID,
	})
}

// Установка основного адреса
// Установка основного адреса
func SetPrimaryAddress(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	idParam1 := c.Param("profile_id")
	profileID, err := strconv.Atoi(idParam1)
	if err != nil || profileID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный профиль ID"})
		return
	}

	idParam2 := c.Param("address_id")
	addressID, err := strconv.Atoi(idParam2)
	if err != nil || addressID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный адрес ID"})
		return
	}

	if err := psql.SetPrimaryAddress(profileID, addressID); err != nil {
		if strings.Contains(err.Error(), "не найден") || strings.Contains(err.Error(), "не принадлежит") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось установить основной адрес"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Основной адрес обновлён",
		"request_id": requestID,
	})
}

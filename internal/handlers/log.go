package handlers

import (
	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
)

func GetLogs(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	var logs []models.DrinksLogs
	if err := database.DB.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}

	var result []dto.DrinkLogWithDetails
	for _, log := range logs {
		var drinks []models.DrinksDetails
		if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
		}
		result = append(result, dto.DrinkLogWithDetails{
			DrinksLogs: log,
			Drinks:     drinks,
		})
	}

	return c.JSON(http.StatusOK, result)
}

func GetLogByID(c echo.Context) error {
	logID := c.Param("id")
	logUUID, err := uuid.Parse(logID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid log id"})
	}
	var log models.DrinksLogs
	if err := database.DB.Where("id = ?", logUUID).First(&log).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "log not found"})
	}

	var drinks []models.DrinksDetails
	if err := database.DB.Where("drink_log_id = ?", log.ID).Find(&drinks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
	}

	return c.JSON(http.StatusOK, dto.DrinkLogWithDetails{
		DrinksLogs: log,
		Drinks:     drinks,
	})
}

func CreateLog(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	var req dto.LogRequest
	if err := c.Bind(&req); err != nil {
		log.Println("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	log := models.DrinksLogs{
		ID:      uuid.New(),
		UserID:  userID,
		Comment: req.Comment,
		// ImageURL:  req.ImageURL,
		DrinkDate: req.DrinkDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&log).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create log"})
	}

	var drinks []models.DrinksDetails
	for i := range req.Drinks {
		req.Drinks[i].ID = uuid.New()
		req.Drinks[i].DrinkLogID = log.ID
		if err := database.DB.Create(&req.Drinks[i]).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create drink detail"})
		}
		drinks = append(drinks, req.Drinks[i])
	}

	return c.JSON(http.StatusCreated, dto.DrinkLogWithDetails{
		DrinksLogs: log,
		Drinks:     drinks,
	})
}

func UpdateLog(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	logID := c.Param("id")
	logUUID, err := uuid.Parse(logID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid log id"})
	}

	var existing models.DrinksLogs
	if err := database.DB.Where("id = ? AND user_id = ?", logUUID, userID).First(&existing).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "log not found"})
	}

	var req dto.LogRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	// Drinkは既存を削除してから新規作成して入れ替える
	if err := database.DB.Where("drink_log_id = ?", logUUID).Delete(&models.DrinksDetails{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update drinks"})
	}

	for i := range req.Drinks {
		req.Drinks[i].ID = uuid.New()
		req.Drinks[i].DrinkLogID = logUUID
		if err := database.DB.Create(&req.Drinks[i]).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create drink detail"})
		}
	}

	existing.Comment = req.Comment
	// existing.ImageURL = req.ImageURL
	existing.DrinkDate = req.DrinkDate
	existing.UpdatedAt = time.Now()

	if err := database.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update log"})
	}

	// 追加: DrinksDetailsを取得して返す
	var drinks []models.DrinksDetails
	if err := database.DB.Where("drink_log_id = ?", logUUID).Find(&drinks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch drink details"})
	}

	return c.JSON(http.StatusOK, dto.DrinkLogWithDetails{
		DrinksLogs: existing,
		Drinks:     drinks,
	})
}

func DeleteLog(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	logID := c.Param("id")
	logUUID, err := uuid.Parse(logID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid log id"})
	}

	// あとで constraint:OnDelete:CASCADE;
	if err := database.DB.Where("drink_log_id = ?", logUUID).Delete(&models.DrinksDetails{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete drinks"})
	}

	if err := database.DB.Where("id = ? AND user_id = ?", logUUID, userID).Delete(&models.DrinksLogs{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete log"})
	}

	return c.NoContent(http.StatusNoContent)
}

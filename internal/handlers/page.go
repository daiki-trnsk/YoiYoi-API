package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/daiki-trnsk/YoiYoi-API/internal/usecase"
)

func GetHome(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	resp, err := usecase.GetHomeInfo(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func GetTimeline(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	resp, err := usecase.GetTimelineInfo(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func GetWeeklyStats(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	resp, err := usecase.GetPeriodStatsInfo(userID, 7)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func GetMonthlyStats(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	resp, err := usecase.GetPeriodStatsInfo(userID, 30)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

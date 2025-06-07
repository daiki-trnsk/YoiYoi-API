package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
	"github.com/daiki-trnsk/YoiYoi-API/internal/repositories"
	"github.com/google/uuid"
)

func GetUseByID(c echo.Context) error {
	id := c.Param("id")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	user, err := repositories.GetUserByID(uuidID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, dto.ToUserResponse(*user))
}
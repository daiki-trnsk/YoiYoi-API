package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
	"github.com/daiki-trnsk/YoiYoi-API/internal/repositories"
)

func GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := repositories.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, dto.ToUserResponse(*user))
}
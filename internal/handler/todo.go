package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	// "github.com/daiki-trnsk/YoiYoi-API/internal/repository"
	"github.com/daiki-trnsk/YoiYoi-API/internal/model"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
)

func GetTodos(c echo.Context) error {
	var todos []model.Todo
	database.DB.Find(&todos)
	return c.JSON(http.StatusOK, todos)
}

func CreateTodo(c echo.Context) error {
	var todo model.Todo
	if err := c.Bind(&todo); err != nil {
		log.Printf("Bind error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}
	database.DB.Create(&todo)
	return c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo model.Todo
	if err := database.DB.First(&todo, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}
	database.DB.Save(&todo)
	return c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.Todo{}, id)
	return c.NoContent(http.StatusNoContent)
}

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	
	"github.com/daiki-trnsk/YoiYoi-API/internal/handler"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
)

func main() {
	_ = godotenv.Load()

	database.Init()

	e := echo.New()

	// 仮でルートだけ
	e.GET("/", handler.Hello)

	e.GET("/todos", handler.GetTodos)
	e.POST("/todos", handler.CreateTodo)
	e.PUT("/todos/:id", handler.UpdateTodo)
	e.DELETE("/todos/:id", handler.DeleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/daiki-trnsk/YoiYoi-API/internal/handler"
)

func main() {
	_ = godotenv.Load()

	e := echo.New()

	// 仮でルートだけ
	e.GET("/", handler.Hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}

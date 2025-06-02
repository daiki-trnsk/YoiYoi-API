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

	conn, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer conn.Close(nil)
	log.Println("Connected to the database successfully")

	e := echo.New()

	// 仮でルートだけ
	e.GET("/", handler.Hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}

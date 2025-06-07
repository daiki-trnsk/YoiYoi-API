package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/daiki-trnsk/YoiYoi-API/internal/handlers"
	customMiddleware "github.com/daiki-trnsk/YoiYoi-API/internal/middleware"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
)

func main() {
	_ = godotenv.Load()

	database.Init()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.HEAD, echo.PATCH},
	}))

	e.Any("/", func(c echo.Context) error {
		return c.String(200, "YoiYoi-API is running")
	})

	e.POST("/auth/login", handlers.Login)
	e.POST("/auth/signup", handlers.Signup)

	auth := e.Group("")
	auth.Use(customMiddleware.Authentication)

	auth.GET("/auth/me", handlers.GetMe)
	// PUTの方がいいかも
	auth.PATCH("/auth/me", handlers.UpdateMe)
	auth.POST("/auth/logout", handlers.Logout)

	auth.GET("users/:id", handlers.GetUserByID)

	auth.GET("/logs/me", handlers.GetLogs)
	auth.GET("/logs/:id", handlers.GetLogByID)
	auth.POST("/logs", handlers.CreateLog)
	auth.PUT("/logs/:id", handlers.UpdateLog)
	auth.DELETE("/logs/:id", handlers.DeleteLog)

	auth.GET("/friends", handlers.GetFriends)
	// auth.GET("/friends/:id", handlers.GetFriendByID)
	auth.POST("/friends/request/:id", handlers.SendFriendRequest)
	auth.PATCH("/friends/accept/:id", handlers.AcceptFriendRequest)
	auth.DELETE("/friends/:id", handlers.DeleteFriendShip)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}

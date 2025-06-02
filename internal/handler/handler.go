package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	fmt.Println("Hello handler called")
	return c.String(http.StatusOK, "Hello, local API!")
}

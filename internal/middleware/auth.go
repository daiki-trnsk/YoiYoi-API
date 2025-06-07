package middleware

import (
	"net/http"

	token "github.com/daiki-trnsk/YoiYoi-API/internal/tokens"
	"github.com/labstack/echo/v4"
)

// 認証が絶対必要な処理はトークンがない場合エラー返す
func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientToken := c.Request().Header.Get("Authorization")
		if clientToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No Authorization Header Provided"})
		}
		claims, err := token.ValidateToken(clientToken, "access")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		c.Set("user_id", claims.UserID)
		return next(c)
	}
}

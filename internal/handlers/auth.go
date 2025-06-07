package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/internal/tokens"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/dto"
)

type SignupRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

type LoginRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type AuthResponse struct {
	User        dto.UserResponse `json:"user"`
	AccessToken string           `json:"access_token"`
}

type UpdatedAtRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	AvatarImg      string    `json:"avatar_img"`
	Bio            string    `json:"bio"`
	FavoriteDrinks string    `json:"favorite_drinks"`
	Motto          string    `json:"motto"`
}

// サインアップ
func Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// パスワードハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to hash password"})
	}

	user := models.Users{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hashed),
		Email:        req.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// DB保存
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user already exists or invalid"})
	}

	// トークン生成
	accessToken, refreshToken, err := tokens.GenerateTokens(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate token"})
	}

	user.RefreshToken = refreshToken
	database.DB.Save(&user)

	return c.JSON(http.StatusOK, AuthResponse{
		User:        dto.ToUserResponse(user),
		AccessToken: accessToken,
	})
}

// ログイン
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	var user models.Users
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user not found"})
	}

	// パスワード照合
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid password"})
	}

	// トークン生成
	accessToken, refreshToken, err := tokens.GenerateTokens(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate token"})
	}

	user.RefreshToken = refreshToken
	user.UpdatedAt = time.Now()
	database.DB.Save(&user)

	return c.JSON(http.StatusOK, AuthResponse{
		User:        dto.ToUserResponse(user),
		AccessToken: accessToken,
	})
}

// ユーザー情報取得
func GetMe(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	fmt.Println("userIDRaw:", userIDRaw)
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}
	fmt.Println("userID:", userID)

	var user models.Users
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// ユーザー情報更新
func UpdateMe(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	var req UpdatedAtRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	var user models.Users
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user not found"})
	}

	user.Username = req.Username
	user.Email = req.Email
	user.AvatarImg = req.AvatarImg
	user.Bio = req.Bio
	user.FavoriteDrinks = req.FavoriteDrinks
	user.Motto = req.Motto
	user.UpdatedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update user"})
	}

	return c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// ログアウト
func Logout(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	var user models.Users
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user not found"})
	}

	// リフレッシュトークンを空にしてログアウト
	user.RefreshToken = ""
	if err := database.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to logout"})
	}

	return c.NoContent(http.StatusNoContent)
}

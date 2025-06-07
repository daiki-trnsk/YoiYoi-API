package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/constants"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// 友達と友達申請を取得
func GetFriends(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	var friendships []models.Friends
	if err := database.DB.Where("follower_id = ? OR followee_id = ?", userID, userID).Find(&friendships).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch friendships"})
	}

	// 友達と友達申請の情報を分けた状態で返す

	return c.JSON(http.StatusOK, friendships)
}

// 友達リクエスト送信
func SendFriendRequest(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	friendID := c.Param("id")
	friendUUID, err := uuid.Parse(friendID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid friend id"})
	}

	var friendship models.Friends
	friendship.ID = uuid.New()
	friendship.FollowerID = userID
	friendship.FolloweeID = friendUUID
	friendship.Status = constants.FriendRequestPending
	friendship.CreatedAt = time.Now()

	if err := database.DB.Create(&friendship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, friendship)
}

// 友達リクエスト承認
func AcceptFriendRequest(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	friendshipID := c.Param("id")
	var friendship models.Friends
	if err := database.DB.Where("id = ? AND followee_id = ?", friendshipID, userID).First(&friendship).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "friend request not found"})
	}

	friendship.Status = constants.FriendRequestAccepted
	friendship.CreatedAt = time.Now()

	if err := database.DB.Save(&friendship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, friendship)
}

// リクエストの拒否と友達関係の削除
func DeleteFriendShip(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "user id not found"})
	}

	friendshipID := c.Param("id")
	var friendship models.Friends
	if err := database.DB.Where("id = ? AND (follower_id = ? OR followee_id = ?)", friendshipID, userID, userID).First(&friendship).Error; err != nil {
		fmt.Println("friendshipID", friendshipID, "userID", userID, "err:", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "friendship not found"})
	}

	if err := database.DB.Delete(&friendship).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

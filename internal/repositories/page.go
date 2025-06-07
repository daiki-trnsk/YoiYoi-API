package repositories

import (
	"time"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"github.com/google/uuid"
)

func GetUserByID(userID uuid.UUID) (*models.Users, error) {
	var user models.Users
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetFriends(userID uuid.UUID) ([]models.Friends, error) {
	var friends []models.Friends
	err := database.DB.
		Where("(follower_id = ? OR followee_id = ?) AND status = ?", userID, userID, "accepted").
		Find(&friends).Error
	return friends, err
}

func GetUsersByIDs(ids []uuid.UUID) ([]models.Users, error) {
	var users []models.Users
	err := database.DB.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func GetDrinkLogs(userID uuid.UUID, fromDate *time.Time, limit int) ([]models.DrinksLogs, error) {
	var logs []models.DrinksLogs
	query := database.DB.Where("user_id = ?", userID)
	if fromDate != nil {
		query = query.Where("drink_date >= ?", *fromDate)
	}
	if limit > 0 {
		query = query.Order("drink_date desc").Limit(limit)
	} else {
		query = query.Order("drink_date desc")
	}
	err := query.Find(&logs).Error
	return logs, err
}

func GetDrinkDetails(logID uuid.UUID) ([]models.DrinksDetails, error) {
	var drinks []models.DrinksDetails
	err := database.DB.Where("drink_log_id = ?", logID).Find(&drinks).Error
	return drinks, err
}

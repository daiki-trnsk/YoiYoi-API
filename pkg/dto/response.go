package dto

import (
	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/google/uuid"
	"time"
)

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	// IconURL         string    `json:"icon_url"`
	// Introduction    string    `json:"introduction"`
	// FavoriteAlcohol string    `json:"favorite_alcohol"`
	// DrinkingHistory string    `json:"drinking_history"`
	// HowDrinking     string    `json:"how_drinking"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DrinkLogWithDetails struct {
	models.DrinksLogs
	Drinks []models.DrinksDetails `json:"drinks"`
}

func ToUserResponse(u models.Users) UserResponse {
	return UserResponse{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		// IconURL:         u.IconURL,
		// Introduction:    u.Introduction,
		// FavoriteAlcohol: u.FavoriteAlcohol,
		// DrinkingHistory: u.DrinkingHistory,
		// HowDrinking:     u.HowDrinking,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

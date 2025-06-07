package dto

import (
	"time"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	// IconURL         string    `json:"icon_url"`
	// Introduction    string    `json:"introduction"`
	// FavoriteAlcohol string    `json:"favorite_alcohol"`
	// DrinkingHistory string    `json:"drinking_history"`
	// HowDrinking     string    `json:"how_drinking"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DrinkLogWithDetails struct {
	models.DrinksLogs
	Drinks []models.DrinksDetails `json:"drinks"`
}

func ToUserResponse(u models.Users) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		// IconURL:         u.IconURL,
		// Introduction:    u.Introduction,
		// FavoriteAlcohol: u.FavoriteAlcohol,
		// DrinkingHistory: u.DrinkingHistory,
		// HowDrinking:     u.HowDrinking,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type WeeklyStats struct {
	TotalAlcoholMl      int            `json:"total_alcohol_ml"`
	TotalNumberOfDrinks int            `json:"total_number_of_drinks"`
	AlcoholByWeekday    map[string]int `json:"alcohol_by_weekday"`
}

type HomeResponse struct {
	UserInfo    UserResponse          `json:"user_info"`
	WeeklyStats WeeklyStats           `json:"weekly_stats"`
	RecentLogs  []DrinkLogWithDetails `json:"recent_logs"`
}

type Timeline struct {
	User     UserResponse        `json:"user"`
	DrinkLog DrinkLogWithDetails `json:"drink_log"`
}

type TimelineResponse struct {
	FriendList []UserResponse `json:"friend_list"`
	Timeline   []Timeline     `json:"timeline"`
}

type PeriodStatsResponse struct {
	TotalAlcoholGram   int            `json:"total_alcohol_gram"`
	AverageAlcoholGram int            `json:"average_alcohol_gram"`
	PeriodDays         int            `json:"period_days"`
	ActualDrinkDays    int            `json:"actual_drink_days"`
	AlcoholByWeekday   map[string]int `json:"alcohol_by_weekday"`
}

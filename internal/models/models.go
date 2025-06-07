package models

import (
	"time"

	"github.com/google/uuid"
)

// ユーザー情報
type Users struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username       string    `gorm:"unique;not null" json:"username"`
	PasswordHash   string    `gorm:"not null;column:password_hash" json:"password_hash"`
	Email          string    `gorm:"unique;not null" json:"email"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	RefreshToken   string    `gorm:"unique;not null" json:"refresh_token"`
	AvatarImg      string    `json:"avatar_img"`
	Bio            string    `json:"bio"`
	FavoriteDrinks string    `json:"favorite_drinks"`
	Motto          string    `json:"motto"`
}

// 飲酒記録
type DrinksLogs struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	// Mood      string         `gorm:"not null" json:"mood"`
	// Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	Comment string `gorm:"type:text" json:"comment"`
	// ImageURL string `gorm:"type:text" json:"image_url"`
	// IsShared  bool           `gorm:"default:false" json:"is_shared"`
	DrinkDate time.Time `gorm:"not null" json:"drink_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 飲酒記録に登録するお酒の種類と量(DrinkLogsとDrinkDetailsは1対多)
type DrinksDetails struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	DrinkLogID uuid.UUID `gorm:"type:uuid;not null" json:"drink_log_id"`
	Name       string    `gorm:"not null" json:"name"`
	AmountMl   float64   `gorm:"not null" json:"amount_ml"`
	Abv        float64   `gorm:"not null" json:"abv"`
}

// 友達関係
type Friends struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FollowerID uuid.UUID `gorm:"type:uuid;not null" json:"follower_id"`
	FolloweeID uuid.UUID `gorm:"type:uuid;not null" json:"followee_id"`
	Status     string    `gorm:"not null" json:"status"`           // "pending", "accepted"
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"` // 申請日時
	// UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 承認日時
}

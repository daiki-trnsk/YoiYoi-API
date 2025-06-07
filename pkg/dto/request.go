package dto

import (
	"time"

	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
)

type LogRequest struct {
	Comment   string                `json:"comment"`
	// ImageURL  string                `json:"image_url"`
	DrinkDate time.Time             `json:"drink_date"`
	Drinks    []models.DrinksDetails `json:"drinks"`
}

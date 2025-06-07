package repositories

import (
	"github.com/daiki-trnsk/YoiYoi-API/internal/models"
	"github.com/daiki-trnsk/YoiYoi-API/pkg/database"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

func GetUseByID(id string) (*models.Users, error) {
    var user models.Users
    if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
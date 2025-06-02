package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/daiki-trnsk/YoiYoi-API/internal/model"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	for i := 0; i < 10; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = db
			break
		}
		log.Println("retrying database connection...", err)
		time.Sleep(2 * time.Second)
	}
	if DB == nil {
		log.Fatal("failed to connect to database")
	}
	log.Println("connected to database successfully")

	DB.AutoMigrate(&model.Todo{})
}

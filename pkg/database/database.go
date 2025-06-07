package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "github.com/daiki-trnsk/YoiYoi-API/internal/models"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	pgCfg := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // これでプリペアドステートメントを使わない
	}

	for i := 0; i < 10; i++ {
		db, err := gorm.Open(postgres.New(pgCfg), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			// PrepareStmt: false, // 不要
		})
		if err == nil {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(1)
			sqlDB.SetMaxOpenConns(3)
			sqlDB.SetConnMaxLifetime(time.Hour)
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
}

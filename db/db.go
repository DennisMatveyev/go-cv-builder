package db

import (
	"cv-builder/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("cv-builder.db"))
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	db.AutoMigrate(models.User{}, models.Profile{}, models.Job{})

	return db
}

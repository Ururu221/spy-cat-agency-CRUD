package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"project1/models"
)

var DB *gorm.DB

func InitDB(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	DB = db

	err = DB.AutoMigrate(&models.Cat{}, &models.Mission{}, &models.Target{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}

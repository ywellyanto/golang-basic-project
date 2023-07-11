package config

import (
	"log"

	"golang_basic_project/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/golang_basic_project?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("error: %v", err)
		panic("failed to connect database")
	}

	DB.AutoMigrate(&models.User{}, &models.UserDetail{}, &models.Order{}, &models.Book{}, &models.OrderBook{})
}

package models

import "gorm.io/gorm"

type UserDetail struct {
	gorm.Model
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  uint   `json:"user_id"`
}

type ResponseUserDetail struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       uint   `json:"role"`
	UserDetail UserDetail
	Orders     []Order
}

type RequestRegister struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     uint   `json:"role" binding:"required"`
}

type TokenRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RequestUser struct {
	Username string `json:"username" binding:"required"`
	Role     uint   `json:"role" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type ResponseUser struct {
	Username   string `json:"name"`
	Email      string `json:"email"`
	Role       uint   `json:"role"`
	UserDetail ResponseUserDetail
	Orders     []ResponseOrderSingle
}

type ResponseUserSingle struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Role     uint   `json:"role"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(reqPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqPass))
	if err != nil {
		return err
	}
	return nil
}

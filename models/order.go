package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderCode   string      `json:"order_code"`
	Description string      `json:"description"`
	UserID      uint        `json:"user_id"`
	User        User        `json:"user"`
	Books       []OrderBook `json:"books"`
}

type RequestOrder struct {
	OrderCode   string `json:"order_code" binding:"required"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id" binding:"required"`
}

type ResponseOrderSingle struct {
	ID          uint      `json:"id"`
	OrderCode   string    `json:"order_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"order_date"`
}

type ResponseOrder struct {
	ID          uint      `json:"id"`
	OrderCode   string    `json:"order_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"order_date"`
	User        ResponseUserSingle
}

type OrderBook struct {
	gorm.Model
	OrderID uint `json:"order_id"`
	Order   Order
	BookID  uint `json:"book_id"`
	Book    Book
	Price   uint `json:"price"`
}

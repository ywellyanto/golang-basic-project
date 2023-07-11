package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name        string      `json:"name"`
	Code        string      `json:"code"`
	Description string      `json:"description"`
	Orders      []OrderBook `json:"orders"`
}

type RequestBook struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

type ResponseBook struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

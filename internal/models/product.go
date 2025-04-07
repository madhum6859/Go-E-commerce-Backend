package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product in the catalog
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null" json:"stock"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category"`
	Images      []Image   `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Image represents a product image
type Image struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"not null" json:"url"`
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
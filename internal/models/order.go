package models

import (
	"time"
)

// Order represents a customer order
type Order struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `json:"user_id"`
	User          User           `json:"user"`
	OrderItems    []OrderItem    `json:"order_items"`
	TotalAmount   float64        `json:"total_amount"`
	Status        string         `gorm:"default:pending" json:"status"` // pending, paid, shipped, delivered, cancelled
	PaymentID     string         `json:"payment_id"`
	ShippingInfo  ShippingInfo   `json:"shipping_info"`
	ShippingInfoID uint           `json:"shipping_info_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ShippingInfo represents shipping information for an order
type ShippingInfo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	PostalCode  string    `json:"postal_code"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
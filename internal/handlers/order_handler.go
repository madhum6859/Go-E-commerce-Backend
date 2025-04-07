package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecommerce/internal/database"
	"github.com/yourusername/ecommerce/internal/models"
)

// OrderHandler handles order-related requests
type OrderHandler struct{}

// NewOrderHandler creates a new order handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// GetOrders returns all orders for the current user
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, _ := c.Get("userID")

	var orders []models.Order
	database.GetDB().Where("user_id = ?", userID).Preload("OrderItems.Product").Preload("ShippingInfo").Find(&orders)

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// GetOrder returns a specific order
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).Preload("OrderItems.Product").Preload("ShippingInfo").First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("userID")

	var orderData struct {
		OrderItems   []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"order_items"`
		ShippingInfo struct {
			Address     string `json:"address"`
			City        string `json:"city"`
			State       string `json:"state"`
			Country     string `json:"country"`
			PostalCode  string `json:"postal_code"`
			PhoneNumber string `json:"phone_number"`
		} `json:"shipping_info"`
	}

	if err := c.ShouldBindJSON(&orderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := database.GetDB().Begin()

	// Create shipping info
	shippingInfo := models.ShippingInfo{
		Address:     orderData.ShippingInfo.Address,
		City:        orderData.ShippingInfo.City,
		State:       orderData.ShippingInfo.State,
		Country:     orderData.ShippingInfo.Country,
		PostalCode:  orderData.ShippingInfo.PostalCode,
		PhoneNumber: orderData.ShippingInfo.PhoneNumber,
	}

	if err := tx.Create(&shippingInfo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipping info"})
		return
	}

	// Create order
	order := models.Order{
		UserID:         userID.(uint),
		ShippingInfoID: shippingInfo.ID,
		Status:         "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items and calculate total
	var totalAmount float64
	for _, item := range orderData.OrderItems {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found: " + strconv.Itoa(int(item.ProductID))})
			return
		}

		// Check if enough stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock for product: " + product.Name})
			return
		}

		// Update stock
		product.Stock -= item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
			return
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	// Update order with total amount
	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order total"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Return the created order
	var createdOrder models.Order
	database.GetDB().Preload("OrderItems.Product").Preload("ShippingInfo").First(&createdOrder, order.ID)

	c.JSON(http.StatusCreated, gin.H{"order": createdOrder})
}
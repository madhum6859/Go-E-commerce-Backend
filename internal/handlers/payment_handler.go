package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecommerce/configs"
	"github.com/yourusername/ecommerce/internal/database"
	"github.com/yourusername/ecommerce/internal/models"
	"github.com/yourusername/ecommerce/internal/payment"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	stripeService *payment.StripeService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(config *configs.Config) *PaymentHandler {
	return &PaymentHandler{
		stripeService: payment.NewStripeService(config),
	}
}

// CreatePaymentIntent creates a payment intent for an order
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var paymentData struct {
		OrderID uint `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get order
	var order models.Order
	if err := database.GetDB().Where("id = ? AND user_id = ?", paymentData.OrderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Check if order is already paid
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is already processed"})
		return
	}

	// Create payment intent
	amount := int64(order.TotalAmount * 100) // Convert to cents
	metadata := map[string]string{
		"order_id": strconv.Itoa(int(order.ID)),
		"user_id":  strconv.Itoa(int(userID.(uint))),
	}

	paymentIntent, err := h.stripeService.CreatePaymentIntent(amount, "usd", metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret": paymentIntent.ClientSecret,
		"payment_id":    paymentIntent.ID,
	})
}

// ConfirmPayment confirms a payment for an order
func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	var confirmData struct {
		PaymentIntentID string `json:"payment_intent_id" binding:"required"`
		OrderID         uint   `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&confirmData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify payment intent
	paymentIntent, err := h.stripeService.ConfirmPayment(confirmData.PaymentIntentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
		return
	}

	// Check payment status
	if paymentIntent.Status != "succeeded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not successful"})
		return
	}

	// Update order status
	var order models.Order
	if err := database.GetDB().First(&order, confirmData.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = "paid"
	order.PaymentID = paymentIntent.ID

	if err := database.GetDB().Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment confirmed successfully",
		"order":   order,
	})
}

// GetPaymentStatus gets the status of a payment
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentIntentID := c.Param("id")

	paymentIntent, err := h.stripeService.ConfirmPayment(paymentIntentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": paymentIntent.Status,
	})
}
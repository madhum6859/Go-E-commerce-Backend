package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecommerce/configs"
	"github.com/yourusername/ecommerce/internal/database"
	"github.com/yourusername/ecommerce/internal/handlers"
	"github.com/yourusername/ecommerce/internal/middleware"
	"github.com/yourusername/ecommerce/internal/models"
)

func main() {
	// Load configuration
	config := configs.LoadConfig()

	// Initialize database
	database.Initialize(config)

	// Auto migrate models
	db := database.GetDB()
	db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Image{},
		&models.Order{},
		&models.OrderItem{},
		&models.ShippingInfo{},
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(config)
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler()
	paymentHandler := handlers.NewPaymentHandler(config)

	// Set up router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// User routes
		user := api.Group("/users")
		user.Use(middleware.AuthMiddleware(config))
		{
			user.GET("/profile", userHandler.GetProfile)
		}

		// Product routes
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)

			// Admin only routes
			products.Use(middleware.AuthMiddleware(config), middleware.AdminMiddleware())
			{
				products.POST("", productHandler.CreateProduct)
				products.PUT("/:id", productHandler.UpdateProduct)
				products.DELETE("/:id", productHandler.DeleteProduct)
			}
		}

		// Order routes
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware(config))
		{
			orders.GET("", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("", orderHandler.CreateOrder)
		}

		// Payment routes
		payments := api.Group("/payments")
		payments.Use(middleware.AuthMiddleware(config))
		{
			payments.POST("/create-intent", paymentHandler.CreatePaymentIntent)
			payments.POST("/confirm", paymentHandler.ConfirmPayment)
			payments.GET("/:id", paymentHandler.GetPaymentStatus)
		}
	}

	// Start server
	log.Printf("Server starting on port %s", config.ServerPort)
	router.Run(":" + config.ServerPort)
}
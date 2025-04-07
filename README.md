# Go-E-commerce-Backend
The backend infrastructure for an online store, including features like user authentication, product catalog management, order processing, and payment gateway integration.

## API Endpoints
### Authentication
- POST /api/auth/register - Register a new user
- POST /api/auth/login - Login a user
### Users
- GET /api/users/profile - Get user profile
### Products
- GET /api/products - Get all products
- GET /api/products/:id - Get a specific product
- POST /api/products - Create a new product (admin only)
- PUT /api/products/:id - Update a product (admin only)
- DELETE /api/products/:id - Delete a product (admin only)
### Orders
- GET /api/orders - Get all orders for the current user
- GET /api/orders/:id - Get a specific order
- POST /api/orders - Create a new order
### Payments
- POST /api/payments/create-intent - Create a payment intent
- POST /api/payments/confirm - Confirm a payment
- GET /api/payments/:id - Get payment status

package payment

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/yourusername/ecommerce/configs"
)

// StripeService handles Stripe payment operations
type StripeService struct {
	config *configs.Config
}

// NewStripeService creates a new Stripe service
func NewStripeService(config *configs.Config) *StripeService {
	stripe.Key = config.StripeKey
	return &StripeService{
		config: config,
	}
}

// CreatePaymentIntent creates a payment intent for an order
func (s *StripeService) CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Metadata: metadata,
	}

	return paymentintent.New(params)
}

// ConfirmPayment confirms a payment intent
func (s *StripeService) ConfirmPayment(paymentIntentID string) (*stripe.PaymentIntent, error) {
	return paymentintent.Get(paymentIntentID, nil)
}
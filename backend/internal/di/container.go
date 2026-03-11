package di

import (
	"fmt"
	"os"
	"time"

	"sghassessment/internal/api"
	"sghassessment/internal/solutions/payment"
	"sghassessment/pkg/store"
)

// Container holds all instantiated dependencies for the application.
// It acts as a central registry for services and handlers.
type Container struct {
	SolutionsHandler *api.SolutionsHandler
	PaymentHandler   *api.PaymentHandler
}

// NewContainer initializes and wires all application dependencies.
func NewContainer() *Container {
	// 1. Load configuration / static data
	solutions, err := api.LoadSolutions()
	if err != nil {
		fmt.Printf("Warning: could not load solutions: %v\n", err)
	}

	delay := 60 * time.Second
	if raw := os.Getenv("PAYMENT_PROCESSING_DELAY"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			delay = parsed
		}
	}

	// 2. Initialize Stores
	txStore := store.New[string, payment.Transaction]()

	// 3. Initialize Services
	paymentSvc := payment.NewService(txStore, delay)

	// 4. Initialize Handlers
	solutionsHandler := api.NewSolutionsHandler(solutions)
	paymentHandler := api.NewPaymentHandler(paymentSvc)

	return &Container{
		SolutionsHandler: solutionsHandler,
		PaymentHandler:   paymentHandler,
	}
}

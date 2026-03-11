package di

import (
	"os"
	"time"

	"sghassessment/internal/api"
	"sghassessment/internal/solutions/codereview1"
	"sghassessment/internal/solutions/payment"
	"sghassessment/internal/solutions/workerpool"
	"sghassessment/pkg/logger"
	"sghassessment/pkg/store"
)

// Container holds all instantiated dependencies for the application.
// It acts as a central registry for services and handlers.
type Container struct {
	Logger             logger.Logger
	SolutionsHandler   *api.SolutionsHandler
	PaymentHandler     *api.PaymentHandler
	WorkerPoolHandler  *api.WorkerPoolHandler
	CodeReview1Handler *api.CodeReview1Handler
}

// NewContainer initializes and wires all application dependencies.
func NewContainer(log logger.Logger) *Container {
	// 1. Load configuration / static data
	log.Debug().Msg("Loading solutions metadata")
	solutions, err := api.LoadSolutions(log)
	if err != nil {
		log.Warn().Err(err).Msg("Could not load solutions")
	}

	delay := 60 * time.Second
	if raw := os.Getenv("PAYMENT_PROCESSING_DELAY"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			delay = parsed
			log.Debug().Dur("delay", delay).Msg("Payment processing delay configured")
		}
	}

	// 2. Initialize Stores
	log.Debug().Msg("Initializing data stores")
	txStore := store.New[string, payment.Transaction]()

	// 3. Initialize Services
	log.Debug().Msg("Initializing services")
	paymentSvc := payment.NewService(txStore, delay)
	workerPoolSvc := workerpool.NewService(5)

	// Code Review 1 Services
	badReview1Svc := codereview1.NewBadService()
	goodReview1Svc := codereview1.NewGoodService(1024*1024, log)
	simReview1Svc := codereview1.NewSimulatorService(badReview1Svc)

	// 4. Initialize Handlers
	log.Debug().Msg("Initializing HTTP handlers")
	solutionsHandler := api.NewSolutionsHandler(solutions)
	paymentHandler := api.NewPaymentHandler(paymentSvc, log)
	workerPoolHandler := api.NewWorkerPoolHandler(workerPoolSvc, log)
	codeReview1Handler := api.NewCodeReview1Handler(badReview1Svc, goodReview1Svc, simReview1Svc, log)

	log.Info().Msg("Dependency injection container initialized")

	return &Container{
		Logger:             log,
		SolutionsHandler:   solutionsHandler,
		PaymentHandler:     paymentHandler,
		WorkerPoolHandler:  workerPoolHandler,
		CodeReview1Handler: codeReview1Handler,
	}
}

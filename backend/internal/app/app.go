package app

import (
	"net/http"
	"os"

	"sghassessment/internal/api"
	"sghassessment/internal/di"
	"sghassessment/pkg/logger"
)

// App represents the main application containing all dependencies.
type App struct {
	router *http.ServeMux
	port   string
	logger logger.Logger
}

// New constructs the application, wires all dependencies, and sets up routes.
func New(log logger.Logger) *App {
	// Initialize the dependency injection container
	container := di.NewContainer(log)

	// Setup Router with injected handlers
	router := api.SetupRouter(
		container.SolutionsHandler,
		container.PaymentHandler,
		container.WorkerPoolHandler,
		container.CodeReview1Handler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &App{
		router: router,
		port:   port,
		logger: log,
	}
}

// Run starts the HTTP server and blocks.
func (a *App) Run() error {
	a.logger.Info().Str("port", a.port).Msg("Starting HTTP server")
	return http.ListenAndServe(":"+a.port, a.router)
}

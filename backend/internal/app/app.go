package app

import (
	"net/http"

	"sghassessment/internal/api"
	"sghassessment/internal/di"
	"sghassessment/pkg/config"
	"sghassessment/pkg/logger"
)

// App represents the main application containing all dependencies.
type App struct {
	router *http.ServeMux
	port   string
	logger logger.Logger
}

// New constructs the application, wires all dependencies, and sets up routes.
func New(cfg config.AppConfig, log logger.Logger) *App {
	// Initialize the dependency injection container
	container := di.NewContainer(cfg, log)

	// Setup Router with injected handlers
	router := api.SetupRouter(
		container.SolutionsHandler,
		container.PaymentHandler,
		container.WorkerPoolHandler,
		container.CodeReview1Handler,
		container.CodeReview2Handler,
		container.Sql1Handler,
		container.Sql2Handler,
	)

	return &App{
		router: router,
		port:   cfg.Port,
		logger: log,
	}
}

// Run starts the HTTP server and blocks.
func (a *App) Run() error {
	a.logger.Info().Str("port", a.port).Msg("Starting HTTP server")
	return http.ListenAndServe(":"+a.port, a.router)
}

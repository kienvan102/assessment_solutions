package app

import (
	"fmt"
	"net/http"
	"os"

	"sghassessment/internal/api"
	"sghassessment/internal/di"
)

// App represents the main application containing all dependencies.
type App struct {
	router *http.ServeMux
	port   string
}

// New constructs the application, wires all dependencies, and sets up routes.
func New() *App {
	// Initialize the dependency injection container
	container := di.NewContainer()

	// Setup Router with injected handlers
	router := api.SetupRouter(
		container.SolutionsHandler,
		container.PaymentHandler,
		container.WorkerPoolHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &App{
		router: router,
		port:   port,
	}
}

// Run starts the HTTP server and blocks.
func (a *App) Run() error {
	fmt.Printf("Starting server on port %s...\n", a.port)
	return http.ListenAndServe(":"+a.port, a.router)
}

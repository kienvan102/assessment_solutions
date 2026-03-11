package main

import (
	"os"

	"sghassessment/internal/app"
	"sghassessment/pkg/config"
	"sghassessment/pkg/logger"
)

func main() {
	// Initialize logger
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	log := logger.NewZerologAdapter(env)

	log.Info().Str("env", env).Msg("Starting application")

	cfg := config.LoadConfig()

	application := app.New(cfg, log)
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

package main

import (
	"os"

	"sghassessment/internal/app"
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

	application := app.New(log)
	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

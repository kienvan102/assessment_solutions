package api

import (
	"sghassessment/internal/models"
	"sghassessment/pkg/config"
	"sghassessment/pkg/logger"
)

// LoadSolutions reads solutions data from a JSON file.
// It checks for the file relative to the current working directory,
// looking in "data/solutions.json" (local dev) or "/app/data/solutions.json" (docker).
func LoadSolutions(log logger.Logger) ([]models.Solution, error) {
	// Possible paths depending on how the binary is executed.
	// In docker, workdir is usually /app, so data/solutions.json works.
	// If not found, we can try absolute paths or fallback locations.
	paths := []string{
		"data/solutions.json",
		"/app/data/solutions.json",
		"../data/solutions.json",
	}

	var solutions []models.Solution
	loadedPath, err := config.LoadJSONFile(&solutions, paths...)
	if err != nil {
		return nil, err
	}

	log.Info().Str("path", loadedPath).Int("count", len(solutions)).Msg("Loaded solutions")
	return solutions, nil
}

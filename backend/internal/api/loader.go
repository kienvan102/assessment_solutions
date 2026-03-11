package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"sghassessment/internal/models"
)

// LoadSolutions reads solutions data from a JSON file.
// It checks for the file relative to the current working directory,
// looking in "data/solutions.json" (local dev) or "/app/data/solutions.json" (docker).
func LoadSolutions() ([]models.Solution, error) {
	// Possible paths depending on how the binary is executed.
	// In docker, workdir is usually /app, so data/solutions.json works.
	// If not found, we can try absolute paths or fallback locations.
	paths := []string{
		"data/solutions.json",
		"/app/data/solutions.json",
		"../data/solutions.json",
	}

	var data []byte
	var err error
	var loadedPath string

	for _, p := range paths {
		data, err = os.ReadFile(p)
		if err == nil {
			loadedPath = p
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open solutions file, tried paths %v: %w", paths, err)
	}
	fmt.Printf("Loaded solutions from: %s\n", loadedPath)

	var solutions []models.Solution
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&solutions); err != nil {
		return nil, fmt.Errorf("failed to parse solutions file: %w", err)
	}

	return solutions, nil
}

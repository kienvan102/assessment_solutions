package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSONFile reads and parses a JSON file from multiple possible paths.
// It tries each path in order and returns the first successful read.
// The result is decoded into the provided target interface.
func LoadJSONFile(target interface{}, paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", fmt.Errorf("no paths provided")
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
		return "", fmt.Errorf("failed to open file, tried paths %v: %w", paths, err)
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(target); err != nil {
		return "", fmt.Errorf("failed to parse JSON file from %s: %w", loadedPath, err)
	}

	return loadedPath, nil
}

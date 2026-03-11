package main

import (
	"fmt"
	"os"

	"sghassessment/internal/app"
)

func main() {
	application := app.New()
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}

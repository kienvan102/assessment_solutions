package api

import (
	"encoding/json"
	"net/http"

	"sghassessment/internal/models"
)

// SolutionsHandler handles GET /api/solutions
type SolutionsHandler struct {
	solutions []models.Solution
}

func NewSolutionsHandler(solutions []models.Solution) *SolutionsHandler {
	return &SolutionsHandler{
		solutions: solutions,
	}
}

func (h *SolutionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(h.solutions)
}

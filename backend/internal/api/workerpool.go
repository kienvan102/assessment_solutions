package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sghassessment/internal/solutions/workerpool"
	"sghassessment/pkg/logger"
)

// WorkerPoolRequest represents the input for the worker pool endpoint.
type WorkerPoolRequest struct {
	NumTasks int `json:"numTasks"`
	PoolSize int `json:"poolSize"`
}

// WorkerPoolResponse represents the output from the worker pool endpoint.
type WorkerPoolResponse struct {
	Results       []workerpool.Result `json:"results"`
	ExecutionTime string              `json:"executionTime"`
	Summary       string              `json:"summary"`
}

// WorkerPoolHandler handles POST /api/workerpool requests.
type WorkerPoolHandler struct {
	service *workerpool.Service
	logger  logger.Logger
}

// NewWorkerPoolHandler creates a new worker pool handler with the given service.
func NewWorkerPoolHandler(service *workerpool.Service, log logger.Logger) *WorkerPoolHandler {
	return &WorkerPoolHandler{
		service: service,
		logger:  log,
	}
}

func (h *WorkerPoolHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn().Str("method", r.Method).Msg("Invalid HTTP method for workerpool endpoint")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WorkerPoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to decode workerpool request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.NumTasks <= 0 {
		h.logger.Warn().Int("numTasks", req.NumTasks).Msg("Invalid numTasks in workerpool request")
		http.Error(w, "numTasks must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.PoolSize <= 0 {
		h.logger.Warn().Int("poolSize", req.PoolSize).Msg("Invalid poolSize in workerpool request")
		http.Error(w, "poolSize must be greater than 0", http.StatusBadRequest)
		return
	}

	h.logger.Info().Int("numTasks", req.NumTasks).Int("poolSize", req.PoolSize).Msg("Processing workerpool request")

	// Process tasks using worker pool
	results, executionTime, err := h.service.ProcessTasks(req.NumTasks, req.PoolSize)
	if err != nil {
		h.logger.Error().Err(err).Int("numTasks", req.NumTasks).Int("poolSize", req.PoolSize).Msg("Worker pool processing failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info().Int("numTasks", req.NumTasks).Int("poolSize", req.PoolSize).Dur("executionTime", executionTime).Msg("Worker pool processing completed")

	// Build response
	response := WorkerPoolResponse{
		Results:       results,
		ExecutionTime: executionTime.String(),
		Summary:       formatSummary(req.NumTasks, req.PoolSize, executionTime.String()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func formatSummary(numTasks, poolSize int, execTime string) string {
	return fmt.Sprintf("Processed %d tasks using %d workers in %s", numTasks, poolSize, execTime)
}

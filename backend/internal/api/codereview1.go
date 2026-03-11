package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"sghassessment/internal/solutions/codereview1"
	"sghassessment/pkg/logger"
)

// CodeReview1Request represents the input for the code review interactive test.
type CodeReview1Request struct {
	Action  string `json:"action"`  // "bad", "good", "simulate", "analysis"
	Payload string `json:"payload"` // Data to send to the handler
}

// CodeReview1Handler handles POST /api/codereview1 requests.
type CodeReview1Handler struct {
	badService  *codereview1.BadService
	goodService *codereview1.GoodService
	simulator   *codereview1.SimulatorService
	logger      logger.Logger
}

func NewCodeReview1Handler(bad *codereview1.BadService, good *codereview1.GoodService, sim *codereview1.SimulatorService, log logger.Logger) *CodeReview1Handler {
	return &CodeReview1Handler{
		badService:  bad,
		goodService: good,
		simulator:   sim,
		logger:      log,
	}
}

func (h *CodeReview1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn().Str("method", r.Method).Msg("Invalid HTTP method for codereview1 endpoint")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the JSON body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req CodeReview1Request
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to decode codereview1 request")
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("action", req.Action).Msg("Processing code review 1 request")

	switch strings.ToLower(req.Action) {
	case "bad":
		// Re-construct a raw HTTP request mimicking a direct hit to the bad handler
		simReq := httptest.NewRequest(http.MethodPost, "/bad", bytes.NewReader([]byte(req.Payload)))
		
		// To demonstrate the lack of Content-Type/JSON parsing, we'll let it write directly to w
		h.badService.Handler(w, simReq)

	case "good":
		// The Good Handler expects direct raw bytes as payload too (as it runs io.ReadAll)
		simReq := httptest.NewRequest(http.MethodPost, "/good", bytes.NewReader([]byte(req.Payload)))
		
		// It will set Content-Type to application/json and write
		h.goodService.Handler(w, simReq)

	case "simulate":
		// Trigger the race condition simulation
		result := h.simulator.SimulateRace()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)

	case "analysis":
		// Return the written problems and improvements
		analysis := codereview1.GetWrittenAnalysis()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(analysis)

	default:
		h.logger.Warn().Str("action", req.Action).Msg("Unknown action requested")
		http.Error(w, "Invalid action. Use: bad, good, simulate, or analysis", http.StatusBadRequest)
	}
}

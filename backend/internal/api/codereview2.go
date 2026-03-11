package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"sghassessment/internal/solutions/codereview2"
	"sghassessment/pkg/logger"
)

// CodeReview2Request represents the input for the codereview2 interactive test.
type CodeReview2Request struct {
	Action  string `json:"action"`  // "bad", "good", "simulate", "analysis"
	Payload string `json:"payload"` // Data to send to the handler
}

// CodeReview2Handler handles POST /api/codereview2 requests.
type CodeReview2Handler struct {
	badService  *codereview2.BadService
	goodService *codereview2.GoodService
	simulator   *codereview2.SimulatorService
	logger      logger.Logger
}

func NewCodeReview2Handler(bad *codereview2.BadService, good *codereview2.GoodService, sim *codereview2.SimulatorService, log logger.Logger) *CodeReview2Handler {
	return &CodeReview2Handler{
		badService:  bad,
		goodService: good,
		simulator:   sim,
		logger:      log,
	}
}

func (h *CodeReview2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn().Str("method", r.Method).Msg("Invalid HTTP method for codereview2 endpoint")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to read body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req CodeReview2Request
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to decode codereview2 request")
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("action", req.Action).Msg("Processing code review 2 request")

	switch strings.ToLower(req.Action) {
	case "bad":
		simReq := httptest.NewRequest(http.MethodPost, "/bad", bytes.NewReader([]byte(req.Payload)))
		h.badService.Handler(w, simReq)

	case "good":
		simReq := httptest.NewRequest(http.MethodPost, "/good", bytes.NewReader([]byte(req.Payload)))
		h.goodService.Handler(w, simReq)

	case "simulate":
		result := h.simulator.SimulateRace()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)

	case "analysis":
		analysis := codereview2.GetWrittenAnalysis()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(analysis)

	default:
		h.logger.Warn().Str("action", req.Action).Msg("Unknown action requested")
		http.Error(w, "Invalid action. Use: bad, good, simulate, or analysis", http.StatusBadRequest)
	}
}

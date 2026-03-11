package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sghassessment/internal/solutions/sql1"
	"sghassessment/pkg/logger"
)

type Sql1Request struct {
	Query string `json:"query"`
}

type Sql1Response struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Error   string          `json:"error,omitempty"`
}

type Sql1Handler struct {
	service *sql1.Service
	logger  logger.Logger
}

func NewSql1Handler(svc *sql1.Service, log logger.Logger) *Sql1Handler {
	return &Sql1Handler{
		service: svc,
		logger:  log,
	}
}

func (h *Sql1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Sql1Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	req.Query = strings.TrimSpace(req.Query)
	if req.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("query", req.Query).Msg("Executing SQL query")

	cols, rows, err := h.service.ExecuteQuery(req.Query)
	
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.logger.Warn().Err(err).Msg("Query execution error")
		w.WriteHeader(http.StatusBadRequest) // User's query might be invalid syntax
		json.NewEncoder(w).Encode(Sql1Response{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Sql1Response{
		Columns: cols,
		Rows:    rows,
	})
}

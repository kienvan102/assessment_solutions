package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sghassessment/internal/solutions/sql2"
	"sghassessment/pkg/logger"
)

type Sql2Request struct {
	Action string `json:"action"` // "query", "explain", "reset", "exec"
	Query  string `json:"query"`
}

type Sql2Response struct {
	Columns []string        `json:"columns,omitempty"`
	Rows    [][]interface{} `json:"rows,omitempty"`
	Explain string          `json:"explain,omitempty"`
	Message string          `json:"message,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type Sql2Handler struct {
	service *sql2.Service
	logger  logger.Logger
}

func NewSql2Handler(svc *sql2.Service, log logger.Logger) *Sql2Handler {
	return &Sql2Handler{
		service: svc,
		logger:  log,
	}
}

func (h *Sql2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Sql2Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	switch strings.ToLower(req.Action) {
	case "query":
		req.Query = strings.TrimSpace(req.Query)
		if req.Query == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: "Query cannot be empty"})
			return
		}

		h.logger.Info().Str("query", req.Query).Msg("Executing SQL2 query")
		cols, rows, err := h.service.ExecuteQuery(req.Query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sql2Response{
			Columns: cols,
			Rows:    rows,
		})

	case "explain":
		req.Query = strings.TrimSpace(req.Query)
		if req.Query == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: "Query cannot be empty"})
			return
		}

		h.logger.Info().Str("query", req.Query).Msg("Explaining SQL2 query")
		explainStr, err := h.service.ExplainQuery(req.Query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sql2Response{
			Explain: explainStr,
		})

	case "exec":
		req.Query = strings.TrimSpace(req.Query)
		if req.Query == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: "Query cannot be empty"})
			return
		}

		h.logger.Info().Str("query", req.Query).Msg("Executing DDL query")
		if err := h.service.ExecQuery(req.Query); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Sql2Response{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sql2Response{
			Message: "Successfully executed",
		})

	case "reset":
		h.logger.Info().Msg("Resetting SQL2 database")
		if err := h.service.ResetDatabase(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Sql2Response{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Sql2Response{
			Message: "Database successfully reset and seeded",
		})

	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Sql2Response{Error: "Invalid action. Use: query, explain, reset"})
	}
}

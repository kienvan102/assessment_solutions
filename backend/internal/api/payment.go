package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sghassessment/internal/solutions/payment"
	"sghassessment/pkg/logger"
)

// PaymentRequest represents the incoming JSON body.
type PaymentRequest struct {
	UserID        string  `json:"userID"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transactionID"`
}

// PaymentResponse represents the response containing idempotency info.
type PaymentResponse struct {
	Transaction payment.Transaction `json:"transaction"`
	Replayed    bool                `json:"replayed"`
}

// PaymentHandler handles POST /api/pay
type PaymentHandler struct {
	paymentService *payment.Service
	logger         logger.Logger
}

func NewPaymentHandler(svc *payment.Service, log logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: svc,
		logger:         log,
	}
}

func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn().Str("method", r.Method).Msg("Invalid HTTP method for payment endpoint")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Failed to decode payment request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.TransactionID == "" {
		h.logger.Warn().Msg("Payment request missing transactionID")
		http.Error(w, "transactionID is required", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		h.logger.Warn().Float64("amount", req.Amount).Msg("Payment request with invalid amount")
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("userID", req.UserID).Float64("amount", req.Amount).Str("transactionID", req.TransactionID).Msg("Processing payment request")

	tx, replayed, err := h.paymentService.ProcessPayment(req.UserID, req.Amount, req.TransactionID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "required") || strings.Contains(msg, "greater than zero") {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		h.logger.Error().Err(err).Msg("Failed to process payment")
		http.Error(w, msg, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PaymentResponse{Transaction: tx, Replayed: replayed})
}

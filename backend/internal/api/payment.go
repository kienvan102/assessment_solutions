package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"sghassessment/internal/solutions/payment"
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
}

func NewPaymentHandler(svc *payment.Service) *PaymentHandler {
	return &PaymentHandler{
		paymentService: svc,
	}
}

func (h *PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, replayed, err := h.paymentService.ProcessPayment(req.UserID, req.Amount, req.TransactionID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "required") || strings.Contains(msg, "greater than zero") {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		http.Error(w, msg, http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PaymentResponse{Transaction: tx, Replayed: replayed})
}

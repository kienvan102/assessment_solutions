package payment

import (
	"errors"
	"time"
)

// Transaction represents a single payment request
type Transaction struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"` // e.g., "pending", "processed"
}

// Store defines the interface for data persistence required by the payment service.
type Store interface {
	Get(key string) (Transaction, bool)
	Set(key string, value Transaction)
	Update(key string, updateFn func(Transaction) Transaction) bool
}

// Service provides payment processing logic.
type Service struct {
	txStore         Store
	processingDelay time.Duration
}

// NewService creates a new payment service.
// It requires a Store interface implementation to implement IoC properly.
func NewService(txStore Store, delay time.Duration) *Service {
	return &Service{
		txStore:         txStore,
		processingDelay: delay,
	}
}

// ProcessPayment handles the core idempotency logic
func (s *Service) ProcessPayment(userID string, amount float64, transactionID string) (Transaction, bool, error) {
	if transactionID == "" {
		return Transaction{}, false, errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return Transaction{}, false, errors.New("amount must be greater than zero")
	}

	// Lock the whole check-and-insert operation to prevent race conditions during insertion.
	// (Our generic store's Get/Set are safe individually, but we need atomicity here).

	// Fast path check
	if existingTx, exists := s.txStore.Get(transactionID); exists {
		if existingTx.UserID != userID || existingTx.Amount != amount {
			return Transaction{}, false, errors.New("transaction ID conflict with different parameters")
		}
		return existingTx, true, nil // Idempotent replay
	}

	newTx := Transaction{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Status:        "pending",
	}

	// Store new transaction.
	s.txStore.Set(transactionID, newTx)

	// Simulate async processing
	delay := s.processingDelay
	go func(txID string) {
		if delay > 0 {
			time.Sleep(delay)
		}
		s.txStore.Update(txID, func(tx Transaction) Transaction {
			tx.Status = "processed"
			return tx
		})
	}(transactionID)

	return newTx, false, nil
}

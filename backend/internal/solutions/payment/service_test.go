package payment

import (
	"testing"
	"time"

	"sghassessment/pkg/store"
)

// ensure memory store satisfies Store interface
var _ Store = (*store.Store[string, Transaction])(nil)

func TestProcessPayment_Success(t *testing.T) {
	txStore := store.New[string, Transaction]()
	svc := NewService(txStore, 10*time.Millisecond)

	tx, replayed, err := svc.ProcessPayment("user_1", 100.50, "tx_1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if replayed {
		t.Errorf("expected replayed=false on first call, got true")
	}
	if tx.Status != "pending" {
		t.Errorf("expected status 'pending' on first call, got: %s", tx.Status)
	}

	time.Sleep(20 * time.Millisecond)
	tx2, replayed2, err2 := svc.ProcessPayment("user_1", 100.50, "tx_1")
	if err2 != nil {
		t.Fatalf("expected no error on retry, got: %v", err2)
	}
	if !replayed2 {
		t.Errorf("expected replayed=true on retry, got false")
	}
	if tx2.Status != "processed" {
		t.Errorf("expected status 'processed' after delay, got: %s", tx2.Status)
	}
}

func TestProcessPayment_Idempotency(t *testing.T) {
	txStore := store.New[string, Transaction]()
	svc := NewService(txStore, 25*time.Millisecond)

	// First request
	_, replayed1, err := svc.ProcessPayment("user_1", 100.50, "tx_1")
	if err != nil {
		t.Fatalf("expected no error on first call, got: %v", err)
	}
	if replayed1 {
		t.Fatalf("expected replayed=false on first call, got true")
	}

	// Second request with same transaction ID
	tx2, replayed2, err2 := svc.ProcessPayment("user_1", 100.50, "tx_1")
	if err2 != nil {
		t.Fatalf("expected no error on retry, got: %v", err2)
	}
	if !replayed2 {
		t.Fatalf("expected replayed=true on retry, got false")
	}
	if tx2.Status != "pending" {
		t.Fatalf("expected status 'pending' on immediate retry, got: %s", tx2.Status)
	}

	time.Sleep(35 * time.Millisecond)
	tx3, replayed3, err3 := svc.ProcessPayment("user_1", 100.50, "tx_1")
	if err3 != nil {
		t.Fatalf("expected no error on retry after delay, got: %v", err3)
	}
	if !replayed3 {
		t.Fatalf("expected replayed=true on retry after delay, got false")
	}
	if tx3.Status != "processed" {
		t.Fatalf("expected status 'processed' after delay, got: %s", tx3.Status)
	}

	// Check if only one transaction is stored
	if txStore.Len() != 1 {
		t.Errorf("expected exactly 1 transaction stored, got: %d", txStore.Len())
	}
}

func TestProcessPayment_Conflict(t *testing.T) {
	txStore := store.New[string, Transaction]()
	svc := NewService(txStore, 10*time.Millisecond)

	// First request
	_, _, _ = svc.ProcessPayment("user_1", 100.50, "tx_1")

	// Second request with same transaction ID but different amount
	_, _, err := svc.ProcessPayment("user_1", 200.00, "tx_1")
	if err == nil {
		t.Fatal("expected conflict error when parameters differ on retry, got none")
	}
}

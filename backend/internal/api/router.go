package api

import (
	"encoding/json"
	"net/http"
)

// SetupRouter creates and configures the HTTP multiplexer with all routes and dependencies injected.
func SetupRouter(
	solutionsHandler *SolutionsHandler,
	paymentHandler *PaymentHandler,
	workerPoolHandler *WorkerPoolHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Backend is running"})
	})

	mux.Handle("/api/solutions", solutionsHandler)
	mux.Handle("/api/pay", paymentHandler)
	mux.Handle("/api/workerpool", workerPoolHandler)

	return mux
}

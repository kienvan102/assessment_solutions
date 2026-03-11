package codereview1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil" // Using this specifically to demonstrate the deprecated warning
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"sghassessment/pkg/logger"
)

// ==========================================
// THE BAD IMPLEMENTATION
// ==========================================

// This is the global variable causing the race condition
var globalData = ""

// BadHandler exactly matches the problematic code.
// We use a method on an empty struct just to namespace it for the router,
// but the core logic is exactly the bad code snippet.
type BadService struct{}

func NewBadService() *BadService {
	return &BadService{}
}

func (s *BadService) Handler(w http.ResponseWriter, r *http.Request) {
	// Problem 1: No size limit
	// Problem 2: Deprecated ioutil
	// Problem 3: Ignores error
	body, _ := ioutil.ReadAll(r.Body)

	// Problem 4: Race condition on globalData
	globalData = string(body)

	// Problem 5: XSS vulnerability, no content type
	fmt.Fprintf(w, "Saved: %s", globalData)

	// Problem 6: defer after reading
	defer r.Body.Close()
}

// ==========================================
// THE GOOD IMPLEMENTATION
// ==========================================

type GoodService struct {
	maxBodySize int64
	logger      logger.Logger
}

func NewGoodService(maxBodySize int64, log logger.Logger) *GoodService {
	if maxBodySize <= 0 {
		maxBodySize = 1024 * 1024 // 1MB default
	}
	return &GoodService{
		maxBodySize: maxBodySize,
		logger:      log,
	}
}

func (s *GoodService) Handler(w http.ResponseWriter, r *http.Request) {
	// Fix: defer close immediately
	defer r.Body.Close()

	// Fix: Limit payload size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)

	// Fix: Use io.ReadAll and handle error
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to read body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fix: Validate input
	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Fix: Use local variable to avoid race conditions
	localData := string(body)

	// Fix: Set proper Content-Type and return JSON safely
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"saved":  localData,
	})
}

// ==========================================
// SIMULATION & ANALYSIS
// ==========================================

type SimulatorService struct {
	badService *BadService
}

func NewSimulatorService(bad *BadService) *SimulatorService {
	return &SimulatorService{badService: bad}
}

type SimulationResult struct {
	TotalRequests int      `json:"totalRequests"`
	Corrupted     int      `json:"corrupted"`
	RaceDetected  bool     `json:"raceDetected"`
	SampleErrors  []string `json:"sampleErrors"`
}

// SimulateRace fires 100 concurrent requests at the BadHandler to demonstrate data corruption.
func (s *SimulatorService) SimulateRace() SimulationResult {
	const numRequests = 100
	var wg sync.WaitGroup
	results := make([]string, numRequests)

	// Reset the global state to empty
	globalData = ""

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			
			// Introduce a tiny random delay to maximize race condition probability
			time.Sleep(time.Duration(index%5) * time.Millisecond)

			payload := fmt.Sprintf("DataPayload-%d", index)
			req := httptest.NewRequest(http.MethodPost, "/bad", strings.NewReader(payload))
			w := httptest.NewRecorder()

			s.badService.Handler(w, req)
			results[index] = w.Body.String()
		}(i)
	}

	wg.Wait()

	corruptedCount := 0
	var samples []string

	for i, res := range results {
		expected := fmt.Sprintf("Saved: DataPayload-%d", i)
		if res != expected {
			corruptedCount++
			if len(samples) < 5 { // collect up to 5 samples
				samples = append(samples, fmt.Sprintf("Req %d expected '%s' but got '%s'", i, expected, res))
			}
		}
	}

	return SimulationResult{
		TotalRequests: numRequests,
		Corrupted:     corruptedCount,
		RaceDetected:  corruptedCount > 0,
		SampleErrors:  samples,
	}
}

// Analysis Models
type Problem struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
}

type Improvement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AnalysisResponse struct {
	Problems     []Problem     `json:"problems"`
	Improvements []Improvement `json:"improvements"`
}

func GetWrittenAnalysis() AnalysisResponse {
	return AnalysisResponse{
		Problems: []Problem{
			{
				ID: 1, Title: "Race Condition on Global Variable", Severity: "Critical",
				Description: "The global variable 'data' is modified by concurrent requests without a mutex lock. This leads to data corruption under load.",
			},
			{
				ID: 2, Title: "Silently Ignored Error", Severity: "High",
				Description: "Using the blank identifier '_' for ioutil.ReadAll's error means network/read failures are silently ignored, causing unexpected behavior.",
			},
			{
				ID: 3, Title: "Misplaced defer", Severity: "Medium",
				Description: "defer r.Body.Close() is placed at the end of the function. It should be placed immediately after the request check to ensure it closes even if reading panics.",
			},
			{
				ID: 4, Title: "Missing Request Body Size Limit", Severity: "High",
				Description: "There is no limit on the payload size. An attacker could send a multi-gigabyte payload and cause an Out-Of-Memory (OOM) crash (Denial of Service).",
			},
			{
				ID: 5, Title: "Potential XSS via fmt.Fprintf", Severity: "Medium",
				Description: "User-supplied data is written directly to the response writer without sanitization and without a Content-Type header. Browsers might interpret this as HTML/JS.",
			},
			{
				ID: 6, Title: "Deprecated ioutil Usage", Severity: "Low",
				Description: "ioutil.ReadAll is deprecated since Go 1.16. Use io.ReadAll instead.",
			},
		},
		Improvements: []Improvement{
			{Title: "Remove Global State", Description: "Use local variables inside the handler so each request has its own isolated memory space."},
			{Title: "Add Error Handling", Description: "Explicitly check 'err != nil' from io.ReadAll and return a 400 or 500 HTTP status code."},
			{Title: "Enforce Size Limits", Description: "Wrap r.Body with http.MaxBytesReader to enforce a hard cap on request body size (e.g., 1MB)."},
			{Title: "Use JSON Responses", Description: "Set Content-Type to application/json and use encoding/json to return structured, safe responses."},
		},
	}
}

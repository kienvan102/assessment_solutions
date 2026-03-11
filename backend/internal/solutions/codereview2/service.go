package codereview2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil" // Used intentionally to match the flawed code
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

// Matches exactly the problematic `var result string`
var result = ""

type BadService struct{}

func NewBadService() *BadService {
	return &BadService{}
}

func (s *BadService) Handler(w http.ResponseWriter, r *http.Request) {
	// Directly mirroring the problematic code issues
	body, _ := ioutil.ReadAll(r.Body)
	result = string(body)
	fmt.Fprintf(w, "Saved: %s", result)
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
	// 1. Defer close immediately
	defer r.Body.Close()

	// 2. Wrap body to prevent large payload memory exhaustion (DoS)
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)

	// 3. Handle errors and use modern io.ReadAll
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to read body in codereview2")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Use a local variable to prevent race conditions
	localResult := string(body)

	// 5. Respond with proper JSON and Content-Type to prevent XSS
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saved successfully",
		"data":    localResult,
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

func (s *SimulatorService) SimulateRace() SimulationResult {
	const numRequests = 100
	var wg sync.WaitGroup
	results := make([]string, numRequests)

	// Reset global state
	result = ""

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Tiny delay to cause overlap
			time.Sleep(time.Duration(index%5) * time.Millisecond)

			payload := fmt.Sprintf("Data-%d", index)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			w := httptest.NewRecorder()

			s.badService.Handler(w, req)
			results[index] = w.Body.String()
		}(i)
	}

	wg.Wait()

	corruptedCount := 0
	var samples []string

	for i, res := range results {
		expected := fmt.Sprintf("Saved: Data-%d", i)
		if res != expected {
			corruptedCount++
			if len(samples) < 5 {
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

type Problem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
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
				ID: 1, Title: "Race Condition on Global Variable 'result'", Severity: "Critical",
				Description: "The global variable 'result' is shared across all HTTP requests. Concurrent requests will overwrite each other's data before it can be written to the response.",
			},
			{
				ID: 2, Title: "Ignored Error on ioutil.ReadAll", Severity: "High",
				Description: "Using the blank identifier '_' drops the error from ioutil.ReadAll. If the client disconnects or sends a malformed body, the code continues execution blindly.",
			},
			{
				ID: 3, Title: "No Request Body Size Limit", Severity: "High",
				Description: "Reading the entire body into memory without a limit makes the server vulnerable to Out-Of-Memory (OOM) Denial of Service attacks.",
			},
			{
				ID: 4, Title: "Misplaced defer r.Body.Close()", Severity: "Medium",
				Description: "The defer statement is at the bottom. If ioutil.ReadAll panics or causes an early exit, the body is never closed, leading to a resource leak.",
			},
			{
				ID: 5, Title: "Cross-Site Scripting (XSS) Vulnerability", Severity: "Medium",
				Description: "Reflecting unescaped user input directly via fmt.Fprintf without setting a Content-Type allows browsers to execute malicious scripts.",
			},
			{
				ID: 6, Title: "Deprecated Package (ioutil)", Severity: "Low",
				Description: "The io/ioutil package has been deprecated since Go 1.16. Modern Go uses io.ReadAll.",
			},
		},
		Improvements: []Improvement{
			{Title: "Remove Global State", Description: "Declare 'result' as a local variable inside the handler function."},
			{Title: "Handle Errors Explicitly", Description: "Check 'err != nil' from ReadAll and return a proper HTTP 400/500 error if it fails."},
			{Title: "Limit Body Size", Description: "Use http.MaxBytesReader to enforce a strict memory cap on incoming payloads."},
			{Title: "Use JSON Responses", Description: "Return structured JSON data and set 'Content-Type: application/json' to neutralize XSS vectors."},
			{Title: "Fix Defer Position", Description: "Place 'defer r.Body.Close()' immediately at the top of the handler."},
		},
	}
}

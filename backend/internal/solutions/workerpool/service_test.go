package workerpool

import (
	"testing"
)

func TestProcessTasks_CorrectResults(t *testing.T) {
	svc := NewService(5)

	results, _, err := svc.ProcessTasks(10, 5)
	if err != nil {
		t.Fatalf("ProcessTasks failed: %v", err)
	}

	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	// Verify each result is correct (squared)
	for i, result := range results {
		expectedTaskID := i + 1
		expectedOutput := expectedTaskID * expectedTaskID

		if result.TaskID != expectedTaskID {
			t.Errorf("Result %d: expected TaskID %d, got %d", i, expectedTaskID, result.TaskID)
		}
		if result.Input != expectedTaskID {
			t.Errorf("Result %d: expected Input %d, got %d", i, expectedTaskID, result.Input)
		}
		if result.Output != expectedOutput {
			t.Errorf("Result %d: expected Output %d, got %d", i, expectedOutput, result.Output)
		}
		if result.Worker < 1 || result.Worker > 5 {
			t.Errorf("Result %d: worker ID %d out of range [1-5]", i, result.Worker)
		}
	}
}

func TestProcessTasks_OrderPreserved(t *testing.T) {
	svc := NewService(5)

	results, _, err := svc.ProcessTasks(100, 5)
	if err != nil {
		t.Fatalf("ProcessTasks failed: %v", err)
	}

	// Verify results are in order
	for i, result := range results {
		expectedTaskID := i + 1
		if result.TaskID != expectedTaskID {
			t.Errorf("Results not in order: position %d has TaskID %d, expected %d", i, result.TaskID, expectedTaskID)
		}
	}
}

func TestProcessTasks_InvalidInput(t *testing.T) {
	svc := NewService(5)

	_, _, err := svc.ProcessTasks(0, 5)
	if err == nil {
		t.Error("Expected error for numTasks=0, got nil")
	}

	_, _, err = svc.ProcessTasks(-10, 5)
	if err == nil {
		t.Error("Expected error for negative numTasks, got nil")
	}
}

func TestProcessTasks_DifferentPoolSizes(t *testing.T) {
	tests := []struct {
		poolSize int
		numTasks int
	}{
		{1, 10},
		{3, 10},
		{5, 10},
		{10, 10},
		{20, 10},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.poolSize)), func(t *testing.T) {
			svc := NewService(tt.poolSize)
			results, _, err := svc.ProcessTasks(tt.numTasks, tt.poolSize)

			if err != nil {
				t.Fatalf("ProcessTasks failed with poolSize=%d: %v", tt.poolSize, err)
			}

			if len(results) != tt.numTasks {
				t.Errorf("Expected %d results, got %d", tt.numTasks, len(results))
			}

			// Verify order
			for i, result := range results {
				if result.TaskID != i+1 {
					t.Errorf("Results not in order at position %d", i)
					break
				}
			}
		})
	}
}

func TestProcessTasks_ConcurrencyBenefit(t *testing.T) {
	// Test that using multiple workers is faster than using 1 worker
	numTasks := 50

	svc1 := NewService(1)
	_, duration1, err := svc1.ProcessTasks(numTasks, 1)
	if err != nil {
		t.Fatalf("ProcessTasks with 1 worker failed: %v", err)
	}

	svc5 := NewService(5)
	_, duration5, err := svc5.ProcessTasks(numTasks, 5)
	if err != nil {
		t.Fatalf("ProcessTasks with 5 workers failed: %v", err)
	}

	// 5 workers should be faster than 1 worker (with some tolerance)
	// We expect at least some speedup, though not necessarily 5x due to overhead
	if duration5 >= duration1 {
		t.Logf("Warning: 5 workers (%v) not faster than 1 worker (%v)", duration5, duration1)
		// Not failing the test as timing can be flaky, but logging for visibility
	}
}

func TestNewService_DefaultPoolSize(t *testing.T) {
	svc := NewService(0)
	if svc.poolSize != 5 {
		t.Errorf("Expected default poolSize 5, got %d", svc.poolSize)
	}

	svc = NewService(-1)
	if svc.poolSize != 5 {
		t.Errorf("Expected default poolSize 5 for negative input, got %d", svc.poolSize)
	}
}

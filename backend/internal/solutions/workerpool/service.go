package workerpool

import (
	"fmt"
	"sort"
	"time"
)

// Task represents a single unit of work to be processed.
type Task struct {
	ID    int
	Input int
}

// Result represents the output of a processed task.
type Result struct {
	TaskID int `json:"taskId"`
	Input  int `json:"input"`
	Output int `json:"output"`
	Worker int `json:"worker"`
}

// Service manages the worker pool execution.
type Service struct {
	poolSize int
}

// NewService creates a new worker pool service with the specified pool size.
func NewService(poolSize int) *Service {
	if poolSize <= 0 {
		poolSize = 5 // default
	}
	return &Service{
		poolSize: poolSize,
	}
}

// ProcessTasks processes the given number of tasks using a worker pool.
// It returns the results in order, the execution time, and any error.
func (s *Service) ProcessTasks(numTasks, poolSize int) ([]Result, time.Duration, error) {
	if numTasks <= 0 {
		return nil, 0, fmt.Errorf("numTasks must be greater than 0")
	}
	if poolSize <= 0 {
		poolSize = s.poolSize
	}

	start := time.Now()

	// Create channels for tasks and results
	tasks := make(chan Task, numTasks)
	results := make(chan Result, numTasks)

	// Spawn worker goroutines
	for w := 1; w <= poolSize; w++ {
		go s.worker(w, tasks, results)
	}

	// Send all tasks to the channel
	for i := 1; i <= numTasks; i++ {
		tasks <- Task{
			ID:    i,
			Input: i,
		}
	}
	close(tasks)

	// Collect all results
	resultSlice := make([]Result, 0, numTasks)
	for i := 0; i < numTasks; i++ {
		resultSlice = append(resultSlice, <-results)
	}
	close(results)

	// Sort results by TaskID to maintain order
	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].TaskID < resultSlice[j].TaskID
	})

	executionTime := time.Since(start)
	return resultSlice, executionTime, nil
}

// worker processes tasks from the tasks channel and sends results to the results channel.
func (s *Service) worker(id int, tasks <-chan Task, results chan<- Result) {
	for task := range tasks {
		// Simulate some work (squaring the number)
		// Add a tiny delay to make concurrency more visible
		time.Sleep(1 * time.Millisecond)
		
		output := task.Input * task.Input
		
		results <- Result{
			TaskID: task.ID,
			Input:  task.Input,
			Output: output,
			Worker: id,
		}
	}
}

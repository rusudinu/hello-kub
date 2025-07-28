package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// fibonacciHandler computes Fibonacci numbers for a specified duration to generate CPU load
func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	// Extract number from URL path /fib/{duration}
	path := strings.TrimPrefix(r.URL.Path, "/fib/")
	if path == "" || path == r.URL.Path {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Please provide duration in minutes in the format /fib/{minutes}"}`)
		return
	}

	// Parse the duration in minutes
	durationStr := path
	durationMins, err := strconv.ParseInt(durationStr, 10, 64)
	if err != nil || durationMins <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid duration. Please provide a positive integer (minutes)."}`)
		return
	}

	// Get number of workers (default to number of CPUs)
	workersStr := r.URL.Query().Get("workers")
	workers := runtime.NumCPU()
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil && w > 0 && w <= 20 {
			workers = w
		}
	}

	// No maximum limit - let Kubernetes handle scaling!

	fmt.Printf("Computing Fibonacci numbers with %d workers for %d minutes\n", workers, durationMins)
	start := time.Now()

	// Compute Fibonacci numbers for specified duration
	result := computeFibonacciLoad(durationMins, workers)

	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, `{
        "duration_requested": "%dm",
        "duration_actual": "%v",
        "workers": %d,
        "total_computations": %d,
        "computations_per_second": %.2f,
        "largest_fibonacci_computed": "%s",
        "cpu_cores": %d
    }`, durationMins, elapsed, workers, result.TotalComputations, float64(result.TotalComputations)/elapsed.Seconds(), result.LargestFib.String(), runtime.NumCPU())
}

type FibResult struct {
	TotalComputations int64
	LargestFib        *big.Int
}

// computeFibonacciLoad generates CPU load by computing Fibonacci numbers for a duration
func computeFibonacciLoad(durationMins int64, workers int) FibResult {
	duration := time.Duration(durationMins) * time.Minute
	endTime := time.Now().Add(duration)

	// Channel to collect results from workers
	results := make(chan FibResult, workers)

	// Start workers
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			computations := int64(0)
			largestFib := big.NewInt(0)

			// Keep computing Fibonacci numbers until time is up
			for time.Now().Before(endTime) {
				// Compute a batch of Fibonacci numbers (memory efficient)
				batchSize := 1000 // Compute 1000 numbers then restart sequence
				a, b := big.NewInt(0), big.NewInt(1)

				for j := 0; j < batchSize && time.Now().Before(endTime); j++ {
					// Calculate next Fibonacci number
					next := new(big.Int)
					next.Add(a, b)

					// Keep track of largest number computed
					if next.Cmp(largestFib) > 0 {
						largestFib.Set(next)
					}

					// Move to next iteration
					a.Set(b)
					b.Set(next)

					computations++
				}
			}

			fmt.Printf("Worker %d completed %d computations\n", workerID, computations)
			results <- FibResult{TotalComputations: computations, LargestFib: largestFib}
		}(i)
	}

	// Collect results from all workers
	totalComputations := int64(0)
	overallLargest := big.NewInt(0)

	for i := 0; i < workers; i++ {
		result := <-results
		totalComputations += result.TotalComputations
		if result.LargestFib.Cmp(overallLargest) > 0 {
			overallLargest.Set(result.LargestFib)
		}
	}

	return FibResult{
		TotalComputations: totalComputations,
		LargestFib:        overallLargest,
	}
}

func heavyComputeHandler(w http.ResponseWriter, r *http.Request) {
	// Get duration parameter (default 10 seconds)
	durationStr := r.URL.Query().Get("duration")
	duration := 10
	if durationStr != "" {
		if d, err := strconv.Atoi(durationStr); err == nil && d > 0 && d <= 60 {
			duration = d
		}
	}

	// Get number of goroutines (default to number of CPUs)
	workersStr := r.URL.Query().Get("workers")
	workers := runtime.NumCPU()
	if workersStr != "" {
		if w, err := strconv.Atoi(workersStr); err == nil && w > 0 && w <= 10 {
			workers = w
		}
	}

	fmt.Printf("Starting heavy computation: %d workers for %d seconds\n", workers, duration)

	start := time.Now()
	done := make(chan bool, workers)

	// Start CPU-intensive work in multiple goroutines
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			result := 0.0
			iterations := 0
			endTime := start.Add(time.Duration(duration) * time.Second)

			for time.Now().Before(endTime) {
				// CPU-intensive mathematical operations
				for j := 0; j < 100000; j++ {
					result += math.Sin(float64(j)) * math.Cos(float64(j))
					result = math.Sqrt(math.Abs(result))
					iterations++
				}
			}

			fmt.Printf("Worker %d completed %d iterations\n", workerID, iterations)
			done <- true
		}(i)
	}

	// Wait for all workers to complete
	for i := 0; i < workers; i++ {
		<-done
	}

	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
       "message": "Heavy computation completed",
       "duration_requested": "%ds",
       "duration_actual": "%v",
       "workers": %d,
       "cpu_cores": %d
    }`, duration, elapsed, workers, runtime.NumCPU())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/fib/", fibonacciHandler)
	http.HandleFunc("/heavy", heavyComputeHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET /              - Hello World")
	fmt.Println("  GET /health        - Health check")
	fmt.Println("  GET /fib/{minutes} - CPU-intensive Fibonacci computation for N minutes")
	fmt.Println("    ?workers=N       - Use N workers (1-20, default CPU count)")
	fmt.Println("  GET /heavy         - CPU intensive task")
	fmt.Println("    ?duration=N      - Run for N seconds (1-60, default 10)")
	fmt.Println("    ?workers=N       - Use N workers (1-10, default CPU count)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

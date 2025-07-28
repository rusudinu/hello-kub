package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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
	http.HandleFunc("/heavy", heavyComputeHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET /         - Hello World")
	fmt.Println("  GET /health   - Health check")
	fmt.Println("  GET /heavy    - CPU intensive task")
	fmt.Println("    ?duration=N - Run for N seconds (1-60, default 10)")
	fmt.Println("    ?workers=N  - Use N workers (1-10, default CPU count)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

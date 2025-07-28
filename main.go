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

// fibonacciHandler computes Fibonacci sequence up to a given number
func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	// Extract number from URL path /fib/{number}
	path := strings.TrimPrefix(r.URL.Path, "/fib/")
	if path == "" || path == r.URL.Path {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Please provide a number in the format /fib/{number}"}`)
		return
	}

	// Parse the target number
	targetStr := path
	target, err := strconv.ParseInt(targetStr, 10, 64)
	if err != nil || target < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Invalid number. Please provide a positive integer."}`)
		return
	}

	// Set reasonable limit to prevent server overload
	if target > 10000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "Number too large. Maximum allowed is 10000."}`)
		return
	}

	fmt.Printf("Computing Fibonacci sequence up to %d\n", target)
	start := time.Now()

	// Compute Fibonacci sequence up to target
	sequence := computeFibonacci(target)

	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")

	// Format response
	sequenceStr := "["
	for i, num := range sequence {
		if i > 0 {
			sequenceStr += ", "
		}
		sequenceStr += num.String()
	}
	sequenceStr += "]"

	fmt.Fprintf(w, `{
        "target": %d,
        "count": %d,
        "sequence": %s,
        "computation_time": "%v",
        "largest_number": "%s"
    }`, target, len(sequence), sequenceStr, elapsed, sequence[len(sequence)-1].String())
}

// computeFibonacci generates Fibonacci numbers up to (and including) the target
func computeFibonacci(target int64) []*big.Int {
	if target < 0 {
		return []*big.Int{}
	}

	sequence := []*big.Int{}

	// Handle base cases
	if target >= 0 {
		sequence = append(sequence, big.NewInt(0))
	}
	if target >= 1 {
		sequence = append(sequence, big.NewInt(1))
	}

	// Generate Fibonacci numbers
	a, b := big.NewInt(0), big.NewInt(1)

	for {
		// Calculate next Fibonacci number
		next := new(big.Int)
		next.Add(a, b)

		// Check if we've exceeded the target
		if next.Cmp(big.NewInt(target)) > 0 {
			break
		}

		sequence = append(sequence, new(big.Int).Set(next))

		// Move to next iteration
		a.Set(b)
		b.Set(next)
	}

	return sequence
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
	fmt.Println("  GET /fib/{number}  - Fibonacci sequence up to number")
	fmt.Println("  GET /heavy         - CPU intensive task")
	fmt.Println("    ?duration=N      - Run for N seconds (1-60, default 10)")
	fmt.Println("    ?workers=N       - Use N workers (1-10, default CPU count)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

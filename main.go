package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

const (
	minNum = 1
	maxNum = 100000
	// Number of worker goroutines to use
	numWorkers = 4
	// How often to report progress (in milliseconds)
	progressInterval = 10
)

// isPrime checks if a number is prime using trial division
// This is optimized to check only up to the square root of n
// and to skip even numbers after 2
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Check divisibility by numbers of form 6k±1 up to sqrt(n)
	i := 5
	for i*i <= n {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}

// worker processes numbers from the jobs channel and sends prime numbers to the results channel
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup, processed *atomic.Int64) {
	defer wg.Done()

	for n := range jobs {
		if isPrime(n) {
			results <- n
		}
		processed.Add(1)
	}
}

// reportProgress periodically logs the progress of the prime finding operation
func reportProgress(ctx context.Context, processed *atomic.Int64, total int) {
	ticker := time.NewTicker(progressInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			current := processed.Load()
			percentage := float64(current) / float64(total) * 100
			fmt.Printf("Progress: %.2f%% (%d/%d)\r", percentage, current, total)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		fmt.Println("\nCancellation requested. Shutting down...")
		cancel()
	}()

	total := maxNum - minNum + 1

	// Channels for jobs and results
	jobs := make(chan int, numWorkers*2) // Buffer the jobs channel
	results := make(chan int, total/10)  // Buffer with a reasonable size

	// Create atomic counter for progress tracking
	processed := &atomic.Int64{}

	go reportProgress(ctx, processed, total)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(i, jobs, results, &wg, processed)
	}

	// Create a channel to signal when all results are collected
	done := make(chan struct{})

	// Collect results in a separate goroutine
	primes := make([]int, 0, total/10) // Pre-allocate with a reasonable capacity
	go func() {
		for prime := range results {
			primes = append(primes, prime)
		}
		close(done)
	}()

	// Send jobs
jobLoop:
	for n := minNum; n <= maxNum; n++ {
		select {
		case <-ctx.Done():
			break jobLoop
		case jobs <- n:
			// Job sent successfully
		}
	}
	close(jobs)

	wg.Wait()

	close(results)

	// Wait for all results to be collected
	<-done

	// Print final results
	fmt.Printf("\nFound %d prime numbers between %d and %d\n", len(primes), minNum, maxNum)
}

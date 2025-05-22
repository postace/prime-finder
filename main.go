package main

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

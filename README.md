# Prime Finder

A high-performance Go application that finds prime numbers within a specified range using concurrent processing.

## Overview

This program finds all prime numbers between a minimum and maximum value (default: 1 to 100,000) using multiple goroutines for parallel processing. It features:

- Concurrent prime number calculation with configurable worker count
- Optimized prime number checking algorithm
- Real-time progress reporting
- Graceful shutdown with CTRL+C

## Requirements

- Go 1.16 or higher

## Running the Program

You can run the program with:

```bash
go run main.go
```

The program will immediately start finding prime numbers between 1 and 100,000 and display a progress bar.

## Configuration

You can modify the following constants in `main.go` to configure the program:

- `minNum`: Minimum number to start searching from (default: 1)
- `maxNum`: Maximum number to search up to (default: 100,000)
- `numWorkers`: Number of concurrent worker goroutines (default: 4)
- `progressInterval`: How often to update the progress display in milliseconds (default: 10)

## How It Works

The program uses the following approach:

1. Creates a pool of worker goroutines that check if numbers are prime
2. Distributes numbers to check across the workers via a channel
3. Collects prime numbers in a separate goroutine
4. Reports progress in real-time
5. Supports graceful cancellation via CTRL+C

## Performance Considerations

- Increasing `numWorkers` may improve performance on multi-core systems
- The prime checking algorithm is optimized for efficiency
- The program uses buffered channels to minimize goroutine blocking

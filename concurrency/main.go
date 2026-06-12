package main

// allof them ...wait groups, worker pools and channels

import (
	"fmt"
	"sync"
	"time"
)

// Email represents our newsletter data payload
type Email struct {
	ID      int
	Address string
}

// emailWorker processes jobs from the channel concurrently
func emailWorker(workerID int, jobs <-chan Email, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement counter when this worker completely runs out of jobs

	// Range loops over the channel, automatically waiting for data
	// It stops cleanly only when close(jobs) is called in main
	for email := range jobs {
		fmt.Printf("[Worker %d] ➔ Sending newsletter to %s...\n", workerID, email.Address)
		
		time.Sleep(200 * time.Millisecond) // Simulate network/email delivery latency
		
		fmt.Printf("[Worker %d] ✓ Successfully sent to %s\n", workerID, email.Address)
	}
}

func main() {
	start := time.Now()

	subscribers := []Email{
		{ID: 1, Address: "alice@example.com"},
		{ID: 2, Address: "bob@example.com"},
		{ID: 3, Address: "charlie@example.com"},
		{ID: 4, Address: "david@example.com"},
	}

	// Create a jobs queue channel with a buffer size matching our subscribers count
	jobQueue := make(chan Email, len(subscribers))
	
	var wg sync.WaitGroup

	// Start EXACTLY 2 workers. This strictly limits our concurrency pool size.
	fmt.Println("Deploying 2 email workers to handle the queue...")
	for i := 1; i <= 2; i++ {
		wg.Add(1) // Track this worker
		go emailWorker(i, jobQueue, &wg)
	}

	// Shove all our email jobs into the queue channel immediately
	for _, sub := range subscribers {
		jobQueue <- sub
	}
	
	// Close the channel to tell workers: "No more emails are coming, shut down when done"
	close(jobQueue)

	// Block main here until BOTH workers completely finish processing the entire queue
	wg.Wait()

	fmt.Printf("\nBroadcast completed successfully! Total time: %v\n", time.Since(start))
}

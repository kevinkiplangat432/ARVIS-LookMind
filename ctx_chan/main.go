package main

import (
	"context"
	"fmt"
	"time"
)

func streamTokens(ctx context.Context, datachan chan int) {
	token := 1
	for {
		select {
		case datachan <- token:
			token++
			time.Sleep(50 * time.Millisecond)

		// Context native shutdown signal (Done() is a channel that closes on timeout)
		case <-ctx.Done():
			fmt.Printf("Context timed out. Safely stopping goroutine: %v\n", ctx.Err())
			return
		}
	}
}

func main() {
	datachan := make(chan int)

	// Create a context that automatically cancels after 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel() // Good practice to clean up resources

	go streamTokens(ctx, datachan)

	for {
		select {
		case token := <-datachan:
			fmt.Printf("Received token %d\n", token)

		// Main loop also stops when the context times out
		case <-ctx.Done():
			fmt.Println("\n[KILL SWITCH] 300ms reached. Main loop exiting.")
			return
		}
	}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// the worker function
func worker(workerID int, jobsChan chan int, wg *sync.WaitGroup) {
	defer wg.Done() // signal the planner when this worker completely done

	// each worker pulls task off the conveyor belt automatically untill it closes
	for shardID := range jobsChan {
		fmt.Printf("[Worker %d] Started fetching shard %d...\n", workerID, shardID)
		time.Sleep(200 * time.Millisecond) // Simulate the 200ms network request
		fmt.Printf("[Worker %d] Finished shard %d\n", workerID, shardID)
	}

}

func main() {
	numJobs := 10
	numworkers := 3

	// create a buffered channel to hold queue of the jobs
	jobsChan := make(chan int, numJobs)
	var wg sync.WaitGroup

	//spwan 3 workers into the background
	for i :=1; i <= numworkers; i++ {
		wg.Add(1) // add this worker to the checklist
		go worker(i, jobsChan, &wg)

	}

	// place the 10shards onto the belt
	for j := 101; j<=110; j++ {
		jobsChan <- j
	}
	// close the channel this signals the worker no more new jobs are allowed
	close(jobsChan)

	// wait for all 3 workers to finish their current jobs 
	wg.Wait()
	fmt.Println("All 10 shards processed safely using only 3 concurrent workers!")

}
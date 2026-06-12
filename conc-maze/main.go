package main

// wait groups 


// same usual stuff imports 
import (
	"fmt"
	"sync"
	"time"
)

// define a function that takes in an int and a waitgroup pointer 

func fetchShardData(id int, wg *sync.WaitGroup) {

	defer wg.Done() // execute when done ...start unpacking the "checklist"

	fmt.Printf("Fetching data from shard %d...\n", id)

	//simulate the real world senario where there is an activity ...in a realprod this will be a no no 
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Done fetching shard %d\n", id)	
}


func main(){

    start := time.Now() // start the timer, use later to get the speed of the call and latency 

	// create a slice
	var shards []int
	// append to the slice a 1000 elemenmts 
	for i := 1001; i <= 2000; i++ {
		shards = append(shards, i)
	}

	//  define the wait group ...
	var wg sync.WaitGroup
	
	// the _ is for the index ..which we do not need. 
	// fetch the shard from the slice and ignore the index 
	for _, shard := range shards {
		// add it to the checklist 
		wg.Add(1)
		// call the fetch ...pass in the shard and the wait group's memory address 
		go fetchShardData(shard, &wg)
	}
	

	wg.Wait()
	fmt.Println("All databases synchronized")
	fmt.Printf("\nBroadcast completed successfully! Total time: %v\n", time.Since(start))
}
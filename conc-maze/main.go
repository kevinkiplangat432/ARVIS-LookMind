package main

import (
	"fmt"
	"sync"
	"time"
)
func fetchShardData(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Fetching data from shard %d...\n", id)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Done fetching shard %d\n", id)	
}
func main(){
    start := time.Now()

	// create a slice
	var shards []int
	for i := 1001; i <= 2000; i++ {
		shards = append(shards, i)
	}
	// 
	var wg sync.WaitGroup

	for _, shard := range shards {
		wg.Add(1)

		go fetchShardData(shard, &wg)
	}
	wg.Wait()
	fmt.Println("All databases synchronized")
	fmt.Printf("\nBroadcast completed successfully! Total time: %v\n", time.Since(start))
}
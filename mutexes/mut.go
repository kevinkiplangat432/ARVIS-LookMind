package main
// mutexes ...  are like the gate keepers ...they decide what happens in the storage who reads and who writes  ...and when they need "privacy"

import (
	"fmt"
	"sync"
	"time"

)

// rule store to hold our applications configurations
type RuleStore struct{
	mu sync.RWMutex // Protect the map underneath
	rules map[string]string
}

// get safety reads a rule using a read lock 
func (rs *RuleStore) Get(key  string) string {
	rs.mu.RLock()  // mutliple goroutines can hold an Rlock at the same time
	defer rs.mu.RUnlock()
	return rs.rules[key]
} 

//  set safely updates a rule using full lock
func (rs *RuleStore) Set(key string,  value string) {
	rs.mu.Lock() // only one goroutine can hold a Lock at a time
	defer rs.mu.Unlock()
	rs.rules[key]= value
}

func main() {
	store := &RuleStore{
		rules: make(map[string]string),

	}

	store.Set("Maintenace_mode", "false")

	var wg sync.WaitGroup

	// spawn 5 concurrent reader workers
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(workerID int) { 
			defer wg.Done()
			for j := 0; j < 10; j++ {
				status := store.Get("maintenance_mode")
				fmt.Printf("[Worker %d ] Read maintenance mode: %s\n", workerID, status)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// spawn 1 writer worker that chnages the values mid-execution
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(25 * time.Millisecond) // wait for readers to start
		fmt.Println("!!! [ADMIN] Updating maintenance mode to true !!!")
		store.Set("maintenance _mode", "true")

	}()

	wg.Wait()
	fmt.Println("All routines finished successfully")
}
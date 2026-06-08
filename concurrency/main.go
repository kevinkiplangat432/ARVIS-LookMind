package main

import (
	"fmt"
	"sync"
	"time"
)

// email represent our newsletter data payload
type Email struct {
	ID int
	Address string
}
// oass a pointer so the wait group can modify the original counter
func sendEmail(email Email, wg *sync.WaitGroup){

	// defer tells go to run wg.Done() right before this function exits
	defer wg.Done()
	// simulate network latency for sending an email
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Newsletter send to %s\n",email.Address)

}

func main() {
	start := time.Now()

	// a mock list of subscribers 

	subscribers := []Email{
		{ID: 1, Address: "alice@example.com"},
		{ID: 2, Address: "bob@example.com"},
		{ID: 3, Address: "charlie@example.com"},
		{ID: 4, Address: "david@example.com"},

	}
	// initialize the wait group
	var wg sync.WaitGroup

	fmt.Println("Starting newsletter broadcast ...")
	//native concurrent implementation (fire and forget)
	for _, sub := range subscribers{
		// increment the wait group counter for each email we plan to send
		wg.Add(1)
		go sendEmail(sub, &wg)
	}
	//block main here untill the counter goes down to 0
	wg.Wait()

	fmt.Printf("Broadcast routine finished in %v\n", time.Since(start))

	
}
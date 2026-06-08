package main

import (
	"fmt"
	"time"
)

// email represent our newsletter data payload
type Email struct {
	ID int
	Address string
}
// oass a pointer so the wait group can modify the original counter
// refactor v2 pass a channel (chan string ) into the function
func sendEmail(email Email, resultChan chan string){

	time.Sleep(100 * time.Millisecond)
	//create a success message string 
	status := fmt.Sprintf("Newsletter send to %s\n",email.Address)

	// shove the massage into the channel using the <- arrow

	resultChan <- status

}

func main() {
	start := time.Now()

	// a mock list of subscribers 
	// buffered channel 
	limitChan := make(chan int, 2)

	totalJobs := 4 

	fmt.Println( "starting rate-limited newsletter broadcast.....")

	for i := 1; i <= totalJobs; i++ {
		// shove a token into the channel buffer before starting the worker 
		// if the waiting room is full it blocks and waits

		limitChan <- i
		fmt.Printf("[Slot] Job %d entered the waiting room\n", i)

		go func(jobID int) {
			// this anonymous function represents our email sender 
			fmt.Printf("  [➔] Processing Email Job %d...\n", jobID)
			time.Sleep(500 * time.Millisecond) // Simulated heavy work
			fmt.Printf("  [✓] Finished Email Job %d\n", jobID)

			// take the token out of the channel when done, freeing up a slot
			<-limitChan

		}(i)

		// a quickie cheat to let the final background jobs print before exiting
		time.Sleep(1500 * time.Millisecond)
		fmt.Printf("Broadcast done! Total time: %v \n ", time.Since(start))
	}
	subscribers := []Email{
		{ID: 1, Address: "alice@example.com"},
		{ID: 2, Address: "bob@example.com"},
		{ID: 3, Address: "charlie@example.com"},
		{ID: 4, Address: "david@example.com"},

	}
	// create a channel that only holds "string" data types 
	resultChan := make(chan string)


	fmt.Println("Starting newsletter broadcast ...")

	//native concurrent implementation (fire and forget)
	for _, sub := range subscribers{
		// increment the wait group counter for each email we plan to send
		go sendEmail(sub, resultChan)
	}

	// pull data out of the channel 4 times 
	for i := 0; i < 5; i++ {
		// This line BLocks and wait untill a go routine shows sth
		msg := <- resultChan
		fmt.Println(msg)
	}

	fmt.Printf("Broadcast routine finished in %v\n", time.Since(start))

	
}
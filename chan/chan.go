package main
// channels, timeouts, data streams

//imports 
import (
	"time"
	"fmt"
)
// the func streamtokens takes in two channels ...one channels is of type int and the other channel is of type bool

// it simulates a never ending data stream

func streamTokens(datachan chan int, quitchan chan bool) {

	token := 1
	// start an infinite loop
	for {
		select{clear 
		case datachan <- token:
			
			token ++
			time.Sleep(50 * time.Millisecond) // sleep for 50 ms so that we can see it 


		// place in another channels that signals that the stream should stop 
		// this case is only selected if quitchan is true ...meaning when we receive a signal.
		case _, ok := <-quitchan:
			if !ok {
				// stop the goroutine if we receive a signal on the quitchan
				fmt.Printf("Received quit signal. Stopping goroutine.\n")
				return
			}
		}
	}
}

// define main 
func main() {
	// create two channels 
	datachan := make(chan int) // assign type int

	quitchan := make(chan bool) // assign type bool

	go streamTokens(datachan, quitchan) // call our function that takes in the channels 

	timeout := time.After(300 * time.Millisecond) // set a timeout so that if the channels do not receive and stream then we can close them 

	// an infinite loop 
	for {

		select{
		case token :=  <- datachan: // Did we receive a token? if True we continue and print it.
			// door one opens we successfully get a token from the worker
			fmt.Printf("Received token %d\n", token)

		case <- timeout: // is the timeout true .... then we send a signal to quitchan 
			//door 2 opens the 300ms clock ran out
			fmt.Println("\n[KILL SWITCH ACTIVATED ] 300ms reached. Shutting down system!")

			close(quitchan) // close the channel

			time.Sleep(10 * time.Millisecond) // give time for the channels to finish closing 
			
			return
		}
	}

}
package main

import (
	"time"
	"fmt"
)

func streamTokens(datachan chan int, quitchan chan bool) {
	token := 1
	for {
		select{
		case datachan <- token:
			token ++
			time.Sleep(50 * time.Millisecond)

		case <- quitchan:
			// stop the goroutine if we receive a signal on the quitchan
			fmt.Printf("Received quit signal. Stopping goroutine.\n")
			return
		}
	}
}
func main() {

	datachan := make(chan int)

	quitchan := make(chan bool)

	go streamTokens(datachan, quitchan)
	timeout := time.After(300 * time.Millisecond)


	for {

		select{
		case token :=  <- datachan:
			// door one opens we successfully get a token from the worker
			fmt.Printf("Received token %d\n", token)

		case <- timeout:
			//door 2 opens the 300ms clock ran out
			fmt.Println("\n[KILL SWITCH ACTIVATED ] 300ms reached. Shutting down system!")

			close(quitchan)
			time.Sleep(10 * time.Millisecond)
			return
		}
	}

}
package main

import (
	"time"
	"fmt"
)

func streamTokens(datachan chan int) {
	token := 1
	for {
		datachan <- token
		token ++
		time.Sleep(50 * time.Millisecond)
	}
}
func main() {

	datachan := make(chan int)

	go streamTokens(datachan)


	for {
		timeout := time.After(300 * time.Millisecond)

		select{
		case token :=  <- datachan:
			// door one opens we successfully get a token from the worker
			fmt.Printf("Received token %d\n", token)

		case <- timeout:
			//door 2 opens the 300ms clock ran out
			fmt.Println("\n[KILL SWITCH ACTIVATED ] 300ms reached. Shutting down system!")
			return
		}
	}

}
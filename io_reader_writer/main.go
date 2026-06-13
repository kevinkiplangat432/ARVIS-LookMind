package main

// io reader and writer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)
// define a main function

func main() {
	// create a large mock file for testing 
	filename := "largedataset.json" 
	mockdata := `[

			{"id": 1, "name": "alice", "role" : "admin"},
			{"id": 2, "name" : "bob", "role" : "user"},
			{"id": 2, "name" : "charlie", "role" : "moderator"},	
			{"id": 2, "name" : "David", "role" : "user"}	

	]`

	// os write file ...will write the data into my file and [] byte will convert it into a stream of bytes.
	err := os.WriteFile(filename, []byte(mockdata), 0644) // the 0644 sets the file permissions
	if err != nil {
		panic(err) // this stops the normal execution of my program 
	}
	defer os.Remove(filename) // clean up the filename when the program ends

	// open the file as an io.reader stream 
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close() // close the file when the program ends

	fmt.Println("stream file token-by token using bufio....")

	// wrap the file in a buffered reader
	reader := bufio.NewReader(file)
	
	for {
		// read the stream up until each closing bracket
		line, err := reader.ReadString('}')

		// clean up formatting artifacts like commas,brackets and newlines
		cleaned := strings.TrimSpace(line)
		cleaned = strings.TrimPrefix(cleaned, "[")
		cleaned = strings.TrimPrefix(cleaned, ",")
		cleaned = strings.TrimSpace(cleaned)

		//print the token if it contains actual Json data
		if cleaned != "" && cleaned != "]" {
			fmt.Printf("Parsed Object: %s\n", cleaned)

		}

		// check for errors after processing what was read
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Reached the end of the stream: %v \n", err)
				break
			}
			fmt.Printf("Error reading stream: %v\n", err)
			break
		}
	}

}
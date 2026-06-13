package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// define a struct to bind the stream data into

type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {
	// setup our mock streaming file

	filename := "stream_test.json"
	mockData := `[
		{"id": 1, "name": "alice", "role": "admin"},
		{"id": 2, "name": "bob", "role": "user"},
		{"id": 3, "name": "charlie", "role": "moderator"}	
	]`


	_ = os.WriteFile(filename, []byte(mockData), 0644)
	defer os.Remove(filename)

	file, _ := os.Open(filename)
	defer file.Close()

	//initialize the streaming decoder directly on the file reader
	decoder := json.NewDecoder(file)

	//read the opening array brackect "["
	t, err := decoder.Token() // 
	if err != nil {
		panic(err)
	}
	fmt.Printf("started parsing stream opening token %v \n", t)
	//loop through the array elements

	for decoder.More() {
		var u User

		//decode a string object from the stream into our struct on the fly
		err := decoder.Decode(&u)
		if err != nil {
			fmt.Printf("Decode error:%v\n", err)
			break
		}

		fmt.Printf("Decode struct into memory -> Id: %d | Name: %s |Role: %s\n ", u.ID, u.Name, u.Role)
	}

	// read the closing array brackect ']'
	t, err = decoder.Token()
	if err != nil && err !=io.EOF{
		panic (err)
	}

	fmt.Printf("Finished parsing stream. closing token: %v\n", t)
}
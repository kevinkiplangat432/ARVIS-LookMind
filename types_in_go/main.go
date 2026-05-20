package main

import(
	"fmt"
)

type User struct {
	ID int
	Username string
	IsActive bool
}

func main(){
	// cteating a struct instance using a literal struct definition
	u := User{
		ID: 101,
		Username: "john_doe",
		IsActive: true,
	}
	// Accessing fields using dot notation
	fmt.Printf("user %d: %s (Active: %t)\n", u.ID, u.Username, u.IsActive)
}


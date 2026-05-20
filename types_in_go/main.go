package main

import(
	"fmt"
)
var count int       // 0
var balance float64 // 0.0
var active bool     // false
var name string  

// just because Go does not want me to set up stuff and not use them actually

func print_stuff(){
	fmt.Printf("count: %d\n", count)
	fmt.Printf("balance: %.2f\n", balance)
	fmt.Printf("active: %t\n", active)
	fmt.Printf("name: '%s'\n", name)
}

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
	print_stuff()
}


// (u User) is a value receiver

func (u User) ProfileMessage() string {
	return "User: " + u.Username
}

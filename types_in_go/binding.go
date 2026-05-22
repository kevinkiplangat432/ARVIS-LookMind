package main
import (
	// "fmt"
)

type Userv2 struct{
	ID int
	Username string

}

// value Receiver 
func (u Userv2) Greet() string{
	return "Hello, " + u.Username
}

// pointer receiver can modify the original struct 
func (u *Userv2) Deactivate() {
	// modify the actual struct instace directly
	u.Username = "deactivated"
}

// func main() {
// 	u := Userv2{ID: 101,
// 	Username: "john_doe",
// }

// 	fmt.Println(u.Greet())
// }
package main

import "fmt"

// define an interface job describtion
type Speaker interface{
	Speak() string
}

// struct 1
type Human struct {
	Name string
	
}


// the reciver
func (h Human) Speak() string {
	return "Hello! my name is " + h.Name
}


// struct b
type Dog struct {
	Breed string
}

// the method for dog 
func (d Dog) Speak() string{
	return "Woof! I am a " + d.Breed

}

// the interface User ( the flexible structure)


func MakeItSpeak(s Speaker) {
	fmt.Println("The speaker says: ", s.Speak())
}

func main(){
	// instantiate
	person := Human{
		Name: "Alice",
	}
	puppy := Dog{
		Breed: "Golden Retriever",
	}

	MakeItSpeak(person)
	MakeItSpeak(puppy)
}


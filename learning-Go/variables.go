package main

import (
	"fmt"
)

func studentInfo() {

	var student1 string = "John" //type string
	var student2 = "jane" //type inferred as string
	x := 10 //type inferred as int
	y := 3.14 //type inferred as float64
	isStudent := true //type inferred as bool

	println(student1)
	println(student2)
	println(x)
	println(y)
	println(isStudent)
}


//variable declaratrion withouth initialization
func var_no_decla(){
	var a int // default value is 0
	var b string
	var c bool


	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}

func assign_them_now(){
	var student1 string
	student1 = "John"
	fmt.Println(student1)
}


func multiplefunctionDeclaration(){
	var a, b, c, d int = 1,2,3,4

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)

}


func withoutTheTypeKeyword() {
  var a, b = 6, "Hello"
  c, d := 7, "World!"

  fmt.Println(a)
  fmt.Println(b)
  fmt.Println(c)
  fmt.Println(d)
}



func declarationInblock(){
	var (
		a int
		b int =1
		c string = "hello"
		d bool = true
	)

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}









//variables
/*
in go three types of variables are allowed
int
float32
string
bool
*/

// declare
/*
there are  two ways "var" ":="
:= can onlybe used inside the fuctipon
var can be used outside the function as well as inside the func
:= can not be declared and assigned differently
*/ 
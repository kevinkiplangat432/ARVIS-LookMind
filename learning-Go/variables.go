package main

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
*/ 


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

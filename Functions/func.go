package main

import (
	"os"
	"log/slog"
)

type Mathop func(int, int) int 

// first class functions and closures 
// write the addition func
func add(a int,b int)int {
	return a + b
}

// write the multiply function
func multiply(a int, b int) int {
	return a * b
}

// write a function that takes the values and the action(add or multiply)
func compute(a int, b int, operation Mathop) int {
	return operation(a, b)

}
// closures
func newcounter() func() int{
	count := 0  // outer var

	// now return an anonymous func This is closure
	// It "traps the "count" 
	return func() int{
		count = count + 1
		return count
	}

}

//defer

func writeLog(){
	slog.Info("starting the function")

	// open the a file 
	file, err := os.Create("log.txt")
	if err != nil {
		return
	}


	// now if we close the file now it will not work below this (code)
	// defer the close 
	defer file.Close()

	slog.Info("writing into the log.txt file....")
	file.WriteString("hello Go")


	slog.Info("reached the end of the function code")

	// this is where file.Close() runs.
}



//panic and recovery 



func main(){
	// pass add and multiply functions as data into compute
	sum := compute(5,3, add)
	product := compute(5,3, multiply)

	slog.Info("Calculation results", "sum",sum, "product",product)


	counterA := newcounter()

	slog.Info("starting counter","count", counterA())

	writeLog()

	
}
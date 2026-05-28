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
	defer func(){
		if err := file.Sync(); err != nil {
			slog.Error("Failed to sync file to disk", "error", err)
		}
		file.Close()
	}()

	slog.Info("writing into the log.txt file....")
	_, err = file.WriteString("hello Go")
	if err != nil {
		slog.Error("failed to write to file", "error", err)
		return 
	}
	


	slog.Info("reached the end of the function code")

	// this is where file.Close() runs.
}



//panic and recovery 
func safeDivide(a int, b int) int {
	// we must defer the recovery function before the panic happens
	defer func(){
		// recover intercepts the point state
		if r := recover(); r != nil {
			slog.Error("saved the program from crashing", "panic_reason", r)
		}
	}()
	slog.Info("Attempting divison", "a", a, "b", b)
	return a/b
}


func main(){
	// pass add and multiply functions as data into compute
	sum := compute(5,3, add)
	product := compute(5,3, multiply)

	slog.Info("Calculation results", "sum",sum, "product",product)


	counterA := newcounter()

	slog.Info("starting counter","count", counterA())

	writeLog()
	// Noraml execution
	result := safeDivide(10,2)
	slog.Info("Division success", "result", result)

	//this will trigger a panic but our code will recover!
	failedResult := safeDivide(10,0)
	// the program run

	slog.Info("Program completed all task successfully!", "final check", failedResult )
}
/*
go has three functions to output text
Print()
Println()
Printf()
*/

package main
import ("fmt")

func printfunc(){
	var i, j string = "hello", "world"
	fmt.Print(i)
	fmt.Print(j)
	//  new lines
	fmt.Print(i, "\n")
	fmt.Print(j, "\n")
	// multiple 
	fmt.Print(i,"\n", j)
	// adding a space
	fmt.Print(i, " ", j )
}

func println_func(){
	var i,j string = "Hello","World"
	fmt.Println(i,j)
}

func print_f_func(){
	var i string = "Hello"
    var j int = 15

    fmt.Printf("i has value: %v and type: %T\n", i, i)
    fmt.Printf("j has value: %v and type: %T", j, j)
}



/*
%v	Prints the value in the default format
%#v	Prints the value in Go-syntax format
%T	Prints the type of the value
%%	Prints the % sign
*/


func formatingverbs() {
  var i = 15.5
  var txt = "Hello World!"

  fmt.Printf("%v\n", i)
  fmt.Printf("%#v\n", i)
  fmt.Printf("%v%%\n", i)
  fmt.Printf("%T\n", i)

  fmt.Printf("%v\n", txt)
  fmt.Printf("%#v\n", txt)
  fmt.Printf("%T\n", txt)
}


/*
%b	Base 2
%d	Base 10
%+d	Base 10 and always show sign
%o	Base 8
%O	Base 8, with leading 0o
%x	Base 16, lowercase
%X	Base 16, uppercase
%#x	Base 16, with leading 0x
%4d	Pad with spaces (width 4, right justified)
%-4d	Pad with spaces (width 4, left justified)
%04d	Pad with zeroes (width 4
*/



func interger_formating() {
  var i = 15
 
  fmt.Printf("%b\n", i)
  fmt.Printf("%d\n", i)
  fmt.Printf("%+d\n", i)
  fmt.Printf("%o\n", i)
  fmt.Printf("%O\n", i)
  fmt.Printf("%x\n", i)
  fmt.Printf("%X\n", i)
  fmt.Printf("%#x\n", i)
  fmt.Printf("%4d\n", i)
  fmt.Printf("%-4d\n", i)
  fmt.Printf("%04d\n", i)
}


//string formating
/*
%s	Prints the value as plain string
%q	Prints the value as a double-quoted string
%8s	Prints the value as plain string (width 8, right justified)
%-8s	Prints the value as plain string (width 8, left justified)
%x	Prints the value as hex dump of byte values
% x	Prints the value as hex dump with spaces
*/

func stringFormattingVerbs() {
  var txt = "Hello"
 
  fmt.Printf("%s\n", txt)
  fmt.Printf("%q\n", txt)
  fmt.Printf("%8s\n", txt)
  fmt.Printf("%-8s\n", txt)
  fmt.Printf("%x\n", txt)
  fmt.Printf("% x\n", txt)
}

// bool formating

// %t  value of the boolean operator in true or false format (same as using %v)


func boolformating() {
  var i = true
  var j = false

  fmt.Printf("%t\n", i)
  fmt.Printf("%t\n", j)
}

// float formatting 
/*
%e	Scientific notation with 'e' as exponent
%f	Decimal point, no exponent
%.2f	Default width, precision 2
%6.2f	Width 6, precision 2
%g	Exponent as needed, only necessary digits
*/


func floatformatting() {
  var i = 3.141

  fmt.Printf("%e\n", i)
  fmt.Printf("%f\n", i)
  fmt.Printf("%.2f\n", i)
  fmt.Printf("%6.2f\n", i)
  fmt.Printf("%g\n", i)
}
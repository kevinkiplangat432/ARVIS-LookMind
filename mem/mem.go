package main

import "log/slog"

func main() {
	slog.Info(" sclice capacity demo")
	// create a slice with a length of 0 but the capacity of 3
	scores := make([]int, 0, 3)
	slog.Info("Initial State", "len", len(scores), "cap", cap(scores))

	scores = append(scores, 90, 85, 95)
	slog.Info("filled capacity", "len", len(scores), "cap", cap(scores))

	// this append will trigger a memory reallocation under the hood
	scores = append(scores, 100)
	slog.Info("Exceeded capacity(Doubled!)", "len", len(scores), "cap", cap(scores))

	slog.Info("MAP GOTCHAS DEMO")
	// creating an initialized map
	ages := make(map[string]int)
	ages["Alice"]= 25
	ages["Bob"]= 30


	// safe deletion
	delete(ages, "Bob")

	// checking if the key exist
	val, exists := ages["Bob"]
	slog.Info("Checking deletion", "exists", exists, "value", val)


	slog.Info("POINTERS & COPYING DEMO")

	num := 10
	slog.Info("Before modifyVal", "num", num)
	modifyVal(num)
	slog.Info("After modifyVal (No change, it copied data)", "num", num)

	slog.Info("Before modifyPointer", "num", num)
	modifyPointer(&num) // Passing the memory address explicitly
	slog.Info("After modifyPointer (Changed, we went to the address)", "num", num)


	subSlice()
	AppendAfterSUbSlice()
}
// Receives a copy of the actual number value
func modifyVal(val int) {
	val = 999 
}

// Receives a copy of the memory address pointing to the number
func modifyPointer(valPtr *int) {
	*valPtr = 999 // De-referencing: alter the value at that address
}

func subSlice() {
	original := []int{10, 20, 30, 40}
	
	// Create a sub-slice (indexes 1 and 2)
	subSlice := original[1:3] // contains [20, 30]

	// Changing the sub-slice mutates the original slice!
	subSlice[0] = 999 

	slog.Info("Mutating sub-slice changes original!", "original", original)
	// Output will show original is now: [10, 999, 30, 40]
}


func AppendAfterSUbSlice() {
	// Capacity is 4. Array is full.
	a := []int{1, 2, 3, 4} 
	sub := a[0:2] // sub is [1, 2], len=2, cap=4 (shares array with 'a')

	// There is room in the capacity! This overwrites the 3 in 'a'
	sub = append(sub, 99) 
	slog.Info("Shared array append", "a", a) // 'a' becomes [1, 2, 99, 4]

	// This append exceeds capacity! 'sub' reallocates to a brand new array.
	sub = append(sub, 88, 77) 
	
	// Now modifying 'sub' does NOT affect 'a' anymore!
	sub[0] = 555
	slog.Info("Snapped link", "a", a, "sub", sub)
}

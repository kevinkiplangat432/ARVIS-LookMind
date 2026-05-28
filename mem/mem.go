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


	slog.Info("--- 3. POINTERS & COPYING DEMO ---")

	num := 10
	slog.Info("Before modifyVal", "num", num)
	modifyVal(num)
	slog.Info("After modifyVal (No change, it copied data)", "num", num)

	slog.Info("Before modifyPointer", "num", num)
	modifyPointer(&num) // Passing the memory address explicitly
	slog.Info("After modifyPointer (Changed, we went to the address)", "num", num)
}
// Receives a copy of the actual number value
func modifyVal(val int) {
	val = 999 
}

// Receives a copy of the memory address pointing to the number
func modifyPointer(valPtr *int) {
	*valPtr = 999 // De-referencing: alter the value at that address
}
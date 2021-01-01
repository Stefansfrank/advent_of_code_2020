package main

import (
	"fmt"
	"time"
)

func elvGame(startSeq []int, terminate int) (called int) {

	// I am keeping a log for the last occurence in form of a map
	log := make(map[int]int)

	// Adding the start sequence without the last number
	for ix := 0; ix < len(startSeq) - 1; ix++ {
		log[startSeq[ix]] = ix
	}

	// The last number of the start sequence starts the loop
	nextCall := startSeq[len(startSeq)-1]
	for ix := len(startSeq)-1; ix < terminate; ix++ {

		// this is the main 'trick' i.e. computing the next call
		// before logging the current call
		called            = nextCall
		position, exists := log[called]
		if exists {
			nextCall = ix - position
		} else {
			nextCall = 0
		}
		log[called]       = ix
}

	return
}


// MAIN ----
func main () {

	start := time.Now()
	seq   := []int{9,19,1,6,0,5,4}

	// Part 1
	fmt.Printf("Last number called after 2020: %v\n", elvGame(seq, 2020))
 	fmt.Printf("Execution time: %v\n", time.Since(start))
	start = time.Now()

	// Part 2
	fmt.Printf("Last number called after 30000000: %v\n", elvGame(seq, 30000000))
 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
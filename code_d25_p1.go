package main

import (
	"fmt"
	"time"
)

// the transformation with a given loop number
// returns the result
func transform(sub, num int) int {
	result := 1
	for i := 0; i < num; i++ {
		result *= sub 
		result = result % 20201227
	}
	return result
}

// running the transformation until a match is found
// returns the loop number
func transformMatch(sub, match int) int {
	result := 1
	ll     := 0
	for {
		ll++
		result *= sub 
		result = result % 20201227
		if result == match {
			return ll
		}
	}
	return -1
}

func main() {

	start   := time.Now()

	//pbKey := []int{5764801, 17807724} // test data
	pbKey := []int{11562782, 18108497} // my data

	// Part 1
	lp0 := transformMatch(7, pbKey[0])
	lp1 := transformMatch(7, pbKey[1])
	fmt.Println("Loop Card: ", lp0)
	fmt.Println("Loop Lock: ", lp1)
	fmt.Println("Enc Key: ", transform(pbKey[0], lp1))
	fmt.Println("Enc Key (validation): ", transform(pbKey[1], lp0))

	fmt.Printf("Execution time: %v\n", time.Since(start))

}
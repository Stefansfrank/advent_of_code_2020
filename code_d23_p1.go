package main

import (
	"fmt"
	"time"
)

var chn  []int  // contains the cups with index = label of cup, value = label of next cup
var chnLen int  // 9 for part 1 and 1,000,000 for part 2

// Initializes the chain from the 9 digit seed given
// returns first value of seed as current
func initChain(vals []int) int {
	chn = make([]int, chnLen+1)
	for i := 0; i < 8; i++ {
		chn[vals[i]] = vals[i+1]
	}
	if chnLen > 9 {
		chn[vals[8]] = 10		
		for i := 10; i < chnLen; i++ {
			chn[i] = i+1
		}
		chn[chnLen] = vals[0]
	} else {
		chn[vals[8]] = vals[0]
	}

	return vals[0]
}

// These are helper functions for readability 
// could be done inline if compact code is desired
func nxt(i int) int {
	return chn[i]
}
func nxt2(i int) int {
	return chn[chn[i]]
}
func nxt4(i int) int {
	return chn[chn[chn[chn[i]]]]
}

// One game move
func move(cur int) int {
	rem      := nxt(cur)  // first num of parked chain of three cups
	chn[cur]  = nxt4(cur) // link the 4th cup after to current to close remaining chain 

	// search for valid destination
	dvl := (cur + chnLen - 2) % chnLen + 1                    // equivalent to dvl = cur - 1 but cyclic and starting with 1
	for dvl  == rem || dvl == nxt(rem) || dvl == nxt2(rem) {  // check whether destination is among the 3 parked cups
		dvl = (dvl + chnLen - 2) % chnLen + 1                 // equivalent to dvl = dvl - 1 but cyclic and starting with 1
	}

	// break at destination and insert parked chain
	brk            := nxt(dvl)
	chn[dvl]        = rem
	chn[nxt2(rem)]  = brk

	return nxt(cur)
}

func main() {

	start   := time.Now()

	//inp := []int{3,8,9,1,2,5,4,6,7} // test
	inp := []int{9,4,2,3,8,7,6,1,5} // real

	// Part 1
	chnLen  = 9
	cur    := initChain(inp)

	// loop
	for i := 0; i < 100; i++ {		
	 	cur = move(cur)
	}

	// result computation
	cur     = 1 
	result := 0
	for i := 0; i < 8; i++ {
		result *= 10
		result += nxt(cur)
		cur     = nxt(cur)
	}
	fmt.Println("Result Part 1: ", result)

	// Part 2
	chnLen  = 1000000
	cur     = initChain(inp)

	// loop
	for i := 0; i < 10000000; i++ {		
	 	cur = move(cur)
	}

	result = chn[1]*chn[chn[1]]
	fmt.Println("Result Part 2: ", result)


	fmt.Printf("Execution time: %v\n", time.Since(start))

}
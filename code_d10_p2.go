package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"sort"
	"math"
//	"strings"
)

type bagDef struct {
	name    string
	content []bagRef
}

type bagRef struct {
	amt  int
	name string
}

// no error handling ...
func readTxtFile2Int (name string) (lines []int) {
	
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error on opening file: %v\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		lines = append(lines, atoi(scanner.Text()))
	}

	return
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// Part 1: count the jolt differences once all adapters are sorted by jolt
// the result ns[i] contains how many steps of i jolts are present in that chain
func countSteps (jolts []int) (ns []int) {
	sort.Ints(jolts)
	ns  = make([]int, 10)
	cur := 0
	for _, jolt := range jolts {
		ns[jolt - cur]++
		cur = jolt
	}
	ns[3]++
	return
}

// Part 2: the analysis in part 1 shows that the only differences occuring are 1 and 3 jolt
// whenever a 3 jolt difference is in that chain of all adapters, the adapters on both ends have to be used
// in all valid solutions as omitting these would create a difference of 4 jolts to be bridged.
// 
// This function counts how often 1-jolt steps occur in a row before an adapter that is part of a 3-jolt step is encountered.
// Only adapters that are part of 1-jolt steps can be permutated and create multiple solution. The permutations are:
//  - if there is only one 1 1-jolt step in a row, both adapters are part of a 3 jolt step on either side and they are both needed
//  - if there are 2 1-jolt steps in a row the middle adapter can be left out or be in so every 2 jolt sequence creates two solutions
//  - if there are 3 1-jolt steps either of the two middle adapters can be present or not creating 4 solutions per sequence
//  - if there are 4 1-jolt steps all three middle adapters can be present or not (8 solutions) with the EXCEPTION that not all three can be missing 
//      -> 7 solutions per sequence
// 
// The result ns[i] contains how often a sequence of i 1 jolt adapters in a row is part of the chain
// it turns out that the maximum number of 1 jolt steps in a row is 4 so the number of all solutions is:
//  2 to the power of (ns[2]) * 4 to the power of (ns[3]) * 7 to the power of (ns[4]) 
func countOneSeqs (jolts []int) (ns []int) {
	sort.Ints(jolts)
	jolts = append(jolts, jolts[len(jolts)-1]+3)
	ns  = make([]int, 10)
	cur := 0
	len := 0
	for _, jolt := range jolts {
		if jolt-cur == 1 {
			len++
		} else {
			ns[len]++
			len = 0
		}
		cur = jolt
	}
	return
}

// MAIN ----
func main () {

	start   := time.Now()

	jolts  := readTxtFile2Int("input_d10_p1.txt")

	// Part 1
	ns := countSteps(jolts)
	fmt.Printf("Part 1: %v\n", ns[1]*ns[3])

	// Part 2
	ns = countOneSeqs(jolts)
	fmt.Printf("Part 2 number of solutions: %v\n", int(math.Pow(2,float64(ns[2])) * 
		math.Pow(4,float64(ns[3])) * math.Pow(7,float64(ns[4])))) // see above for explanation

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
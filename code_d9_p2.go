package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
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
	file, _ := os.Open(name)
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

// Part 1: find first number that is the sum of two within the last 25
// While looping through the range of numbers after the preamble,
// I keep a map [i]n that counts how often each number i occurs in the previous 'pre'(25) numbers 
// at every step, I reduce that count for the number that drops out of range and
// increase it for each new number that enters that range 
func findFirstNumber(cypher []int, pre int) int {

	cyLen := len(cypher)
	preMap := make(map[int]int)

	for i:=0; i<pre; i++ {
		preMap[cypher[i]] += 1
	}

	for j:=pre; j<cyLen; j++ {

		found := false
		for i:=j-pre; i<j; i++ {
			if preMap[cypher[j]-cypher[i]] > 0 {
				found = true
				break
			}
		}

		if !found {
			return cypher[j]
		}

		preMap[cypher[j-pre]] -= 1
		preMap[cypher[j]] += 1
	}

	return -1
}

// Part 2: find continous range summed up to be the result of part 1
// while looping through the range of numbers, I keep a sum over a range of previous numbers that is smaller than the target.
// Each iteration, I add the next number. If the sum becomes greater than the target, I take away numbers from
// the lower end of the range until the sum is smaller again. This is done until adding the next number is hitting the target.
func findRanges(cypher []int, target int) {

	cyLen := len(cypher)
	cypher = append(cypher, target+1) // avoid out of bounds on last iteration

	from := 0
	to   := 1
	sum  := cypher[0]+cypher[1]

	for to < cyLen {

		if sum == target {
			fmt.Printf("Range for part 2 found: Ix%v(%v) to Ix%v(%v): %v\n",
				 from, cypher[from], to, cypher[to], minMaxSumRange(cypher, from, to))
			sum -= cypher[from]
			from++
			continue
		}

		if sum < target {
			to++
			sum += cypher[to]	
			continue		
		}

		if sum > target {
			sum -= cypher[from]
			from++
			if from == to {
				to++
				sum += cypher[to]	
			}		
		}
	}

}

// only used to compute the resulting checksum for part 2
func minMaxSumRange(cypher []int, from, to int) int {
	min := cypher[from]
	max := cypher[from]
	for i:= from+1; i <= to; i++ {
		if cypher[i] > max {
			max = cypher[i]
		} else {
			if cypher[i] < min {
				min = cypher[i]
			}
		}
	}
	return max + min
}


// MAIN ----
func main () {

	start   := time.Now()

	cypher  := readTxtFile2Int("input_d9_p1.txt")

	// part 1
	weakNr  := findFirstNumber(cypher, 25)
	fmt.Println("First weak number (part 1):", weakNr)
	
	// part 2
	findRanges(cypher, weakNr)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
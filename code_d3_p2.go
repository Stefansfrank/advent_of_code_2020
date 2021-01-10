package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile (name string) (lines []string) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}
	return
}

// parse map into a 2 dimensional map of 0s and 1s in order to have a simple sum later
func toBitMap (cMap []string) (bMap [][]int) {
	bMap = [][]int{}

	// loop through map lines
	for _, mLine := range cMap {
		bMapLine := []int{}

		// loop through individual line
		for _, cc := range mLine {			
			if cc == '#' {
				bMapLine = append(bMapLine, 1)
			} else {
				bMapLine = append(bMapLine, 0)
			}
		}

		bMap = append(bMap, bMapLine)
	}

	return
}

// traverse map and count the values in the cells as parsed
func traverseMap(bMap [][]int, dx, dy int) (trees int) {

	wid    := len(bMap[0])
	hgt    := len(bMap)
	x      := 0
	trees   = 0	

	for y := 0; y < hgt; y += dy {
		trees += bMap[y][x]
		x = (x + dx) % wid
	}

	return
}

// MAIN ----
func main () {

	start := time.Now()

	cMap := readTxtFile("input_d3_p1.txt")
	bMap := toBitMap(cMap)

	trees := traverseMap(bMap,3,1)
	fmt.Printf("Trees for Part 1: %v\n", trees)

	trees *= traverseMap(bMap,1,1)
	trees *= traverseMap(bMap,5,1)
	trees *= traverseMap(bMap,7,1)
	trees *= traverseMap(bMap,1,2)
	fmt.Printf("Trees for Part 2: %v\n", trees)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
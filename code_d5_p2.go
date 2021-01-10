package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
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

// convert boarding pass by using binary conversion of the rules
func convertBp(bp string) (int) {

	// row in binary
	rowBin  := strings.ReplaceAll(strings.ReplaceAll(bp[0:7],"B","1"),"F","0")
	seatBin := strings.ReplaceAll(strings.ReplaceAll(bp[7:10],"R","1"),"L","0")
	rowInt, _   := strconv.ParseUint(rowBin, 2, 7)
	seatInt, _  := strconv.ParseUint(seatBin, 2, 3)
	return int(rowInt*8) + int(seatInt)

}

// find seat between two occupied and detect max seat on the way
func findSeat(brdDoc []string) (int, int) {

	seatOcc := make(map[int]bool) 
	max := 0

	// detect max and fill occupied map
	for _,bp := range brdDoc {
		cur := convertBp(bp)
		seatOcc[cur] = true
		if cur > max {
			max = cur
		}
	}

	// find unoccupied seat
	for i:=1; i<max-1; i++ {
		if seatOcc[i-1] && !seatOcc[i] && seatOcc[i+1] {
			return max, i
		}
	}

	return max, -1
}

// MAIN ----
func main () {

	start  := time.Now()

	brdDoc  := readTxtFile("input_d5_p1.txt")
	max, seat := findSeat(brdDoc)

	fmt.Printf("Highest Seat: %v\n", max)
	fmt.Printf("Your Seat: %v\n", seat)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
	"strings"
	"math/big"
)

// generic text file reader (should be standard in go ...)
func readTxtFile (name string) (lines []string) {
	
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error on opening file: %v\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}

	return
}

// Buses helper struct (big.Int for part 2)
type line struct {
	bigNum  *big.Int  // the number of the line
	ix      *big.Int  // the index in the input list
	match   bool      // whether the line was matched (part 2 only)
} 

// Parses bus string 
// creating the line sturctures
func parseBuses(sheet []string) (now int64, buses []line) {
	eol      := false
	num      := -1
	var err error

	// line 1
	tmp, _ := strconv.Atoi(sheet[0])
	now     = int64(tmp)

	// line 2
	buses = []line{}
	ix    := int64(0)
	for !eol {
		// parsing the input
		cIx     := strings.IndexRune(sheet[1],',')
		if cIx == -1 {  
			num, err = strconv.Atoi(sheet[1])
			eol = true
		} else {
			num, err = strconv.Atoi(sheet[1][:cIx])
			sheet[1] = sheet[1][cIx+1:]
		}
		if err == nil {
			buses = append(buses, line{ 
				bigNum: big.NewInt(int64(num)),
				ix:     big.NewInt(ix),
				match:  false })
		}
		ix++
	}
	return
}

// PART 1
// simple computation using modulo
func nextBus(now int64, buses []line) {
	minTime := int64(10000)
	minLine := int64(-1)

	for _,bus := range buses {
		num  := bus.bigNum.Int64()
		next := num - (now % num)
		if next < minTime {
			minTime = next
			minLine = num
		}
	}

	fmt.Printf("Next bus is the line %v at %v i.e. in %v minutes. Checksum: %v\n", 
		minLine, now+minTime, minTime, minTime * minLine)
}

// PART 2
// The basic idea for part 2 is to find a match for a subset of buses and then continue
// in steps of the LCM of all matched buses to try to match the next
//
// this uses the big/int libraries for two reasons:
// - protect from overflow
// - use built in GCT function in order to get LCM
func findTime(buses []line) (big.Int) {

	// temporary big Int variables needing initialization
	gcd       := big.NewInt(0)
	bmod      := big.NewInt(0)
	bigZero   := big.NewInt(0)

	// NOTE: assumes the first line in the input not to be 'x' !!!
	curLcm    := buses[0].bigNum      // the LCM of all matched buses
	bestTime  := big.NewInt(0)        // the time that works for all matched buses
	numMatch  := 1                    // the amount of matched buses
	buses[0].match = true
	fmt.Printf("First Bus %v\n", buses[0].bigNum.String())

	// loop through time in steps of the LCM of matched buses
	done := false
	for !done {
		bestTime.Add(bestTime, curLcm)

		for i, bus := range buses {
			// only unmatched buses
			if !bus.match {

				// compute the relative arrival time
				// and compare with the desired delay
				bmod.Mod(bmod.Add(bestTime,bus.ix), bus.bigNum)
				if bmod.Cmp(bigZero) == 0 {

					// we got a match!
					buses[i].match = true
					numMatch++

					if numMatch == len(buses) {
						done = true
					} else {
						// compute new LCM using LCM(a,b) = abs(a*b)/GCD(a,b)
						// without abs() as we deal with positive numbers
						gcd.GCD(nil, nil, curLcm, bus.bigNum)
						curLcm.Mul(curLcm, bus.bigNum)
						curLcm.Div(curLcm, gcd)
					}

					fmt.Printf("Match for bus %v at %v mins - Next LCM:%v\n",
					 bus.bigNum.String(), bestTime.String(), curLcm.String())
				}
			}
		}
	}

	return *bestTime
}

// MAIN ----
func main () {

	start 		:= time.Now()
	fmt.Println("Part 1")

	// needed for both tasks ...
	sheet  		:= readTxtFile("input_d13_p1.txt")
	now, buses  := parseBuses(sheet)

	// part 1
	nextBus(now, buses)

	fmt.Printf("Execution time: %v\n\nPart 2\n", time.Since(start))
	start 		 = time.Now()

	// part 2
	best        := findTime(buses)
	fmt.Printf("Time everything lines up: %v\n", best.String())

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"strings"
	"time"
)

// no error handling ...
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

// short inline capable / no error handling Atoi / Itoa
func atoi(x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

var patterns []string         // the patterns parsed from the input
var rawRules map[int][][]int  // maps rule number to the rule numbers it builds on (once resolved into the vldRules below, this will be empty)
var vldRules map[int][]string // maps rule number to the set of patterns it covers

// building rawRules and seed vldRules with the leaf entries
func parseSheet(sheet []string) {
	patterns = []string{}
	vldRules = make(map[int][]string)
	rawRules = make(map[int][][]int)

	ruleSec  := true
	for i,line := range sheet {

		if len(line) == 0 {
			ruleSec = false
			continue
		}

		// parse rules
		if ruleSec {

			// parse rule number
			var tmp []int
			ix   := strings.IndexByte(line, ':')
			num  := atoi(line[:ix])
			line  = line[ix+2:]

			// parse leaf entries straight into vldRules
			if line[0] == '"' {
				vldRules[num] = []string{ string(line[1]) }
				continue
			}

			// is there an '|'?
			ix = strings.IndexByte(line, '|')
			if ix == -1 {
				// single value
				ix2 := strings.IndexByte(line, ' ')
				if ix2 == -1 {
					tmp = []int{ atoi(line) }
				// two (or more values)
				} else {
					if line == "4 1 5" { // special casing the only case where three rules are indicated
						tmp = []int{ 4,1,5 }
					} else {				
						tmp = []int{ atoi(line[:ix2]), atoi(line[ix2+1:]) }
					}			
				}
				rawRules[num] = [][]int { tmp }

			// no or
			} else {
				ix2 := strings.IndexByte(line, ' ')
				if ix2 == ix-1 {
					tmp = []int{ atoi(line[:ix2]) }
				} else {
					tmp = []int{ atoi(line[:ix2]), atoi(line[ix2+1:ix-1]) }
				}
				rawRules[num] = [][]int { tmp }
				line = line[ix+2:]
				ix2  = strings.IndexByte(line, ' ')
				if ix2 == -1 {
					tmp = []int{ atoi(line) }
				} else {
					tmp = []int{ atoi(line[:ix2]), atoi(line[ix2+1:]) }
				}
				rawRules[num] = append(rawRules[num], tmp)
			}

		// pattern section
		} else {
			patterns = sheet[i:]
			break
		}
	}
}

// returns all combinations of all strings in a with all strings in b
func crossComb(a []string, b []string) (r []string) {
	r = []string{}
	for _, aa := range a {
		for _, bb := range b {
			r = append(r, aa + bb)
		}
	}
	return
}

// going through the dependencies in rawRules and creates vldRules
func resolveRules() {
	for len(rawRules) > 0 {
		for num, rule := range rawRules {

			canResolve := true

			// check whether rule can be resolved
			for _, rp := range rule {
				for _, ri := range rp {
					_, ok := vldRules[ri]
					if !ok {
						canResolve = false
					}
				}
			}

			if canResolve {
				vldRules[num] = []string{}
				// loop through or conditions
				for _, rp := range rule {

					// loop through parts of individual patterns
					tmp := []string {""}
					for _, ri := range rp {
						tmp = crossComb(tmp, vldRules[ri])
					}
					vldRules[num] = append(vldRules[num], tmp...)
				}

				delete(rawRules, num)
			}
		}
	}
}

// tests whether a given pattern obeys a rule
func fitsRule(patt string, rule int) bool {
	for _, rPatt := range vldRules[rule] {
		if rPatt == patt {
			return true
		}
	}
	return false
}

// MAIN ----
func main() {

	start  := time.Now()

	sheet := readTxtFile("input_d19_p1.txt")
	parseSheet(sheet)
	resolveRules()

	// part 1 just compares the given patterns with the covered ones in vldRules[0]
	numValidPatts := 0
	for _, patt := range patterns {
		for _, zPatt := range vldRules[0] {
			if patt == zPatt {
				numValidPatts++
				break
			}			
		}
	}
	fmt.Println("Number of valid patterns in Part 1:", numValidPatts)

	// Part 2 is solve quite "manually" looking at the specifics of the instructions and patterns at hand
	// The added rule logic for 8 and 11 is not added to vldRules[] but the whole match logic is implemented
	// in code for rules 0,8 and 11
	//
	// The important conclusions after looking at the patterns and rules are:
	// - Rule 0 is made up from a pattern obeying rule 8 combined with a pattern obeying rule 11
	// - Rule 8 is now dynamic as it can be an arbitrary amount of patterns obeying rule 42
	// - Rule 11 is also dynamic defined recursively in a way that it consists of 
	//   an arbitrary number of patterns obeying 42 followed by the SAME amount of patterns obeying rule 31
	// - Patterns obeying either rule 42 or rule 31 are all 8 characters long (variable "sub" below)
	// - The longest pattern found in the test set is 96 (12*8) characters long
	//
	// Thus the solution below goes through all patterns, throws out patterns with a length not divisible by 8 and shorter than 3*8
	// and then goes through the 8 char subpatterns of valid patterns checking the following:
	// - at least the first half of subpatterns must obey 42 rounding up for odd numbers of 8 char subpatterns
	// - the remaining subpatterns can obey 42 or 31 with the following exceptions:
	// - once a subpattern is encountered that obeys only 31 the remaining subpatterns need to obey 31
	// - the last subpattern must always obey 31

	numValidPatts = 0
	sub := 8
	
	patternloop:
	for _, patt := range patterns {
		ln    := len(patt)

		// must be divisible by sub and a minimum size of 3*sub
		if ln % sub != 0 || len(patt) < 3*sub {
			continue
		}

		// at least the first half of the patterns rounded up must obey 42 
		split := sub * (int((ln / sub) / 2) + 1)  
		for i:=0; i<split; i+=sub {
			if !fitsRule(patt[i:i+sub], 42) {
				continue patternloop
			}
		}

		// besides the last, the next ones can be either, but:
		// once the first pattern only obeying 31 is encountered only patterns from 31 count 
		enc31 := false 
		for i:=split; i<ln-8; i+=sub {
			ob42 := fitsRule(patt[i:i+sub], 42)
			ob31 := fitsRule(patt[i:i+sub], 31)
			if !ob42 && !ob31 {
				continue patternloop
			} else {
				if !ob31 && enc31 {
					continue patternloop
				} else if !ob42 {
					enc31 = true
				}
			}
		}

		// the last one needs to be 31
		if fitsRule(patt[ln-sub:], 31) {
			numValidPatts++
		}
	}
	fmt.Println("Number of valid patterns in Part 2:", numValidPatts)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
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

// creates a map with the labels of required fields returning 1 for these
func requiredFields() (fieldMap map[string]int) {
	fieldMap = make(map[string]int)
	fieldMap["byr"] = 1
	fieldMap["iyr"] = 1
	fieldMap["eyr"] = 1
	fieldMap["hgt"] = 1
	fieldMap["hcl"] = 1
	fieldMap["ecl"] = 1
	fieldMap["pid"] = 1
	return	
}

// returns one for a valid field
func validField(key, val string) int {
	return 1
}

// just count valid passports
func countValidPps(ppDoc []string) (numPP int) {

	fieldCnt := 0
	fieldMap := requiredFields()

	for _, ppLine := range ppDoc {

		scIx := strings.IndexRune(ppLine, ':')
		spIx := strings.IndexRune(ppLine, ' ')

		// empty line
		if scIx == -1 {
			fieldCnt = 0
			fieldMap = requiredFields()
			continue
		}

		for scIx >= 0 {

			// determine key/value
			if spIx == -1 {
				spIx = len(ppLine)
			}
			key := ppLine[scIx-3:scIx]
			val := ppLine[scIx+1:spIx]

			// determine valid field
			fieldCnt += fieldMap[key] * validField(key, val)

			// prevent double counting of fields
			fieldMap[key] = 0

			// cut off the parsed key
			if spIx == len(ppLine) {
				scIx = -1
			} else {
				ppLine = ppLine[spIx+1:]
				scIx = strings.IndexRune(ppLine, ':')
				spIx = strings.IndexRune(ppLine, ' ')
			}
		}

		// test if 7 valid fields are present
		if fieldCnt == 7 {
			numPP++
			fieldCnt = 0 // required as there could be another line with an optional CID field
		}
	}

	return
}

// MAIN ----
func main () {

	start  := time.Now()

	ppDoc  := readTxtFile("input_d4_p1.txt")

	fmt.Printf("%v valid Passports\n", countValidPps(ppDoc))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
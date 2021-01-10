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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
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

// creates a map for testing valid eye colors
func validEyeColor(ecl string) (valid bool) {
	ecls := map[string]int{
 	   "amb": 0,
 	   "blu": 0,
 	   "brn": 0,
 	   "gry": 0,
 	   "grn": 0,
 	   "hzl": 0,
 	   "oth": 0,
	}
	_, valid = ecls[ecl]
	return
}

// returns one for a valid field
func validField(key, val string) int {

	valid := false

	switch key {
	case "byr":
		year := atoi(val)
		valid = year > 1919 && year < 2003 
	case "iyr":
		year := atoi(val)
		valid = year > 2009 && year < 2021 
	case "eyr":
		year := atoi(val)
		valid = year > 2019 && year < 2031 
	case "hgt":
		unIx := strings.Index(val, "cm")
		if unIx > -1 {
			hgt := atoi(val[:unIx])
			valid = hgt > 149 && hgt <194 
			break
		}
		unIx = strings.Index(val, "in")
		if unIx > -1 {
			hgt := atoi(val[:unIx])
			valid = hgt > 58 && hgt <77 
			break
		}
	case "hcl":
		if val[:1] != "#" || len(val) != 7 {
			break
		}
		_, err := strconv.ParseInt(val[1:], 16, 48)
 		valid = err == nil
	case "ecl":
		valid = validEyeColor(val)
	case "pid":
		if len(val) != 9 {
			break
		}
		_, err := strconv.ParseInt(val, 10, 0)
 		valid = err == nil
 	}

 	//fmt.Printf("%v:%v - %v\n", key, val, valid)
	if valid {
		return 1
	}

	return 0
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
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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// tests passwords for validity
func pwdValid(pwd string) bool {

	hyPo := strings.IndexRune(pwd, '-')
	spPo := strings.IndexRune(pwd, ' ')
	scPo := strings.IndexRune(pwd, ':')

	from := atoi(pwd[:hyPo])
	to   := atoi(pwd[hyPo+1:spPo])
	char := pwd[spPo+1:scPo]
	pw   := pwd[scPo+2:]

	cnt  := strings.Count(pw, char)

    return (cnt >= from && cnt <= to)
}

// MAIN ----
func main () {

	start := time.Now()

	pwds := readTxtFile("input_d2_p1.txt")
	valid := 0

	for _, pwd := range pwds {
		if pwdValid(pwd) {
			valid++
		}
	}

	fmt.Printf("%v valid Passwords found.", valid)

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
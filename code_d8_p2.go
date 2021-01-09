package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
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

// detects loops in the code 
// the first return parameter contains acc
// the second return parameter returns whether a loot was found 
func detectLoop(code []string) (int, bool) {

	execLines := make(map[int]bool) // map tracing lines executed
	nextLine  := 0
	acc       := 0

	for !execLines[nextLine] {

		cmd := code[nextLine][:3]
		prm := code[nextLine][4:]
		execLines[nextLine] = true

		switch cmd {
		case "acc":
			acc += atoi(prm)
			nextLine++
		case "nop":
			nextLine++
		case "jmp":
			nextLine += atoi(prm)
		}

		if nextLine == len(code) {
			return acc, false
		}
	}

	return acc, true
}

// loop through the code trying to replace every JMP and NOP
func findValidCode(code []string) int {
	for i, line := range code {
		cmd := line[:3]
		switch cmd {
		case "jmp":
			code[i] = strings.ReplaceAll(line, "jmp", "nop")
			acc, loop := detectLoop(code)
			if loop {
				code[i] = strings.ReplaceAll(line, "nop", "jmp")
			} else {
				return acc
			}
		case "nop":
			code[i] = strings.ReplaceAll(line, "nop", "jmp")
			acc,loop := detectLoop(code)
			if loop {
				code[i] = strings.ReplaceAll(line, "jmp", "nop")
			} else {
				return acc
			}
		}
	}
	return -1
}

// MAIN ----
func main () {

	start := time.Now()

	code := readTxtFile("input_d8_p1.txt")

	// part 1 
	acc, _ := detectLoop(code)
	fmt.Printf("Code encounters loop at acc = %v\n", acc)

	// part 2
	acc = findValidCode(code)
	fmt.Printf("Acc after corrected code runs through: %v\n", acc)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
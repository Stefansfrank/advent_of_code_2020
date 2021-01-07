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
	file, _ := os.Open(name)
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
func itoa(i int) string {
	return strconv.Itoa(i)
}

// this takes care of the nested parenthesis and such and extracts
// paranthesis free substrings to be evaluated
func parseLine(line string, part2 bool) (result int) {
	opPar := strings.LastIndexByte(line, '(')
	for opPar > -1 {
		clPar := strings.IndexByte(line[opPar+1:],')') + opPar + 1
		line   = line[:opPar] + evalLine(line[opPar+1:clPar], part2) + line[clPar+1:]
		opPar  = strings.LastIndexByte(line, '(')
	}
	return atoi(evalLine(line, part2))
}

// evaluates by replacing each operation "a op b" by it's result "c"
// finally providing just one number
func evalLine(line string, part2 bool) string {
	
	// part 2 i.e. sum priority
	if part2 {
		return evalMult(evalSums(line))
	}
	
	// part 1 i.e. left-2-right execution
	mulIx := strings.IndexByte(line, '*')
	sumIx := strings.IndexByte(line, '+')
	if mulIx + sumIx < -1 {
		return line
	} 
	for true {
		spIx   := strings.IndexByte(line, ' ')
		result := atoi(line[:spIx])
		spIx2  := strings.IndexByte(line[spIx+3:], ' ')
		if spIx2 == -1 {
			spIx2 = len(line)
		} else {
			spIx2 += spIx + 3
		}
		if mulIx == -1 || (sumIx > -1 && mulIx > sumIx) {
			result += atoi(line[spIx+3:spIx2])
		} else {
			result *= atoi(line[spIx+3:spIx2])			
		}
		if spIx2 < len(line) {
			line = itoa(result) + line[spIx2:]
		} else {
			line = itoa(result)
			break
		}
		mulIx = strings.IndexByte(line, '*')
		sumIx = strings.IndexByte(line, '+')
	}
	return line
}

// evaluates sums within a paranthesis free string 
func evalSums(line string) (string) {
	sumIx := strings.IndexByte(line, '+')
	if sumIx == -1 {
		return line
	}
	for sumIx != -1 {
		spIx := strings.LastIndexByte(line[:sumIx-1], ' ')
		spIx2 := strings.IndexByte(line[sumIx+2:], ' ') 
		if spIx2 == -1 {
			spIx2 = len(line)
		} else {
			spIx2 += sumIx + 2
		}
		result := atoi(line[spIx+1:sumIx-1]) + atoi(line[sumIx+2:spIx2])
		line  = line[:spIx+1] + itoa(result) + line[spIx2:]
		sumIx = strings.IndexByte(line, '+')
	}
	return line
}

// evaluates multiplications within a paranthesis free string
func evalMult(line string) (string) {
	mulIx := strings.IndexByte(line, '*')
	if mulIx == -1 {
		return line
	}
	for mulIx != -1 {
		spIx := strings.LastIndexByte(line[:mulIx-1], ' ')
		spIx2 := strings.IndexByte(line[mulIx+2:], ' ') 
		if spIx2 == -1 {
			spIx2 = len(line)
		} else {
			spIx2 += mulIx + 2
		}
		result := atoi(line[spIx+1:mulIx-1]) * atoi(line[mulIx+2:spIx2])
		line  = line[:spIx+1] + itoa(result) + line[spIx2:]
		mulIx = strings.IndexByte(line, '*')
	}
	return line
}

// MAIN ----
func main() {

	start  := time.Now()

	homework := readTxtFile("input_d18_p1.txt")

	// both parts togeher 
	sum1 := 0
	sum2 := 0
	for _,line := range homework {
		sum1 += parseLine(line, false)
		sum2 += parseLine(line, true)
	}
	fmt.Println("Homework sum part 1:", sum1)
	fmt.Println("Homework sum part 2:", sum2)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
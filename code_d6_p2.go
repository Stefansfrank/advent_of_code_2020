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

// Part 1: Questions that anyone in the group answered yes to
func countGroupAnswersOr(answers []string) (count int) {

	letterMap := make(map[rune]bool) // collects all letters for one group

	for _, line := range answers {

		if len(line) == 0 {          // next group, initialize things and count
			count += len(letterMap)
			letterMap = make(map[rune]bool)
			continue
		}

		for _, letter := range line { // collect letters in this line
			letterMap[letter] = true
		}
	}

	count += len(letterMap)
	return
}

// Part2: Questions that all in the group answered yes to
func countGroupAnswersAnd(answers []string) (count int) {

	letterMap    := make(map[rune]bool) // group specific map
	nxtLetterMap := make(map[rune]bool) // line specific map
	firstLine := true

	for _, line := range answers {

		if len(line) == 0 {             // next group, initialize things snd count
			count += len(letterMap)
			letterMap    = make(map[rune]bool)
			firstLine = true
			continue
		}

		for _, letter := range line {   // collect letter in this line
			nxtLetterMap[letter] = true
		}

		if firstLine {                  // merge the found letters with previous once in the group 
			letterMap = nxtLetterMap
			firstLine = false
		} else {
		    letterMap = mergeMaps(letterMap, nxtLetterMap)
		}

		// reset the line specific map
		nxtLetterMap = make(map[rune]bool)
	}

	count += len(letterMap)
	return
}

// For Part 2: takes two maps and creates a new map with only entries that are true in both
func mergeMaps(m1,m2 map[rune]bool) (nm map[rune]bool) {
	nm = make(map[rune]bool)
	for key, _ := range m1 {
		if m2[key] {
			nm[key] = true
		}
	}
	return
}

// MAIN ----
func main () {

	start  := time.Now()

	answers  := readTxtFile("input_d6_p1.txt")

	fmt.Printf("Part 1: Number of questions anyone answered yes to: %v\n", countGroupAnswersOr(answers))
	fmt.Printf("Part 2: Number of questions everyone answered yes to: %v\n", countGroupAnswersAnd(answers))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
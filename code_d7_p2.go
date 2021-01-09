package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
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

// helper structs
type bagDef struct { // a bag's name and list of references to other bags
	name    string
	content []bagRef
}

type bagRef struct { // a reference to another box with the amount and name
	amt  int
	name string
}

// parses the rules from the input 
func parseRules(rules []string) (bags map[string]bagDef) {
	re  := regexp.MustCompile(`([0-9]*)? ?([a-z]* [a-z]*) bags?`)
	bags = make(map[string]bagDef)

	for _, rule := range rules {
		newBag := bagDef{}
		newBag.content = []bagRef{}

		match := re.FindAllStringSubmatch(rule, -1)
		newBag.name = match[0][2]
		for i := 1; i < len(match); i++ {
			newBag.content = append(newBag.content, bagRef{ amt: atoi(match[i][1]), name: match[i][2] })	
		}
		bags[newBag.name] = newBag
	}
	return
}

// Part 1: start of the recursive match
func countRecBagMatch(bags map[string]bagDef, pattern string) (count int) {
	for _, bag := range bags {
		if recursiveMatch(bags, bag.content, pattern) {
			count++
		}
	}
	return
}

// Part 1: recursive depth first traversal identifying a match pattern
func recursiveMatch(bags map[string]bagDef, content []bagRef, pattern string) bool {
	for _, bagRef := range content {
		if bagRef.name == pattern {
			return true
		}
		if recursiveMatch(bags, bags[bagRef.name].content, pattern) {
			return true
		}
	}
	return false
}

// Part 2: recursive depth first traversal of all contained bags
func recursiveCount(bags map[string]bagDef, content []bagRef) (counter int) {	
	for _, bagRef := range content {
		counter += bagRef.amt
		counter += bagRef.amt*recursiveCount(bags, bags[bagRef.name].content) 
	}
	return counter
}


// MAIN ----
func main () {

	start  := time.Now()

	rules  := readTxtFile("input_d7_p1.txt")
    bags := parseRules(rules)

	// Part 1
    fmt.Println("Part 1:", countRecBagMatch(bags, "shiny gold"), "bag colors can contain shiny gold")

    // Part 2
    fmt.Printf("Part 2: Total bags in Shiny Gold: %v\n", recursiveCount(bags, bags["shiny gold"].content))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
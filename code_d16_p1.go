package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
	"regexp"
)

// standard file read ...
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
func atoi(inp string) int {
	i, _ := strconv.Atoi(inp)
	return i
}

// Helper structs and methods
type rule struct { // A rule has a parameter name and two validity ranges
	pNam string
	rng  dblRange
}

type dblRange struct { // A double range consists of a pair of from-to parameters
	from [2]int
	to   [2]int
}

type ticket struct { // A ticket has a list of parameters and a validity flag
	valid bool
	paras []int
}

func (dr *dblRange) contains(num int) bool { // Checks wether an int is in one of the ranges
	return ((num >= dr.from[0]) && (num <= dr.to[0])) || 
			((num >= dr.from[1]) && (num <= dr.to[1]))
}

func intersect(as []int, bs []int) []int { // Computes the intersection of two []int slices
	m  := make(map[int]bool)
	cs := []int{}
	for _, a := range as {
		m[a] = true
	}
	for _, b := range bs {
		if m[b] {
			cs = append(cs, b)
		}
	}
	return cs
}

// Parsing of input into the structures above
func parseInput(sheet []string) (rules []rule, tickets []ticket) {
	for ix, line := range sheet {
		if len(line) == 0 {
			rules    = parseRules(sheet[:ix])
			tickets  = parseTickets(append([]string{sheet[ix+2]}, sheet[ix+5:]...))
			break
		}
	}
	return
} 

func parseTickets(tickStr []string) (tickets []ticket) {
	for _, tkStr := range tickStr {
		tkt := ticket{ valid: true }
		prms := strings.Split(tkStr,",")
		for _, p := range prms {
			tkt.paras = append(tkt.paras, atoi(p))
		}
		tickets = append(tickets, tkt)
	}
	return
}

func parseRules(ruleStr []string) (rules []rule) {
	re   := regexp.MustCompile(`([[:alnum:]]+):\s([[:digit:]]+)-([[:digit:]]+)\sor\s([[:digit:]]+)-([[:digit:]]+)`)
	rules = []rule{}
	for _, line := range ruleStr {
		prms := re.FindAllStringSubmatch(line, -1)
		rl   := rule{ pNam: prms[0][1] , rng: dblRange{ from: [2]int{ atoi(prms[0][2]), atoi(prms[0][4]) },
														 to: [2]int{ atoi(prms[0][3]), atoi(prms[0][5]) } } }
		rules = append(rules, rl)
	}
	return
}

// creates a map with the key being all nubmers from 0 to 999 
// and the value a list of rules that the number fulfills
func validityMap(rules []rule) map[int][]int {
	valid := make(map[int][]int)
	for i:=0; i<1000; i++ {        
		valid[i] = []int{}
		for rIx, rl := range rules {
			if rl.rng.contains(i) {
				valid[i] = append(valid[i], rIx)
			}
		}
	}
	return valid
}

// finds all numbers in tickets that do not have a valid rule
// (where the validity map shows zero rules fulfilles) and invalidates these
func findErrorRate(tickets []ticket, valid map[int][]int) (int) {
	errRate := 0
	for i,tkt := range tickets {
		for _, para := range tkt.paras {
			if len(valid[para]) == 0 {
				errRate += para
				tickets[i].valid = false
			}
		}
	}
	return errRate
}

// maps each parameter index to all rules it fulfills
func paraRuleMap(valid map[int][]int, tickets []ticket) (map[int][]int) {
	prMap := make(map[int][]int)

	// initial values
	for i := 0; i < len(tickets[0].paras); i++ {
		prMap[i] = []int{}
		for j := 0; j < len(tickets[0].paras); j++ {
			prMap[i] = append(prMap[i], j)
		}
	}

	// reduce the rules using validity map
	for _, tkt := range tickets {
		if tkt.valid {
			for ix, para := range tkt.paras {
				prMap[ix] = intersect(prMap[ix], valid[para])
			} 
		}
	} 
	return prMap
}

// this maps each rule to the parameter index
// it assumes that there is one parameter that fulfills only one rule
// one parameter that fulfills only two rules (one being the one above) and so forth
func ruleParaIx(prMap map[int][]int) map[int]int {
	ruleParaIx := make(map[int]int)

	for i:=0; i<len(prMap); i++ {

		// this loop is a clunky way to find the
		// solution with exactly i+1 rules 
		// a possible optimization could be some kind of sorting
		for pix, rls := range prMap {
			if len(rls) == i+1 {
				for _, r := range rls {
					if _, exists := ruleParaIx[r]; !exists {
						ruleParaIx[r] = pix
					}
				}
				break
			}
		}
	}
	return ruleParaIx
}

// MAIN ----
func main () {

	start  := time.Now()

	sheet  := readTxtFile("input_d16_p1.txt")
    rules, tickets  := parseInput(sheet)

    // Part 1 needs only the validity map
    valid  := validityMap(rules)
    fmt.Println("Part 1 - Error Rate:", findErrorRate(tickets, valid))

    // Part 2 needs a bit of cross referencing maps
    prMap      := paraRuleMap(valid, tickets)
    ruleParaIx := ruleParaIx(prMap) 

	// the first six rules are the ones needed 
    part2result := 1
    for i:=0; i<6; i++ {
    	part2result *= tickets[0].paras[ruleParaIx[i]]
    }
    fmt.Println("Part 2 - Product of the departure parameters:", part2result)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
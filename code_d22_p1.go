package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
)

// standard text file read ...
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

// parsing the input generating the hands
func parseSheet(sheet []string) (hand [][]int) {
	hand = [][]int{ []int{}, []int{} }
	num := 0
	for _,line := range sheet {
		if len(line) == 0 {
			num++
			continue
		} 
		if (line[0] != 'P') {
			i, _ := strconv.Atoi(line)
			hand[num] = append(hand[num], i)
		}
	}
	return
}

// little helper
var bit = map[bool]int{false: 0, true: 1}

// a game holding the hand and a cache of already encountered hands
type game struct {
	hand  [][]int
	cache map[string]bool
	win   bool // win is true if player 2 wins
}

// comupute the key for the cache
func (g *game) key() string {
	return fmt.Sprint(g.hand[0], g.hand[1])
}

// play the game and return the win
func (g *game) play() bool {

	for len(g.hand[0]) > 0 && len(g.hand[1]) >0 {

		// infinite loop protection
		if g.cache[g.key()] {
			return false
		}
		g.cache[g.key()] = true

		// determine winner (incl. recursion)
		if g.hand[0][0] < len(g.hand[0]) && g.hand[1][0] < len(g.hand[1]) {
			gm  := game { hand:[][]int{ []int{}, []int{} } }
			gm.hand[0] = append(gm.hand[0], g.hand[0][1:1+g.hand[0][0]]...) // creates a true copy
			gm.hand[1] = append(gm.hand[1], g.hand[1][1:1+g.hand[1][0]]...) // creates a true copy
			gm.cache   = make(map[string]bool)
			g.win = gm.play()
		} else {
			g.win = g.hand[1][0]>g.hand[0][0] 
		}

		// manipulate hand according to outcome
		if g.win {
			g.hand[1] = append(g.hand[1][1:], g.hand[1][0], g.hand[0][0])
			g.hand[0] = g.hand[0][1:]			
		} else {
			g.hand[0] = append(g.hand[0][1:], g.hand[0][0], g.hand[1][0])
			g.hand[1] = g.hand[1][1:]
		}

	}
	return g.win
}

// counts the cards
func (g *game) count() (cnt int) {
	for i, pt := range(g.hand[bit[g.win]]) {
		cnt += pt * (len(g.hand[bit[g.win]]) - i)
	}
	return
}

func main () {

	start   := time.Now()
	sheet   := readTxtFile("input_d22_p1.txt")
	gm      := game{}
	gm.hand  = parseSheet(sheet)
	gm.cache = make(map[string]bool)
	gm.play()
	fmt.Println("Player", bit[gm.win]+1, "won with", gm.count(), "points.")

	fmt.Printf("Execution time: %v\n", time.Since(start))
}

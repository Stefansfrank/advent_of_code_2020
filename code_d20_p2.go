package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
	"regexp"
	"math/bits"
	"math"
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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// 10 bit inverse (as I represent the individual tiles as arrays of 10 bit integers)
func inv10(i uint16) uint16 {
	return bits.Reverse16(i) >> 6
}

// I represent the map information as bits in an int on the individual cards
type card struct {
	id    int
	rows  [10]uint16
	cols  [10]uint16
	edge  [4]uint16
	rot   int
	uniq  int
}

// Global variables (for lazyness
var emp map[uint16][]int // map that links a given edge value to the index of the tiles fulfilling it
var cards []card         // the list of tiles
var mapsz int            // the size of the global map (3 for test data, 12 for puzzle input) 

// sets the edge variables from the data
// (doubling information for easier readability)
func (c *card) setEdge() {
	c.edge[0] = c.cols[9]
	c.edge[1] = c.rows[9]
	c.edge[2] = c.cols[0]
	c.edge[3] = c.rows[0]
}

// tests whether a given edge value is unique i.e. no other tile has that edge value
func (c *card) edgeUnique(e int) bool {
	return len(emp[c.edge[e]]) == 1
}

// performs the next transformation in a sequence of the 8 transformations
// that cover all possible orientations 
func (c *card) next() {
	if c.rot % 4 == 3 {
		c.flip()
	} else {
		c.rotate()
	}
}

// rotate 90 degrees to ccw
func (c *card) rotate() {
	oc := c.cols
	or := c.rows
	for i := 0; i < 10; i++ {
		c.cols[i]   = inv10(or[i])
		c.rows[9-i] = oc[i]
	}
	c.rot = (c.rot + 1) % 8
	c.setEdge()
}

// flip around the vertical axis
func (c *card) flip() {
	oc := c.cols
	or := c.rows
	for i := 0; i < 10; i++ {
		c.rows[i]   = inv10(or[i])
		c.cols[9-i] = oc[i]
	}
	c.rot = (c.rot + 1) % 8
	c.setEdge()
}

// DEBUG helper only: dumping the content of a tile
func (c *card) dump() {
	fmt.Println("ID:", c.id)
	fmt.Println("Rows:")
	for _,row := range c.rows {
		fmt.Printf("%010b\n",row)
	}
	fmt.Println("Cols:")
	for _,col := range c.cols {
		fmt.Printf("%010b\n",col)
	}
	fmt.Println("Edges:")
	for i := 0; i < 4; i++ {
		fmt.Printf("%010b(%v)\n", c.edge[i], c.edge[i])
	}
	fmt.Println("Rotation:", c.rot)
}

// Parses the tiles from the input into cards[]
// and building the emp map to remember which tiles fulfillwhich edge values
func parseCards(sheet []string) {
	re := regexp.MustCompile(`Tile\s(\d+):`)

	cards = []card{}
	var cd card
	var cl [10]uint16
	var ix int
	var rw uint16
	for _,line := range sheet {
		if re.MatchString(line) {
			cd = card{ id:atoi(re.FindStringSubmatch(line)[1]) }	
			cl = [10]uint16{0,0,0,0,0,0,0,0,0,0}		
			cd.rows = [10]uint16{0,0,0,0,0,0,0,0,0,0}
			ix = 0		
			continue
		}
		if len(line) > 0 {
			tmp,_ := strconv.ParseUint(strings.ReplaceAll(strings.ReplaceAll(line,".","0"),"#","1"),2,10)
			rw     = uint16(tmp)
			cd.rows[ix] = rw
			for i := 9; i > -1; i-- {
				cl[i] <<= 1
				cl[i] |= rw & uint16(1)
				rw >>= 1
			}
			ix++
		} else {
			cd.cols  = cl
			cd.setEdge()
			cards    = append(cards, cd)
			for _,e := range cd.edge {
				ee := inv10(e)
				emp[e] =  append(emp[e], len(cards)-1)
				emp[ee] = append(emp[ee], len(cards)-1)
			}
		}
	}

	return
}

// finds the card that matches the top and left cards given
// edge parameter means: 1 = left edge, 2 = top edge, 0 cards on left and top
func findCard(top, left, edge int) int {
	lEdge := cards[left].edge[0]
	tEdge := cards[top].edge[1]
	var rng []int
	if edge == 1 {
		rng = emp[tEdge]
	} else {
		rng = emp[lEdge]
	}
	for _, e := range rng {
		if e != left && e != top {
			for j := 0; j < 8; j++ {
				if (edge == 0 && cards[e].edge[2] == lEdge && cards[e].edge[3] == tEdge) ||
				   (edge == 1 && cards[e].edgeUnique(2)    && cards[e].edge[3] == tEdge) ||
				   (edge == 2 && cards[e].edge[2] == lEdge && cards[e].edgeUnique(3)) {
					return e
				} else {
					cards[e].next()
				}
			}
		}
	}
	return -1			
}

// arranging the tiles with a given starting left top corner
// that must already be in the right orientation
func arrangeCards(c int) [][]int {

	// start top row
	mp  := [][]int{}
	row := []int{ c }

	// finish top row
	for i := 0; i < mapsz-1; i++ {
		e := findCard(0, row[i], 2)
		row = append(row, e)	
	} 
	mp = append(mp, row)

	// go through all remaining rows
	for i := 0; i < mapsz-1; i++ {

		last := mp[len(mp)-1]
		row  := []int{}
		e := findCard(last[0], 0, 1)
		row = append(row, e)	

		// then the remaining pieces
		for j := 0; j < mapsz-1; j++ {
			e := findCard(last[j+1], row[j], 0)
			row = append(row, e)
		}
		mp = append(mp, row)
	} 

	return mp
}

// matching and counting monsters and rocks
// since the length of each line is longer than 64, I can't no longer deal with ints 
// so I convert back to srings but stick with '0' and '1' so I can
// match to the monster by converting substrings back into binary and use a simple binary AND mask
func countMonsters(gmp [][]int) (num, tot int) {
	//                   # 
	// #    ##    ##    ###
	//  #  #  #  #  #  # 
	msk1 := uint64(0b00000000000000000010)
	msk2 := uint64(0b10000110000110000111)
	msk3 := uint64(0b01001001001001001000)


	// create map (as strings with the chars '0' and '1')
	gmap := []string{}
	for i := 0; i < mapsz; i++ {
		rmap := []string{"","","","","","","",""}
		for j := 0; j < mapsz; j++ {
			cd := cards[gmp[i][j]]
			for k := 0; k<8; k++ {
				rmap[k] = rmap[k] + fmt.Sprintf("%08b", cd.rows[k+1] & 0b0111111110 >> 1)
			}		
		}
		for j := 0; j < 8; j++ {
			gmap = append(gmap, rmap[j])
		}
	}

	// map created, count monsters and rocks
	// since I do not loop over the first and last line
	// they are manually added
	tmp := strings.ReplaceAll(gmap[0],"0","")
	tot += len(tmp)
	for i := 1; i < len(gmap)-1; i++ {
		for j := 0; j<len(gmap[i])-19; j++ {
			t2, _ := strconv.ParseUint(gmap[i][j:j+20],2,20)  
			if t2 & msk2 == msk2 {
				t3, _ := strconv.ParseUint(gmap[i+1][j:j+20],2,20)  
				if t3 & msk3 == msk3 {
					t1, _ := strconv.ParseUint(gmap[i-1][j:j+20],2,20) 
					if t1 & msk1 == msk1 {
						num++
					}
				}
			}
		}
		tmp = strings.ReplaceAll(gmap[i],"0","")
		tot += len(tmp)
	}
	tmp = strings.ReplaceAll(gmap[len(gmap)-1],"0","")
	tot += len(tmp) - num*15

	return
}

// MAIN ----
func main () {

	start  := time.Now()

	emp     = make(map[uint16][]int) // the global edge map mapping an edge value to the indices of the tiles that have that edge
	crn    := []int{} 				 // a list of corner tiles (listing their indices)
	sheet  := readTxtFile("input_d20_p1.txt")
	mapsz   = int(math.Sqrt(float64(len(sheet))/12)) // needed in order to quickly switch between test file and puzzle input by only changing the filename
	          parseCards(sheet)

	// NOTE: after playing with the puzzle input, it becomes clear that each edge value is used in 2 tiles maximum so there is only one matching tile
	// The input is also designed in a way that the corners have two unique edge values and edge pieces one unique edge value
	// Otherwise the code below would not work and I would have to traverse through permutations

	// Determine corner candidates
	// look for tiles with 2 unique edges for corners 
	// i.e. only 4 matches as each match has an inverted match ....
	for i,cd := range cards {
		fits := 0
		for _, e := range cd.edge {
			fits += len(emp[e])-1
			fits += len(emp[inv10(e)])-1
		}		
		if fits == 4 {
			crn = append(crn, i)
		}
	}

	// Part 1 solution
	idMul  := 1
	for _,i := range crn {
		fmt.Println("Corner ID:", cards[i].id)
		idMul *= cards[i].id
	}
	fmt.Println("Product of corners: ", idMul)

	// Part 2 solution
	// I loop trough all 8 transformations of the big map to hunt for monsters
	// that is done by looping through all 4 corners as the top left starting point
	// and using both orientations that one card can have in order to have the unique edges top and left
	var gmp [][]int
	fmt.Println("Map Variations:")
	for _, c := range crn {
		for i := 0; i < 8; i++ {
			if cards[c].edgeUnique(2) && cards[c].edgeUnique(3) {
				gmp = arrangeCards(c)
				num, tot := countMonsters(gmp)
				fmt.Println("Monsters:", num, "/ Total Rocks:", tot)
			}
			cards[c].next()
		}
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
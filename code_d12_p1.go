package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
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
func abs(x int) int {
	if x < 0 { return -x }
	return x
}

// ship
type ship struct {
	deg  int
	x    int
	y    int
}

var vects = map[int][]int{ 90: []int{1,0}, 180: []int{0,1}, 270: []int{-1,0}, 0: []int{0,-1}}
var degs  = map[byte]int{ 'E': 90, 'N': 0, 'W': 270, 'S': 180 }

// sail
func sail(shp ship, navs []string) (int, int) {
	for _, nav := range navs {
		cmd := nav[0]
		prm := atoi(nav[1:])
		switch cmd {
		case 'L':
			shp.deg = (shp.deg + 360 - prm) %360
		case 'R':
			shp.deg = (shp.deg + prm) % 360
		case 'E', 'N', 'W', 'S':
			shp.x += vects[degs[cmd]][0] * prm
			shp.y += vects[degs[cmd]][1] * prm
		case 'F':
			shp.x += vects[shp.deg][0] * prm
			shp.y += vects[shp.deg][1] * prm
		}
	}
	return shp.x, shp.y
}

// MAIN ----
func main () {

	start  := time.Now()

	navs  := readTxtFile("input_d12_p1.txt")
	shp   := ship{ deg: 90 }
	
	x,y := sail(shp, navs)
	fmt.Printf("Ship ends at %v,%v - Manhattan Distamce is %v\n", x, y, abs(x)+abs(y))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
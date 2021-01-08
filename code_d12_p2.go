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

// simple int abs function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// coordinate
type point struct {
	x    int
	y    int
}

// vectors representing cardinal directions
var vects = map[byte][]int{	'E': []int{1,0}, 'S': []int{0,1}, 'W': []int{-1,0}, 'N': []int{0,-1} }

// sail
func sail(shp, wpt point, navs []string) (int, int) {
	for _, nav := range navs {
		cmd := nav[0]
		prm := atoi(nav[1:])
		switch cmd {
		case 'L', 'R':
			if cmd == 'L' {
				prm = 360 - prm
			}
			switch prm {
			case 90:
				wpt.x, wpt.y = wpt.y, wpt.x
				wpt.x = -wpt.x
			case 180:
				wpt.x = -wpt.x
				wpt.y = -wpt.y
			case 270:
				wpt.x, wpt.y = wpt.y, wpt.x
				wpt.y = -wpt.y
			}
		case 'E', 'N', 'W', 'S':
			wpt.x += vects[cmd][0] * prm
			wpt.y += vects[cmd][1] * prm
		case 'F':
			shp.x += wpt.x * prm
			shp.y += wpt.y * prm
		}
	}
	return shp.x, shp.y
}

// MAIN ----
func main () {

	start  := time.Now()

	navs  := readTxtFile("input_d12_p1.txt")
	shp   := point{}
	wpt   := point{x:10,y:-1}

	x,y := sail(shp, wpt, navs)
	fmt.Printf("Ship ends at %v,%v - Manhattan Distamce is %v\n", x, y, abs(x)+abs(y))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
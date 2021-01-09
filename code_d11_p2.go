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

// floor space
type space struct {
	seat  bool  // there is a seat
	occ   bool  // the seat is occupied
	visN  int   // the amount of neighbouring seats visible
}

// map
type smap struct {
	wdt    int   // width
	hgt    int   // height
	numOcc int   // total number of occupied seats
	data   [][]space // all the spaces
	changed bool // indicates whether the seating has changed 
}

// parsing the raw data into the structs above
func initData(wdt, hgt int) (data [][]space) {
	data = make([][]space, hgt + 2)
	for i, _ := range data {
		data[i] = make([]space, wdt +2)
	}
	return
}
func parseRawMap(raw []string) (pMap smap) {
	pMap.wdt = len(raw[0])
	pMap.hgt = len(raw)
	pMap.data = initData(pMap.wdt, pMap.hgt)

	for rIx, row := range raw {
		for cIx, spc := range row {
			if spc == 'L' {
				pMap.data[rIx+1][cIx+1].seat = true
			}
		}
	}
	pMap.changed = true
	return
}

// DEBUGGING only: prints the current map
func dumpMap(cMap smap) {
	for y := 1; y <= cMap.hgt; y++ {
		for x := 1; x <= cMap.wdt; x++ {
			if !cMap.data[y][x].seat {
				fmt.Print(".")
			} else {
				if cMap.data[y][x].occ {
					fmt.Print("#")
				} else {
					fmt.Print("L")
				}
			}
		}
		fmt.Println()
	}
	fmt.Printf("Changed: %v / Occupied seats: %v\n", cMap.changed, cMap.numOcc)
}

// define the 8 directions to loop over for the detection of visual neighbors 
var dircs = [][]int { []int{-1,-1},[]int{-1,0},[]int{0,-1},[]int{1,1},
						[]int{0,1},[]int{1,0},[]int{1,-1},[]int{-1,1} }

// test function to detect whether a coordinate hits the edge of the map
func inside(x,y,wdt,hgt int) bool {
	return x > 0 && x <= wdt && y > 0 && y <= hgt
}

// iterates once over the map
func iterateMap(oMap smap) (nMap smap) {
	nMap.data = initData(oMap.wdt, oMap.hgt)
	nMap.wdt, nMap.hgt = oMap.wdt, oMap.hgt

	for y := 1; y <= nMap.hgt; y++ {
		for x := 1; x <= nMap.wdt; x++ {

			// do nothing if no seat
			nMap.data[y][x].seat = oMap.data[y][x].seat
			if !oMap.data[y][x].seat {
				continue
			}

			// calculate visible neighbours by traversing in all 8 directions
			// until either an empty chair or occupied chair or the edge is reached
			for _,d := range dircs {
				xx := x
				yy := y 
				for inside(xx,yy,oMap.wdt,oMap.hgt)	{
					xx += d[0]
					yy += d[1]
					if oMap.data[yy][xx].seat {
						if oMap.data[yy][xx].occ {
							nMap.data[y][x].visN += 1
						}
						break
					}
				}
			}

			// not occupied and no neigbors -> occupy
			// occupied and less than 5 visible neighbours -> occupy again
			if (nMap.data[y][x].visN == 0 && !oMap.data[y][x].occ) ||
				(nMap.data[y][x].visN < 5 && oMap.data[y][x].occ) {

				nMap.data[y][x].occ = true
				nMap.numOcc += 1
			}

			nMap.changed = nMap.changed || (nMap.data[y][x].occ != oMap.data[y][x].occ)
		}
	}

	return
}

// MAIN ----
func main () {

	start  := time.Now()

	rawMap  := readTxtFile("input_d11_p1.txt")
	curMap  := parseRawMap(rawMap)

	for curMap.changed {
		curMap  = iterateMap(curMap)
	}
	fmt.Println("Occupied Seats:", curMap.numOcc)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
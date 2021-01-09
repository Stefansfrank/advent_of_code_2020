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
	occN  int   // the amount of neighbouring seats occupied
}

// map
type smap struct {
	wdt    int   // width
	hgt    int   // height
	numOcc int   // total number of occupied seats
	data   [][]space // all the spaces
	changed bool // indicates whether the seating has changed 
}

// parsing the raw input into these objects
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

// DEBUGGING only - prints the map
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

// iterates once on the map
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

			// not occupied and no neigbors -> occupy
			// occupied and less than 4 neighbours -> occupy again
			if (oMap.data[y][x].occN == 0 && !oMap.data[y][x].occ) ||
				(oMap.data[y][x].occN < 4 && oMap.data[y][x].occ) {

				nMap.data[y][x].occ = true
				nMap.numOcc += 1

				// add to the neighbour counter of all neighbours
				for yy := y-1; yy < y+2; yy++ {
					for xx := x-1; xx < x+2; xx++ {
						nMap.data[yy][xx].occN += 1
					}
				}			
				nMap.data[y][x].occN--
			}

			// track whether changes have occured
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
	// dumpMap(curMap)
	fmt.Println("Occupied Seats:", curMap.numOcc)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
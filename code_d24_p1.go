package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile(name string) (lines []string) {
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}
	return
}

// global variables representing the floor and a helper map
var tls = map[string][]int{} // the key is a string based on x,y and the value is []int{x,y} 
var dlt = map[string][]int{ "w":[]int{-1,0},"e":[]int{1,0},
							"nw":[]int{0,-1},"ne":[]int{1,-1},
							"sw":[]int{-1,1},"se":[]int{0,1} }

// parsing the directions into individual directions
func parseDirs(sheet []string) [][]string {
	dirs := [][]string{}
	for _, line := range sheet {
		dir := []string{}
		for i := 0; i < len(line); i++ {
			d := string(line[i])
			if d[0] == 'n' || d[0] == 's' {
				d = d + string(line[i+1])
				i++
			}
			dir = append(dir, d) 
		}
		dirs = append(dirs, dir)
	}
	return dirs
}

// computing a string key for the tile map based on coordinates
func key(x,y int) string {
	return fmt.Sprint(x, y)
}

// traverses the tiles for Part 1 to create the initial pattern
func traverse(dirList [][]string) {
	for _,dirs := range dirList {
		x, y := 0,0
		for _,dir := range dirs {
			dl := dlt[dir]
			x += dl[0]
			y += dl[1]
		}
		ky  := key(x,y)
		col := len(tls[ky]) > 0
		if col {
			delete(tls, ky)
		} else {
			tls[ky] = []int{x,y}
		}
	}
}

// iterates the pattern once (i.e. simulates one day)
func iterate() {
	// make a copy of the map 
	newTls := make(map[string][]int)
	for k,v := range tls {
  		newTls[k] = v
	}

	// determine tiles to be checked i.e. all black tiles and all their imminent neighbors
	tlList := map[string][]int{}
	for k,v := range tls {
		tlList[k] = v
		for _, dl := range dlt {
			tlList[key(v[0]+dl[0],v[1]+dl[1])] = []int{v[0]+dl[0],v[1]+dl[1]}
		}
	}

	// check the nighbouring tiles and set copy accordingly
	for k,v := range tlList {
		num := 0
		for _, dl := range dlt {
			num += len(tls[key(v[0] + dl[0], v[1] + dl[1])])
		}
		if num == 0 || num > 4 {
			delete(newTls, k)
		} else if num == 4 {
			newTls[k] = v
		}
	}

	// copy back
	tls = newTls
}

// MAIN ----
func main () {

	start := time.Now()

	sheet   := readTxtFile("input_d24_p1.txt")
	dirList := parseDirs(sheet)

	// Part 1
	traverse(dirList)
	fmt.Println("Number of initial black tiles: ", len(tls))

	// Part 2
	for i := 0; i < 100; i++ {
		iterate()
	}
	fmt.Println("Number of black tiles after 100 days: ", len(tls))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}

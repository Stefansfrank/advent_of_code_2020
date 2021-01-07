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

// helper struct holding the space (as a map to boolean) and the limits of the space
// defined as the extension of coordinates containing all active cells
type space4D struct {
	xto   int
	xfrom int
	yto   int
	yfrom int
	zto   int
	zfrom int
	wfrom int
	wto   int
	mp    map[string]bool
}
type space struct {
	xto   int
	xfrom int
	yto   int
	yfrom int
	zto   int
	zfrom int
	mp    map[string]bool
}

// the key to my map
func key4D(x,y,z,w int) string {
	return fmt.Sprint(x) + "/" + fmt.Sprint(y) + "/" + fmt.Sprint(z) + "/" + fmt.Sprint(w)
}
func key(x,y,z int) string {
	return fmt.Sprint(x) + "/" + fmt.Sprint(y) + "/" + fmt.Sprint(z) 
}

// building the map
func parseInit4D(init []string) space4D {

	spc := space4D { xto: len(init[0]),
		xfrom: 0,
		yto: len(init),
		yfrom: 0,
		zto:   1,
		zfrom: 0,
		wto:   1,
		wfrom: 0,
		mp: make(map[string]bool)}

	for y, line := range init {
		for x, rune := range line {
			if rune == '#' {
				spc.mp[key4D(x,y,0,0)] = true
			}
		}
	}

	return spc
}
func parseInit(init []string) space {

	spc := space { xto: len(init[0]),
		xfrom: 0,
		yto: len(init),
		yfrom: 0,
		zto:   1,
		zfrom: 0,
		mp: make(map[string]bool)}

	for y, line := range init {
		for x, rune := range line {
			if rune == '#' {
				spc.mp[key(x,y,0)] = true
			}
		}
	}

	return spc
}

// determines whether a certain space should be active in the next iteration
// ... since this method counts the space itself, the conditions look slightly
// different than in the description
func active4D(spc space4D, x,y,z,w int) bool {
	cnt := 0
	act := spc.mp[key4D(x,y,z,w)]
	for xx:=x-1; xx<x+2; xx++ {
		for yy:=y-1; yy<y+2; yy++ {
			for zz:=z-1; zz<z+2; zz++ {
				for ww:=w-1; ww<w+2; ww++ {
					if spc.mp[key4D(xx,yy,zz,ww)] {
						cnt++
					}
				}
			}
		}
	}
	return cnt == 3 || (cnt == 4 && act) 
}
func active(spc space, x,y,z int) bool {
	cnt := 0
	act := spc.mp[key(x,y,z)]
	for xx:=x-1; xx<x+2; xx++ {
		for yy:=y-1; yy<y+2; yy++ {
			for zz:=z-1; zz<z+2; zz++ {
				if spc.mp[key(xx,yy,zz)] {
					cnt++
				}
			}
		}
	}
	return cnt == 3 || (cnt == 4 && act) 
}

// iterate on space and return new space
func iterate4DLife(spc space4D) (new space4D) {

	// make the new space one bigger in all directions
	new = space4D { xto: spc.xto+1,
		xfrom: spc.xfrom-1,
		yto:   spc.yto+1,
		yfrom: spc.yfrom-1,
		zto:   spc.zto+1,
		zfrom: spc.zfrom-1,
		wto:   spc.wto+1,
		wfrom: spc.wfrom-1,
		mp: make(map[string]bool)}

	// looping through the whole space
	// NOTE: if I would want to optimize my code
	// I would build a list of candidates only containing
	// cells in the vicinity of an active one
	for x:=new.xfrom; x<new.xto; x++ {
		for y:=new.yfrom; y<new.yto; y++ {
			for z:=new.zfrom; z<new.zto; z++ {
				for w:=new.wfrom; w<new.wto; w++ {
					if active4D(spc,x,y,z,w) {
						new.mp[key4D(x,y,z,w)] = true
					}
				}
			}
		}
	}

	return
}
func iterateLife(spc space) (new space) {

	// make the new space one bigger in all directions
	new = space { xto: spc.xto+1,
		xfrom: spc.xfrom-1,
		yto:   spc.yto+1,
		yfrom: spc.yfrom-1,
		zto:   spc.zto+1,
		zfrom: spc.zfrom-1,
		mp: make(map[string]bool)}

	// looping through the whole space
	// NOTE: if I would want to optimize my code
	// I would build a list of candidates only containing
	// cells in the vicinity of an active one
	for x:=new.xfrom; x<new.xto; x++ {
		for y:=new.yfrom; y<new.yto; y++ {
			for z:=new.zfrom; z<new.zto; z++ {
				if active(spc,x,y,z) {
					new.mp[key(x,y,z)] = true
				}
			}
		}
	}

	return
}

// MAIN ----
func main () {

	start  := time.Now()

	init   := readTxtFile("input_d17_p1.txt")
	spc1   := parseInit(init)
	spc2   := parseInit4D(init)

	n := 6
	for i:=0; i<n; i++ {
		spc1  = iterateLife(spc1)
		spc2  = iterate4DLife(spc2)
	}

	fmt.Printf("Active cubes after %v cycles: 3D space %v, 4D space %v\n", n, len(spc1.mp), len(spc2.mp))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
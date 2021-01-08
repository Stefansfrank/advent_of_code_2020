package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
)

// Standard text file parser 
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

// Execution of the program
// I do use a pair of masks:
// - one called oneMask  that has a 1 wherever a 1 needs enforced (using OR)
// - one called zeroMask that has a 0 wherever a 0 needs enforced (using AND)
// For part 2 I create multiple mask pairs for each possible permutation created by an "X"

func runInpPart1(codeLines []string) (mem map[int]uint64) {

	var oneMask, zeroMask uint64
	mem = make(map[int]uint64)
	for _, code := range codeLines {

		switch code[:3] {
		case "mas":
			oneMask, zeroMask = 0, 0
			for _, bit := range code[7:] {
				switch bit {
				case '1': 
					oneMask  += 1
					zeroMask += 1 // in order to not enforce a zero, the zero mask needs a 1 as well
				case 'X':
					zeroMask += 1
				} // no explicit case '0' needed as not setting a bit enforces zero
				oneMask  <<= 1
				zeroMask <<= 1
			}
			oneMask  >>= 1
			zeroMask >>= 1
		case "mem":
			ix := strings.IndexRune(code, ']')
			memKey,_ := strconv.Atoi(code[4:ix])
			memVal,_ := strconv.ParseUint(code[ix+4:], 10, 64)
			mem[memKey] = (memVal & zeroMask) | oneMask
		}
	}
	return
}

func runInpPart2(codeLines []string) (mem map[uint64]uint64) {

	var oneMask, zeroMask []uint64
	mem = make(map[uint64]uint64)
	for _, code := range codeLines {

		switch code[:3] {
		case "mas":
			oneMask, zeroMask = []uint64{0}, []uint64{0}
			for _, bit := range code[7:] {

				// shit the existing masks to the left
				for i,_ := range oneMask {
					oneMask[i]  <<= 1
					zeroMask[i] <<= 1
				}

				// apply the last bit according to the rules
				switch bit {
				case '0':
					for i,_ := range zeroMask {						   
						zeroMask[i] +=1 // in order to not enforce a zero, the zero mask needs a 1
					}
				case '1':
					for i,_ := range oneMask {
						oneMask[i] +=1
						zeroMask[i] +=1 // in order to not enforce a zero, the zero mask needs a 1
					}
				case 'X':
					len := len(oneMask)
					for i := 0; i < len; i++ {
						oneMask  = append(oneMask, oneMask[i])   // add a new permutation at the end
						zeroMask = append(zeroMask, zeroMask[i]) // not setting any bits enforces a zero
						oneMask[i]  +=1 // enforce a 1 on the existing permutation
						zeroMask[i] +=1 // i.e. zero Mask needs to be 1 as well
					}
				}
			}
		case "mem":
			ix := strings.IndexRune(code, ']')
			memKey,_ := strconv.ParseUint(code[4:ix], 10, 64)
			memVal,_ := strconv.ParseUint(code[ix+4:], 10, 64)
			for i,_ := range oneMask {
				mem[(memKey & zeroMask[i]) | oneMask[i]] = memVal
			}
		}
	}
	return
}


// MAIN ----
func main () {

	start := time.Now()

	code  := readTxtFile("input_d14_p1.txt")

	mem    := runInpPart1(code)
	result := uint64(0)
	for _, val := range mem {
		result += val
	}
	fmt.Printf("Mem sum part 1: %v\n", result)

	mem2  := runInpPart2(code)
	result = uint64(0)
	for _, val := range mem2 {
		result += val
	}
	fmt.Printf("Mem sum part 2: %v\n", result)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
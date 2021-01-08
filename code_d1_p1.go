 package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile2Int (name string) (nums []int) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		nums = append(nums, atoi(scanner.Text()))
	}
	return
}
// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}


// MAIN ----
func main () {

	start := time.Now()

	exps := readTxtFile2Int("input_d1_p1.txt")
	exLen := len(exps)

	for i := 0; i < exLen-1; i++ {
		for j := i+1; j < exLen; j++ {
			if exps[i] + exps[j] == 2020 {
				fmt.Printf("%v * %v = %v\n", exps[i], exps[j], exps[i]*exps[j])
			}
		}
	}
	
	for i := 0; i < exLen-2; i++ {
		for j := i+1; j < exLen-1; j++ {
			for k := j+1; k < exLen; k++ {
				if exps[i] + exps[j] + exps[k] == 2020 {
					fmt.Printf("%v * %v * %v = %v\n", exps[i], exps[j], exps[k], exps[i]*exps[j]*exps[k])
				}
			}
		}
	}

	fmt.Printf("Execution time: %v\n", time.Since(start))
}
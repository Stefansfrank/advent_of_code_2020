package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"sort"
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

// structures
type food struct {
	ings []string
	alls []string
}
var foods  []food
var allMap  map[string][]int
var ingMap  map[string][]int
var ing2all map[string][]string // contains POTENTIAL allergens per ingredient
var all2ing map[string]string   // contains the actual link from allergen to ingredient

// intersect two slices of strings
func intersect(a, b []string) (c []string) {
	c = []string{}
	for _, aa := range a {
		for _, bb := range b {
			if aa == bb {
				c = append(c, aa)
			}
		}
	}
	return
}

func parseSheet(sheet []string) {

	foods  = []food{}
	allMap = make(map[string][]int)
	ingMap = make(map[string][]int)
	for _, line := range sheet {

		tmp  := strings.Split(line, " (contains ")
		ings := strings.Split(tmp[0], " ")
		alls := strings.Split(tmp[1], " ")

		for _, ing := range ings {
			ingMap[ing] = append(ingMap[ing], len(foods))
		}
		for i, all := range alls {
			alls[i] = strings.Trim(all, ",)")
			allMap[alls[i]] = append(allMap[alls[i]], len(foods))
		}

		foods = append(foods, food{ ings: ings, alls: alls})
	}
}

func detectAllergenFree() (free []string, numFd int) {

	// a list mapping ingredients to POTENTIAL allergens
	ing2all  = make(map[string][]string)
	free     = []string{}
	ingLst  := []string{}

	for all, fdIxLst := range allMap {
		for i, fdIx := range fdIxLst {
			if i == 0 {
				ingLst = foods[fdIx].ings
			} else {
				ingLst = intersect(ingLst, foods[fdIx].ings)
			}
		}
		for _, ing := range ingLst {
			ing2all[ing] = append(ing2all[ing], all)
		}
	}

	for ing, alls := range ingMap {
		if len(ing2all[ing]) == 0 {
			free = append(free, ing)
			numFd += len(alls)
		}
	}
	return
}

// uses the map of potential allergens for each ingredient map ing2all
// to derive the 1:1 map from allergens to ingredients
// by starting to resolve entries from ing2all that are of len() 1 (direct match) 
// and then remove matched allergens from the unresolved entries in ing2all 
// which will again create len() 1 entries and so forth
func mapAll2Ing() {
	all2ing = make(map[string]string)

	for len(ing2all) > len(all2ing) {
		for ing, alls := range ing2all {
			if len(alls) == 1 {
				all2ing[alls[0]] = ing
				ing2all[ing] = []string{}
			} else {
				new := []string{}
				for _, all := range alls {
					if len(all2ing[all]) == 0 {
						new = append(new, all)
					}
				}
				ing2all[ing] = new
			}
		}
	}
}


func main() {

	start  := time.Now()
	sheet  := readTxtFile("input_d21_p1.txt")
	parseSheet(sheet)

	// Part 1
	_, num := detectAllergenFree()
	fmt.Println("Number of allergen free ingredients", num)

	// Part 2
	mapAll2Ing()

	// sort the keys
	keys := []string{}
	for k := range all2ing {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// run through in a sorted way
	result := ""
	for _,k := range keys {
		result += all2ing[k] + ","
	}
	fmt.Println("Canonical dangerous ingredients list:",result[:len(result)-1])

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}

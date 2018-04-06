package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	csvPath     = "warehouse-grid.csv"
	shelfLength = 1.0
	shelfWidth  = 1.0
	pathWidth   = 1.0
)

//ReadCSV returns a 2D array of string from the csv file
func ReadCSV(path string) ([][]string, error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records, err
}

// Product defines the information of a product
type Product struct {
	id, x, y   int
	l, r, u, d bool
	//num int
}

func coordinateConverter(x, y int) (int, int) {
	return 2*x + 1, 2*y + 1
}

func posAssigner(prod *Product) *Product {
	prod.l, prod.r = true, true
	return prod
}

// ParseProductInfo returns a map that includes product info
func ParseProductInfo(path string) map[int]Product {
	records, err := ReadCSV(path)
	if err != nil {
		log.Fatal(err)
	}
	var m map[int]Product
	for _, s := range records {
		var temp [3]int
		var err error
		m = make(map[int]Product)
		for i := range temp {
			s[i] = strings.TrimSpace(s[i])
			switch i {
			case 0:
				temp[i], err = strconv.Atoi(s[i])
			default:
				temp[i], err = strconv.Atoi(strings.Split(s[i], ".")[0])
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		temp[1], temp[2] = coordinateConverter(temp[1], temp[2])
		prod := Product{id: temp[0], x: temp[1], y: temp[2]}
		m[temp[0]] = *posAssigner(&prod)
	}
	return m
}

// ReadInput returns 3 int from stdin
func ReadInput() (int, int, int) {
	var strInput [3]string
	var input [3]int
	_, err := fmt.Scan(&strInput[0], &strInput[1], &strInput[2])
	if err != nil {
		log.Fatal(err)
	}
	for i := range strInput {
		strInput[i] = strings.TrimSpace(strInput[i])
		input[i], err = strconv.Atoi(strInput[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	return input[0], input[1], input[2]
}

// FindDest returns the destination given init position & product to fetch
func FindDest(x, y int, prod Product) (int, int) {
	if x < prod.x {
		return prod.x - 1, prod.y
	}
	return prod.x + 1, prod.y
}

func main() {
	m := ParseProductInfo(csvPath)
	x, y, id := ReadInput()

	fmt.Print(x, y, m[id])
}

// CONSIDER ON THE SAME LINE SITUATION! e.g.: (0,1) -> (3,1)

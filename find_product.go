package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var csvPath = "warehouse-grid.csv"

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
	id int
	x  int
	y  int
	//num int
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
		prod := Product{temp[0], temp[1], temp[2]}
		m[temp[0]] = prod
	}
	return m
}

func main() {
	m := ParseProductInfo(csvPath)
	fmt.Println(m[2629382])
}

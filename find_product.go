package main

import (
	"fmt"
	"os"
	"log"
	"encoding/csv"
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

func main() {
	fmt.Println(ReadCSV(csvPath))
}

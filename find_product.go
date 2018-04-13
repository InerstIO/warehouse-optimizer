package main

import (
	"fmt"
	"warehouse-optimizer/warehouse"
	"log"
)

const (
	gridPath     = "warehouse-grid.csv"
)

func main() {
	m := warehouse.ParseProductInfo(gridPath)
	fmt.Print("Please input x y prod_id:\n")
	x, y, id := warehouse.ReadInput()
	if x*y%2 == 1 {
		log.Fatal("Cannot start on a shelf.")
	}

	prod, ok := m[id]
	if ok {
		src := warehouse.Point{x, y}
		path, length := warehouse.FindPath(src, warehouse.FindDest(src, prod))
		fmt.Printf("%v\nlength: %v\n", path, length)
	} else {
		fmt.Print("prod_id not exist.\n")
	}
	//fmt.Print(warehouse.ReadOrder())
}

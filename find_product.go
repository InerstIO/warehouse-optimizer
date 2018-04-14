package main

import (
	"fmt"
	"log"
	"warehouse-optimizer/warehouse"
)

const (
	gridPath = "warehouse-grid.csv"
)

func main() {
	m := warehouse.ParseProductInfo(gridPath)
	pathInfo := warehouse.BuildPathInfo(gridPath)
	fmt.Println("Hello User, where is your worker? e.g.:\"2 4\"")
	x, y := warehouse.ReadInput()
	start := warehouse.Point{X: x, Y: y}
	if x*y%2 == 1 {
		log.Fatal("Cannot start on a shelf.")
	}
	fmt.Println("What is your worker's end location? e.g.:\"0 18\"")
	x, y = warehouse.ReadInput()
	end := warehouse.Point{X: x, Y: y}
	if x*y%2 == 1 {
		log.Fatal("Cannot end on a shelf.")
	}
	fmt.Println("Hello User, what items would you like to pick? (separate by space)")
	orders := warehouse.ReadOrder(m)
	fmt.Println("Here is the optimal picking order:")
	optimalOrder := warehouse.BruteForceOrderOptimizer(orders[0], start, end, m, pathInfo)
	fmt.Println(optimalOrder)
	fmt.Println("Here is the optimal path:")
	dest := warehouse.FindDest(start, m[optimalOrder[0]])
	s := fmt.Sprintf("%v->", warehouse.FindPath(start, dest))
	s += fmt.Sprintf("[pick up %v from %v]->", optimalOrder[0], m[optimalOrder[0]].Pos)
	var src warehouse.Point
	for _, prod := range optimalOrder[1:] {
		src = dest
		dest = warehouse.FindDest(start, m[prod])
		s += fmt.Sprintf("%v->", warehouse.FindPath(src, dest))
		s += fmt.Sprintf("[pick up %v from %v]->", prod, m[prod].Pos)
	}
	src = dest
	s += fmt.Sprint(warehouse.FindPath(src, end))
	fmt.Println(s)
	fmt.Printf("Total distance traveled: %v\n", warehouse.RouteLength(optimalOrder, start, end, m, pathInfo))

	/*prod, ok := m[id]
	if ok {
		path, length := warehouse.FindPath(src, warehouse.FindDest(src, prod))
		fmt.Printf("%v\nlength: %v\n", path, length)
	} else {
		fmt.Print("prod_id not exist.\n")
	}*/
	//fmt.Println(warehouse.BruteForceOrderOptimizer(warehouse.Order{4,123, 67}, warehouse.Point{0, 0}, warehouse.Point{4, 6}, warehouse.ParseProductInfo(gridPath)))
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"warehouse-optimizer/warehouse"
)

const (
	gridPath  = "warehouse-grid.csv"
	dimPath   = "item-dimensions-tabbed.txt"
	timeLimit = 10.0
)

func main() {
	dim := warehouse.ParesDimensionInfo(dimPath)
	m := warehouse.ParseProductInfo(gridPath, dim)
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
	var op, t, iter int
	fmt.Println("Type 0 for Nearest Neighbor Optimizer, type 1 for Branch & Bound Optimizer (slow!!)")
	var strInput string
	_, err := fmt.Scan(&strInput)
	if err != nil {
		log.Fatal(err)
	}
	strInput = strings.TrimSpace(strInput)
	op, err = strconv.Atoi(strInput)
	optimizer := func(op int, o warehouse.Order, start, end warehouse.Point, m map[int]warehouse.Product,
		pathInfo map[warehouse.Point]map[warehouse.Point]float64, iteration ...int) warehouse.Order {
		if op == 0 {
			return warehouse.NNIOrderOptimizer(o, start, end, m, pathInfo, iteration...)
		} else {
			return warehouse.BnBOrderOptimizer(o, start, end, m, pathInfo, timeLimit)
		}
	}
	if op == 0 {
		fmt.Println("What's the max number of iterations you want? (0 for max available)")
		_, err := fmt.Scan(&strInput)
		if err != nil {
			log.Fatal(err)
		}
		strInput = strings.TrimSpace(strInput)
		iter, err = strconv.Atoi(strInput)
	}

	for {
		fmt.Println("Type 1 to manual input, type 2 to file input.")
		_, err := fmt.Scan(&strInput)
		if err != nil {
			log.Fatal(err)
		}
		strInput = strings.TrimSpace(strInput)
		t, err = strconv.Atoi(strInput)
		if err != nil {
			log.Fatal(err)
		}
		if t == 1 || t == 2 {
			break
		}
	}
	if t == 1 {
		fmt.Println("Hello User, what items would you like to pick? (separate by space)")
		orders := warehouse.ReadOrder(m)
		fmt.Println("Here is the optimal picking order:")
		//optimalOrder := warehouse.BruteForceOrderOptimizer(orders[0], start, end, m, pathInfo)
		optimalOrder := optimizer(op, orders[0], start, end, m, pathInfo, iter)
		fmt.Println(optimalOrder)
		fmt.Println("Here is the optimal path:")
		s := warehouse.Route2String(optimalOrder, start, end, m)
		fmt.Println(s)
		fmt.Printf("Total distance traveled: %v\n", warehouse.RouteLength(optimalOrder, start, end, m, pathInfo))
		if effort, missWeightData := warehouse.RouteEffort(optimalOrder, start, end, m, pathInfo); missWeightData {
			fmt.Printf("There are some item(s) with no weight data, and the effort of this path is at least %v.\n", effort)
		} else {
			fmt.Printf("The effort is %v.\n", effort)
		}
	} else if t == 2 {
		fmt.Println("Please list file of orders to be processed:")
		ordersPath := warehouse.ReadString()
		orders := warehouse.ParesOrderInfo(ordersPath)
		fmt.Println("Please list output file:")
		outputPath := warehouse.ReadString()
		fmt.Println("Computing...")
		outputFile, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		var ods []warehouse.Order
		for i, order := range orders {
			var result warehouse.Order
			if op == 0 {
				result = warehouse.NNIOrderOptimizer(order, start, end, m, pathInfo)
			} else {
				result = warehouse.BnBOrderOptimizer(order, start, end, m, pathInfo, timeLimit)
			}
			ods = append(ods, result)
			fmt.Println(i)
		}
		if err := ioutil.WriteFile(outputPath, warehouse.Routes2JSON(ods, start, end, m), 0777); err != nil {
			log.Fatalln("error writing results to json:", err)
		}
	}

	/*prod, ok := m[id]
	if ok {
		path, length := warehouse.FindPath(src, warehouse.FindDest(src, prod))
		fmt.Printf("%v\nlength: %v\n", path, length)
	} else {
		fmt.Print("prod_id not exist.\n")
	}*/
	//fmt.Println(warehouse.BruteForceOrderOptimizer(warehouse.Order{4,123, 67}, warehouse.Point{0, 0}, warehouse.Point{4, 6}, warehouse.ParseProductInfo(gridPath)))
}

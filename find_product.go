package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"warehouse-optimizer/warehouse"
)

const (
	gridPath = "warehouse-grid.csv"
	dimPath = "item-dimensions-tabbed.txt"
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
			return warehouse.BnBOrderOptimizer(o, start, end, m, pathInfo)
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

		writer := csv.NewWriter(outputFile)
		defer writer.Flush()

		orderCtr := 0
		for i, order := range orders {
			if err := writer.Write([]string{"##Order Number##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{strconv.Itoa(i + 1)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Worker Start Location##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{fmt.Sprint(start)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"## Worker End Location##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{fmt.Sprint(end)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Original Parts Order##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write(warehouse.Order2csv(order)); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Optimized Parts Order##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			optimalOrder := optimizer(op, order, start, end, m, pathInfo, iter)
			if err := writer.Write(warehouse.Order2csv(optimalOrder)); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Original Parts Total Distance##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{strconv.FormatFloat(warehouse.RouteLength(order, start, end, m, pathInfo), 'G', -1, 64)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Optimized Parts Total Distance##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{strconv.FormatFloat(warehouse.RouteLength(optimalOrder, start, end, m, pathInfo), 'G', -1, 64)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"##Path of optimized order##"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{warehouse.Route2String(optimalOrder, start, end, m)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			if err := writer.Write([]string{"------------------------------------------------"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
			orderCtr++
		}
		fmt.Printf("%v orders processed.", orderCtr)
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

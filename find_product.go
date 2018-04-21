package main

import (
	//"encoding/csv"
	"flag"
	//"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime"
	//"strconv"
	//"strings"
	"warehouse-optimizer/warehouse"
)

const (
	gridPath = "warehouse-grid.csv"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	m := warehouse.ParseProductInfo(gridPath)
	pathInfo := warehouse.BuildPathInfo(gridPath)
	//fmt.Println("Hello User, where is your worker? e.g.:\"2 4\"")
	//x, y := warehouse.ReadInput()
	x, y := 0, 0
	start := warehouse.Point{X: x, Y: y}
	if x*y%2 == 1 {
		log.Fatal("Cannot start on a shelf.")
	}
	//fmt.Println("What is your worker's end location? e.g.:\"0 18\"")
	//x, y = warehouse.ReadInput()
	x, y = 0, 0
	end := warehouse.Point{X: x, Y: y}
	if x*y%2 == 1 {
		log.Fatal("Cannot end on a shelf.")
	}
	var t int
	for {
		/*fmt.Println("Type 1 to manual input, type 2 to file input.")
		var strInput string
		_, err := fmt.Scan(&strInput)
		if err != nil {
			log.Fatal(err)
		}
		strInput = strings.TrimSpace(strInput)
		t, err = strconv.Atoi(strInput)
		if err != nil {
			log.Fatal(err)
		}*/
		t = 2
		if t == 1 || t == 2 {
			break
		}
	}
	if t == 1 {
		prods := []int{46071, 379019, 70172, 1321, 2620261}
		for i:=0;i<10000000;i++{
			for _, prod := range prods {
				dest := warehouse.FindDest(start, m[prod])
				warehouse.FindPath(start, dest)
				warehouse.FindPath(dest, end)
			}
		}
		
		/*fmt.Println("Hello User, what items would you like to pick? (separate by space)")
		orders := warehouse.ReadOrder(m)
		fmt.Println("Here is the optimal picking order:")
		//optimalOrder := warehouse.BruteForceOrderOptimizer(orders[0], start, end, m, pathInfo)
		optimalOrder := warehouse.NearestNeighbourOrderOptimizer(orders[0], start, end, m, pathInfo)
		fmt.Println(optimalOrder)
		fmt.Println("Here is the optimal path:")
		s := warehouse.Route2String(optimalOrder, start, end, m)
		fmt.Println(s)
		fmt.Printf("Total distance traveled: %v\n", warehouse.RouteLength(optimalOrder, start, end, m, pathInfo))*/
	} else if t == 2 {
		//fmt.Println("Please list file of orders to be processed:")
		//ordersPath := warehouse.ReadString()
		ordersPath := "warehouse-orders-v01.csv"
		orders := warehouse.ParesOrderInfo(ordersPath)
		orderList := []int{9,11,68}
		for j:=0; j<100000; j++{
			for _, i := range orderList {
				warehouse.NearestNeighbourOrderOptimizer(orders[i], start, end, m, pathInfo)
			}
		}
		
		/*fmt.Println("Please list output file:")
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
			optimalOrder := warehouse.NearestNeighbourOrderOptimizer(order, start, end, m, pathInfo)
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
		fmt.Printf("%v orders processed.", orderCtr)*/
	}

	if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
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

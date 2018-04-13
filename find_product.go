package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"bufio"
)

const (
	gridPath     = "warehouse-grid.csv"
	shelfLength = 1.0
	shelfWidth  = 1.0
	pathWidthX  = 1.0
	pathWidthY  = 1.0
)

// Product defines the information of a product
type Product struct {
	id         int
	pos        Point
	l, r, u, d bool
	//num int
}

// Point defines the location of a point
type Point struct {
	x, y int
}

// Path is a slice of Points
type Path []Point

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

func coordinateConverter(x, y int) (int, int) {
	return 2*x + 1, 2*y + 1
}

func posAssigner(prod *Product) *Product {
	prod.l, prod.r = true, true
	return prod
}

// ParseProductInfo returns a map that includes product info
// TO-DO: ALSO FIND MAX/MIN INFO
// MAYBE NOT NECESSARY?
func ParseProductInfo(path string) map[int]Product {
	records, err := ReadCSV(path)
	if err != nil {
		log.Fatal(err)
	}
	var m map[int]Product
	m = make(map[int]Product)
	for _, s := range records {
		var temp [3]int
		var err error
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
		prod := Product{id: temp[0], pos: Point{temp[1], temp[2]}}
		m[temp[0]] = *posAssigner(&prod)
	}
	return m
}

// ParesOrderInfo returns a list of orders
func ParesOrderInfo(path string) [][]int {
	records, err := ReadCSV(path)
	if err != nil {
		log.Fatal(err)
	}
	var orders[][] int
	for _, s := range records {
		var err error
		s = strings.Split(strings.TrimSpace(s[0]), "\t")
		order := make([]int, len(s))
		for i := range s {
			s[i] = strings.TrimSpace(s[i])
			order[i], err = strconv.Atoi(s[i])
			if err != nil {
				log.Fatal(err)
			}
		}
		orders = append(orders, order)
	}
	return orders
}

// ReadOredr returns a list of "an" order to be compatible with ParesOrderInfo
// product_id should be separated by space from stdin
func ReadOrder() [][]int {
	r := bufio.NewReader(os.Stdin)
	strInput, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(strInput, " ")
	order := make([]int, len(s))
		for i := range s {
			s[i] = strings.TrimSpace(s[i])
			order[i], err = strconv.Atoi(s[i])
			if err != nil {
				log.Fatal(err)
			}
		}
	return [][]int{order}
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
func FindDest(src Point, prod Product) Point {
	if src.x < prod.pos.x {
		return Point{prod.pos.x - 1, prod.pos.y}
	}
	return Point{prod.pos.x + 1, prod.pos.y}
}

// FindPath returns the array of turning points on the path
// inclduing source and destination & length of the path
func FindPath(src Point, prod Product) (Path, float64) {
	dest := FindDest(src, prod)
	var path Path
	switch {
	case src.x == dest.x && src.y == dest.y:
	case src.x == dest.x:
		path = []Point{src, dest}
	case src.y%2 == 1 && src.y < dest.y:
		path = []Point{src, {src.x, src.y + 1}, {dest.x, src.y + 1}, dest}
	case src.y%2 == 1 && src.y >= dest.y:
		path = []Point{src, {src.x, src.y - 1}, {dest.x, src.y - 1}, dest}
	default:
		path = []Point{src, {dest.x, src.y}, dest}
	}
	return path, PathLength(path)
}

// PathLength returns the length of the path
func PathLength(path Path) float64 {
	if cap(path) < 1 {
		return 0.0
	}
	var dx, dy float64
	for i := range path[1:] {
		dx += math.Abs(float64(path[i+1].x-path[i].x)) * (shelfLength + pathWidthX) / 2
		dy += math.Abs(float64(path[i+1].y-path[i].y)) * (shelfWidth + pathWidthY) / 2
	}
	return dx + dy
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func (path Path) String() string {
	if cap(path) < 1 {
		return "Don't need to move."
	}
	s := fmt.Sprint(path[0])
	for _, p := range path[1:] {
		s += fmt.Sprintf("->%v", p)
	}
	return s
}

func main() {
	/*m := ParseProductInfo(gridPath)
	fmt.Print("Please input x y prod_id:\n")
	x, y, id := ReadInput()
	if x*y%2 == 1 {
		log.Fatal("Cannot start on a shelf.")
	}

	prod, ok := m[id]
	if ok {
		path, length := FindPath(Point{x, y}, prod)
		fmt.Printf("%v\nlength: %v\n", path, length)
	} else {
		fmt.Print("prod_id not exist.\n")
	}*/
	fmt.Print(ReadOrder())
}

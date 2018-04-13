package warehouse

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	gridPath    = "warehouse-grid.csv"
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
	X, Y int
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

// BuildPathInfo return a nested map that records the distances between Points
func BuildPathInfo(path string) map[Point]map[Point]float64 {
	var m map[Point]map[Point]float64
	m = make(map[Point]map[Point]float64)
	for i := 0; i <= 38; i++ {
		for j := 0; j <= 22; j++ {
			if i*j%2 == 0 {
				src := Point{i, j}
				//srcstr := fmt.Sprintf("(%v %v)", i, j)
				var m2 map[Point]float64
				m2 = make(map[Point]float64)
				for p := 0; p <= 38; p++ {
					for q := 0; q <= 22; q++ {
						if p*q%2 == 0 {
							dest := Point{p, q}
							//deststr := fmt.Sprintf("(%v %v)", p, q)
							_, m2[dest] = FindPath(src, dest)
						}
					}
				}
				m[src] = m2
			}
		}
	}
	return m
}

// ParesOrderInfo returns a list of orders
func ParesOrderInfo(path string) [][]int {
	records, err := ReadCSV(path)
	if err != nil {
		log.Fatal(err)
	}
	var orders [][]int
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

// ReadInput returns 2 int from stdin
func ReadInput() (int, int) {
	var strInput [2]string
	var input [2]int
	_, err := fmt.Scan(&strInput[0], &strInput[1])
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
	return input[0], input[1]
}

// FindDest returns the destination given init position & product to fetch
func FindDest(src Point, prod Product) Point {
	if src.X < prod.pos.X {
		return Point{prod.pos.X - 1, prod.pos.Y}
	}
	return Point{prod.pos.X + 1, prod.pos.Y}
}

// FindPath returns the array of turning points on the path
// inclduing source and destination & length of the path
func FindPath(src Point, dest Point) (Path, float64) {
	var path Path
	switch {
	case src.X == dest.X && src.Y == dest.Y:
	case src.X == dest.X:
		path = []Point{src, dest}
	case src.Y%2 == 1 && src.Y < dest.Y:
		path = []Point{src, {src.X, src.Y + 1}, {dest.X, src.Y + 1}, dest}
	case src.Y%2 == 1 && src.Y >= dest.Y:
		path = []Point{src, {src.X, src.Y - 1}, {dest.X, src.Y - 1}, dest}
	default:
		path = []Point{src, {dest.X, src.Y}, dest}
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
		dx += math.Abs(float64(path[i+1].X-path[i].X)) * (shelfLength + pathWidthX) / 2
		dy += math.Abs(float64(path[i+1].Y-path[i].Y)) * (shelfWidth + pathWidthY) / 2
	}
	return dx + dy
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
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

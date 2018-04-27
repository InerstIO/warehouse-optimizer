package warehouse

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cznic/mathutil"
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
	Pos        Point
	l, r, u, d bool
	pseudo     bool
	pseudoIn   Point
	pseudoOut  Point
	//num int
}

// Point defines the location of a point
type Point struct {
	X, Y int
}

// Path is a slice of Points
type Path []Point

// Order is a list of int that represents the products
type Order []int

func (o Order) Len() int           { return len(o) }
func (o Order) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Order) Less(i, j int) bool { return o[i] < o[j] }

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
		prod := Product{id: temp[0], Pos: Point{temp[1], temp[2]}, pseudo: false}
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
							m2[dest] = PathLength(FindPath(src, dest))
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
func ParesOrderInfo(path string) []Order {
	records, err := ReadCSV(path)
	if err != nil {
		log.Fatal(err)
	}
	var orders []Order
	for _, s := range records {
		var err error
		s = strings.Split(strings.TrimSpace(s[0]), "\t")
		order := make(Order, len(s))
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

// ReadOrder returns a list of "an" order to be compatible with ParesOrderInfo
// product_id should be separated by space from stdin
func ReadOrder(m map[int]Product) [][]int {
	r := bufio.NewReader(os.Stdin)
	strInput, err := r.ReadString('\n')
	strInput, err = r.ReadString('\n') // an ugly fix to avoid empty line from stdin
	strInput = strings.TrimSpace(strInput)
	if len(strInput) == 0 {
		log.Fatal("Empty input.")
	}
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
		_, ok := m[order[i]]
		if !ok {
			log.Fatalf("Item id %v not exist.", order[i])
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

// ReadString returns the string without space from stdin
func ReadString() string {
	var strInput string
	_, err := fmt.Scan(&strInput)
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(strInput)
}

func orderDeepCopy(o Order) Order {
	var newOrder Order
	for _, v := range o {
		newOrder = append(newOrder, v)
	}
	return newOrder
}

// BruteForceOrderOptimizer returns the Order with min total distance
func BruteForceOrderOptimizer(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) Order {
	var i sort.Interface = o
	mathutil.PermutationFirst(i)
	order := i.(Order)
	j := 0
	minIndex := j
	length := RouteLength(order, start, end, m, pathInfo)
	min := length
	for {
		ok := mathutil.PermutationNext(i)
		if !ok {
			break
		}
		order = i.(Order)
		length = RouteLength(order, start, end, m, pathInfo)
		j++
		if min > math.Min(min, length) {
			min = length
			minIndex = j
		}
	}

	mathutil.PermutationFirst(i)
	if minIndex == 0 {
		return i.(Order)
	}
	for ; minIndex > 0; minIndex-- {
		mathutil.PermutationNext(i)
	}
	return i.(Order)
}

// NearestNeighbourOrderOptimizer returns the Order by finding nearest neighbours
func NearestNeighbourOrderOptimizer(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) Order {
	var newOrder Order
	src := start
	for len(o) > 0 {
		minIndex := 0
		dest := FindDest(src, m[o[0]])
		length := pathInfo[src][dest]
		min := length
		minDest := dest
		for i, prod := range o[1:] {
			dest = FindDest(src, m[prod])
			length = pathInfo[src][dest]
			if min > math.Min(min, length) {
				min = length
				minIndex = i + 1
				minDest = dest
			}
		}
		newOrder = append(newOrder, o[minIndex])
		o = append(o[:minIndex], o[minIndex+1:]...)
		src = minDest
	}
	return newOrder
}

// NNIOrderOptimizer Nearest Neighbor With Iterations Order Optimizer.
// If no iteration varible given then iteration == len(order)
func NNIOrderOptimizer(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64, iteration ...int) Order {
	pseudoProd := Product{pseudo: true, pseudoIn: end, pseudoOut: start}
	iter := len(o)
	if len(iteration) > 0 {
		iter = iteration[0]
	}
	var newOrder Order
	minTotal := math.Inf(1)
	prods := []Product{pseudoProd}
	for _, p := range o {
		prods = append(prods, m[p])
	}
	for i := 0; i < iter; i++ {
		var src Point
		srcPoint := prods[i]
		ps := make([]Product, len(prods))
		copy(ps, prods)
		ps = append(ps[:i], ps[i+1:]...)
		if srcPoint.pseudo {
			src = srcPoint.pseudoOut
			nnOrder := nearestNeighborRing(ps, src, srcPoint, pathInfo)
			length := RouteLength(nnOrder, start, end, m, pathInfo)
			if length < minTotal {
				minTotal = length
				newOrder = nnOrder
			}
		} else {
			src = Point{srcPoint.Pos.X - 1, srcPoint.Pos.Y}
			nnOrder := nearestNeighborRing(ps, src, srcPoint, pathInfo)
			length := RouteLength(nnOrder, start, end, m, pathInfo)
			if length < minTotal {
				minTotal = length
				newOrder = nnOrder
			}
			src = Point{srcPoint.Pos.X + 1, srcPoint.Pos.Y}
			nnOrder = nearestNeighborRing(ps, src, srcPoint, pathInfo)
			length = RouteLength(nnOrder, start, end, m, pathInfo)
			if length < minTotal {
				minTotal = length
				newOrder = nnOrder
			}
		}
	}
	return newOrder
}

func nearestNeighborRing(prods []Product, src Point, srcProd Product, pathInfo map[Point]map[Point]float64) Order {
	ps := make([]Product, len(prods))
	copy(ps, prods)
	prodsOrder := []Product{srcProd}
	for len(ps) > 0 {
		minIndex := 0
		var length float64
		min := math.Inf(1)
		var newSrc Point
		for i, prod := range ps {
			dest := FindDest(src, prod)
			length = pathInfo[src][dest]
			if min > math.Min(min, length) {
				min = length
				minIndex = i
				if prod.pseudo {
					newSrc = prod.pseudoOut // for the new src
				} else {
					newSrc = dest
				}
			}
		}
		prodsOrder = append(prodsOrder, ps[minIndex])
		ps = append(ps[:minIndex], ps[minIndex+1:]...)
		src = newSrc
	}
	var startIndex int
	for i, prod := range prodsOrder {
		if prod.pseudo {
			startIndex = i
			break
		}
	}
	prodsOrder = append(prodsOrder[startIndex+1:], prodsOrder[:startIndex]...)
	var order Order
	for _, prod := range prodsOrder {
		order = append(order, prod.id)
	}
	return order
}

// RouteLength returns the length of the route for a specific Order
func RouteLength(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) float64 {
	var length float64
	var prevPos Point
	pos := FindDest(start, m[o[0]])
	prevPos = pos
	length += pathInfo[start][pos]
	for i := range o[1:len(o)] {
		prevPos = pos
		pos = FindDest(prevPos, m[o[i+1]])
		length += pathInfo[prevPos][pos]
	}
	length += pathInfo[pos][end]
	return length
}

// FindDest returns the destination given init position & product to fetch
func FindDest(src Point, prod Product) Point {
	if prod.pseudo {
		return prod.pseudoIn
	}
	if src.X < prod.Pos.X {
		return Point{prod.Pos.X - 1, prod.Pos.Y}
	}
	return Point{prod.Pos.X + 1, prod.Pos.Y}
}

// FindPath returns the array of turning points on the path
// inclduing source and destination
func FindPath(src Point, dest Point) Path {
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
	return path
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

// Route2String returns the string representation of the route
func Route2String(order Order, start, end Point, m map[int]Product) string {
	dest := FindDest(start, m[order[0]])
	s := fmt.Sprintf("%v->", FindPath(start, dest))
	s += fmt.Sprintf("[pick up %v from %v]->", order[0], m[order[0]].Pos)
	var src Point
	for _, prod := range order[1:] {
		src = dest
		dest = FindDest(start, m[prod])
		s += fmt.Sprintf("%v->", FindPath(src, dest))
		s += fmt.Sprintf("[pick up %v from %v]->", prod, m[prod].Pos)
	}
	src = dest
	s += fmt.Sprint(FindPath(src, end))
	return s
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

func (o Order) String() string {
	s := fmt.Sprint(o[0])
	for _, prod := range o[1:] {
		s += fmt.Sprintf(", %v", prod)
	}
	return s
}

// Order2csv returns a list of strings in csv compatible format
func Order2csv(o Order) []string {
	var ls []string
	for _, prod := range o {
		ls = append(ls, strconv.Itoa(prod))
	}
	return ls
}

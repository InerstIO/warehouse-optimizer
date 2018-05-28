package warehouse

import (
	"container/heap"
	"math"
	"time"
)

type vertex struct {
	matrix [][]float64
	cost   float64
	path   []int
}

type priorityQueue []*vertex

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	if pq[i].cost == pq[j].cost {
		return len(pq[i].path) > len(pq[j].path)
	}
	return pq[i].cost < pq[j].cost
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	v := x.(*vertex)
	*pq = append(*pq, v)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	v := old[n-1]
	*pq = old[0 : n-1]
	return v
}

func deepCopy2DMatrix(m [][]float64) [][]float64 {
	newMatrix := make([][]float64, len(m))
	for j := range newMatrix {
		newMatrix[j] = make([]float64, len(m[j]))
		copy(newMatrix[j], m[j])
	}
	return newMatrix
}

// reduceMatrix returns the reduced matrix and cost
func reduceMatrix(m [][]float64) ([][]float64, float64) {
	var cost float64
	newMatrix := deepCopy2DMatrix(m)
	for j := 0; j < len(newMatrix); j++ {
		min := math.Inf(1)
		for i := 0; i < len(newMatrix[j]); i++ {
			if newMatrix[j][i] < min {
				min = newMatrix[j][i]
				if min == 0.0 {
					break
				}
			}
		}
		if min == 0.0 || math.IsInf(min, 1) {
			continue
		}
		cost += min
		for i := 0; i < len(newMatrix[j]); i++ {
			newMatrix[j][i] -= min
		}
	}
	for i := 0; i < len(newMatrix[0]); i++ {
		min := math.Inf(1)
		for j := 0; j < len(newMatrix); j++ {
			if newMatrix[j][i] < min {
				min = newMatrix[j][i]
				if min == 0.0 {
					break
				}
			}
		}
		if min == 0.0 || math.IsInf(min, 1) {
			continue
		}
		cost += min
		for j := 0; j < len(newMatrix); j++ {
			newMatrix[j][i] -= min
		}
	}
	return newMatrix, cost
}

func explore(src vertex, dest int, m [][]float64, infSlice []float64) [][]float64 {
	newMatrix := deepCopy2DMatrix(m)
	newMatrix[src.path[len(src.path)-1]] = infSlice
	newMatrix[dest][src.path[len(src.path)-1]] = math.Inf(1)
	for j := 0; j < len(newMatrix); j++ {
		newMatrix[j][dest] = math.Inf(1)
	}
	for _, p := range src.path {
		newMatrix[dest][p] = math.Inf(1)
	}
	return newMatrix
}

func checkNext(dest int, parent *vertex, infSlice []float64) vertex {
	matrix, cost := reduceMatrix(explore(*parent, dest, parent.matrix, infSlice))
	newPath := make([]int, len(parent.path))
	copy(newPath, parent.path)
	return vertex{
		matrix: matrix,
		cost:   parent.cost + cost + parent.matrix[parent.path[len(parent.path)-1]][dest],
		path:   append(newPath, dest),
	}
}

func buildEdgeMatrixBnB(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) [][]float64 {
	prods := []Product{Product{Pos: start, pseudo: true, pseudoIn: end}}
	for _, p := range o {
		prod := m[p.prodID]
		prod.orderID = p.orderID
		prods = append(prods, prod)
	}
	matrix := make([][]float64, len(prods))
	for i := range matrix {
		matrix[i] = make([]float64, len(prods))
	}
	var length float64
	for j := range matrix {
		for i := 0; i < len(prods); i++ {
			if i != j {
				src := prods[j].Pos
				dest := FindDest(src, prods[i])
				if !prods[j].pseudo {
					src = FindDest(dest, prods[j]) // It will result in MST containing impossible edges
					// Since always choosing the smaller one from the left/right of the shelf
				}
				length = pathInfo[src][dest]
				matrix[j][i] = length
			} else {
				matrix[j][i] = math.Inf(1)
			}
		}
	}
	return matrix
}

func buildEdgeMatrixBnBLR(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) [][]float64 {
	prods := []Product{Product{Pos: start, pseudo: true, pseudoIn: end}}
	for _, p := range o {
		prod := m[p.prodID]
		prod.orderID = p.orderID
		prods = append(prods, prod)
	}
	matrix := make([][]float64, len(prods)*2-1)
	for i := range matrix {
		matrix[i] = make([]float64, len(prods)*2-1)
	}
	var length float64
	pm := []int{1, -1}
	for j := range matrix {
		for i := range matrix[j] {
			if (i+1)/2 != (j+1)/2 {
				src := prods[(j+1)/2].Pos
				if !prods[(j+1)/2].pseudo {
					src.X += pm[j%2]
				}
				dest := prods[(i+1)/2].Pos
				if !prods[(i+1)/2].pseudo {
					dest.X += pm[i%2]
				}
				length = pathInfo[src][dest]
				matrix[j][i] = length
			} else {
				matrix[j][i] = math.Inf(1)
			}
		}
	}
	return matrix
}

// BnBOrderOptimizer Branch and Bound Order Optimizer
func BnBOrderOptimizer(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64, timeLimit float64) Order {
	t := time.Now()
	matrix := buildEdgeMatrixBnB(o, start, end, m, pathInfo)
	infSlice := make([]float64, len(matrix[0]))
	for i := range infSlice {
		infSlice[i] = math.Inf(1)
	}
	var cost float64
	matrix, cost = reduceMatrix(matrix)
	initial := vertex{
		matrix: matrix,
		cost:   cost,
		path:   []int{0},
	}
	pq := priorityQueue{&initial}
	heap.Init(&pq)
	newOrder := NNIOrderOptimizer(o, start, end, m, pathInfo)
	min := reconstructCost(newOrder, o, start, end, m, pathInfo)
	realMin := RouteLength(newOrder, start, end, m, pathInfo)
	for pq.Len() > 0 {
		if time.Since(t).Seconds() > timeLimit {
			break
		}
		p := heap.Pop(&pq).(*vertex)
		var v vertex
		remain := 0
		if p.cost <= min {
			for i := range p.matrix {
				if math.IsInf(p.matrix[i][0], 1) {
					continue
				}
				remain++
				cv := checkNext(i, p, infSlice)
				if cv.cost <= min {
					heap.Push(&pq, &cv)
				}
				v = cv
			}
			if remain == 1 && v.cost <= min {
				min = v.cost
				var tempOrder Order
				for _, k := range v.path[1:] {
					tempOrder = append(tempOrder, o[k-1])
				}
				tempOrderLen := RouteLength(tempOrder, start, end, m, pathInfo)
				if tempOrderLen < realMin {
					realMin = tempOrderLen
					newOrder = tempOrder
				}
			}
		} else {
			break
		}
	}
	return newOrder
}

func (slice Order) pos(value Item) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func reconstructCost(o Order, ori Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) float64 {
	indices := make([]int, len(o))
	for i, item := range o {
		indices[i] = ori.pos(item)
	}
	matrix := buildEdgeMatrixBnB(ori, start, end, m, pathInfo)
	infSlice := make([]float64, len(matrix[0]))
	for i := range infSlice {
		infSlice[i] = math.Inf(1)
	}
	var cost float64
	matrix, cost = reduceMatrix(matrix)
	p := vertex{
		matrix: matrix,
		cost:   cost,
		path:   []int{0},
	}
	for _, i := range indices {
		p = checkNext(i+1, &p, infSlice)
	}
	return p.cost
}

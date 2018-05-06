package warehouse

import (
	"container/heap"
	"math"
)

type vertex struct {
	parent *vertex
	matrix [][]float64
	cost   float64
	path   []int
}

type priorityQueue []*vertex

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
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
	copy(newMatrix, m)
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
	return vertex{
		parent: parent,
		matrix: matrix,
		cost:   parent.cost + cost + parent.matrix[parent.path[len(parent.path)-1]][dest],
		path:   append(parent.path, dest),
	}
}

func buildEdgeMatrixBnB(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) [][]float64 {
	prods := []Product{Product{Pos: start, pseudo: true, pseudoIn: end}}
	for _, p := range o {
		prods = append(prods, m[p])
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

// BnBOrderOptimizer Branch and Bound Order Optimizer
func BnBOrderOptimizer(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) Order {
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
	min := math.Inf(1)
	realMin := math.Inf(1)
	var newOrder Order
	for pq.Len() > 0 {
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
				heap.Push(&pq, &cv)
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

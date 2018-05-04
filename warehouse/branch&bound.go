package warehouse

import (
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
		if min == 0.0 {
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
		if min == 0.0 {
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

func checkNext(src vertex, dest int, parent *vertex, infSlice []float64) vertex {
	matrix, cost := reduceMatrix(explore(src, dest, parent.matrix, infSlice))
	return vertex{
		parent: parent,
		matrix: matrix,
		cost:   parent.cost + cost + parent.matrix[parent.path[len(parent.path)-1]][dest],
		path:   append(parent.path, dest),
	}
}

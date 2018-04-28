package warehouse

import (
	"math"
)

// buildEdgeMatrix returns a 2D array with the all possible edge values in the order
func buildEdgeMatrix(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) [][]float64 {
	prods := []Product{Product{Pos: start, pseudo: true, pseudoIn: start}}
	if start != end {
		prods = append(prods, Product{Pos: end, pseudo: true, pseudoIn: end})
	}
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
				src := prods[i].Pos
				dest := FindDest(src, prods[j])
				if !prods[i].pseudo {
					src = FindDest(dest, prods[i]) // It will result in MST containing impossible edges
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

// lowerBoundGeneric returns the lower bound of the length of the route.
// start != end
func lowerBoundGeneric(matrix [][]float64) float64 {
	var minIndex int
	var sum float64
	var scanJ []int
	for j := range matrix {
		scanJ = append(scanJ, j)
	}
	for len(scanJ) > 0 {
		min := math.Inf(1)
		for k, j := range scanJ {
			for i := j + 1; i < len(matrix[j]); i++ {
				if matrix[j][i] < min {
					min = matrix[j][i]
					minIndex = k
				}
			}
		}
		if !math.IsInf(min, 1) {
			sum += min
		}
		scanJ = append(scanJ[:minIndex], scanJ[minIndex+1:]...)
	}
	return sum
}

// lowerBoundSE returns the lower bound of the length of the route.
// start == end. Running iterations.
func lowerBoundSE(matrix [][]float64) float64 {
	maxsum := math.Inf(-1)
	var sum float64
	if len(matrix) > 2 {
		for i := range matrix {
			newMatrix := make([][]float64, len(matrix))
			copy(newMatrix, matrix)
			newMatrix = append(newMatrix[:i], newMatrix[i+1:]...)
			sum = lowerBoundGeneric(newMatrix)
			min1, min2 := math.Inf(1), math.Inf(1)
			for _, num := range matrix[i] {
				if num < min1 {
					min1 = num
					min2 = min1
				}
			}
			sum += min1 + min2
			if sum > maxsum {
				maxsum = sum
			}
		}
	} else if len(matrix) == 2 {
		maxsum = matrix[0][1] * 2
	} else {
		maxsum = 0
	}
	return maxsum
}

// LowerBound returns the lower bound of the length of the route
func LowerBound(o Order, start, end Point, m map[int]Product, pathInfo map[Point]map[Point]float64) float64 {
	matrix := buildEdgeMatrix(o, start, end, m, pathInfo)
	if start == end {
		return lowerBoundSE(matrix)
	}
	return lowerBoundGeneric(matrix)
}

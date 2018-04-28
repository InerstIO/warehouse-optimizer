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

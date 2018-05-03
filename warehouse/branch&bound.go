package warehouse

import (
	"math"
)

// ReduceMatrix returns the reduced matrix and cost
func ReduceMatrix(m [][]float64) ([][]float64, float64) {
	var cost float64
	newMatrix := make([][]float64, len(m))
	copy(newMatrix, m)
	for j := range newMatrix {
		newMatrix[j] = make([]float64, len(m[j]))
		copy(newMatrix[j], m[j])
	}
	for j := 0; j<len(newMatrix); j++{
		min := math.Inf(1)
		for i:= 0; i<len(newMatrix[j]); i++{
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
		for i:= 0; i<len(newMatrix[j]); i++ {
			newMatrix[j][i] -= min
		}
	}
	for i:= 0; i<len(newMatrix[0]); i++ {
		min := math.Inf(1)
		for j:= 0;j<len(newMatrix); j++ {
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
		for j:=0; j<len(newMatrix); j++{
			newMatrix[j][i] -= min
		}
	}
	return newMatrix, cost
}

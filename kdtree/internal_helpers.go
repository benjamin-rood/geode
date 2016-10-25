package kdtree

import "math/rand"

func randomFloatInRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func sum(set []float64) float64 {
	var result float64
	for i := range set {
		result += set[i]
	}
	return result
}

func sumValuesAlongAxis(ds Datapoints, axis int) float64 {
	var sum float64
	for i := range ds {
		sum += ds[i].set[axis]
	}
	return sum
}

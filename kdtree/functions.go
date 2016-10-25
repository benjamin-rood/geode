package kdtree

import "math"

// Distance returns the Euclidean length of the line connecting any two Datapoints
func Distance(p, q *Datapoint) float64 {
	return math.Sqrt(DistanceSq(p, q))
}

// DistanceSq returns the length squared of the line connecting any two Datapoints
func DistanceSq(p, q *Datapoint) float64 {
	var differences = make([]float64, len(p.set), len(p.set))
	for i := range p.set {
		v := q.set[i] - p.set[i]
		differences[i] = v * v
	}
	return sum(differences)
}

package kdtree

import "math"

// Convert uses the Importable interface to cleanly produce a kdtree
// from a slice of some type which has implented ToDataPoint(), and
// according to the pivot algorithm (PivotFunc).
func Convert(c []Importable, sorting bool, pivotDef PivotFunc) (*Branch, error) {
	var points = make(Datapoints, len(c), len(c))
	// basedim := len(c[0].ToDatapoint().set)
	for i := range c {
		points[i] = c[i].ToDatapoint()
		// TODO:
		// Implement build error types.
		// e.g.
		// if points[i].Dimensionality() != basedim {
		// 	return nil, DimClashError
		// }
	}

	if pivotDef == nil {
		pivotDef = LazyAverage
	}

	b := Build(points, 0, pivotDef)
	return b, nil
}

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

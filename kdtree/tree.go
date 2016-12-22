package kdtree

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Branch is a Binary Tree Node
type Branch struct {
	Datapoints
	pivot       float64
	depth       int
	left, right *Branch
}

// PivotFunc calculates the pivot value
type PivotFunc func(Datapoints, int) float64

// Set of pre-defined functions which match the prototype of `PivotFunc`
var (
	// LazyAverage implements a simple fast average split to produce the pivot value
	// Sufficient for large numbers of Datapoints with uniformly distributed values
	LazyAverage = func(ds Datapoints, axis int) float64 {
		sz := len(ds)
		first := ds[0].set[axis]
		last := ds[sz-1].set[axis]
		pivot := (first + last) / 2
		return pivot
	}

	// Median implements a true median split on a sorted set for the pivot value
	// note: will be significantly slower
	Median = func(ds Datapoints, axis int) float64 {
		By(Comparator(axis)).Sort(ds)
		midpoint := len(ds) / 2
		return ds[midpoint].set[axis]
	}

	// Mean implements a true mean (average) calculation to determine the pivot value
	Mean = func(ds Datapoints, axis int) float64 {
		sz := float64(len(ds))
		return sumValuesAlongAxis(ds, axis) / sz
	}
)

// Build constructs the k-d tree from a set of assumed to be valid Datapoints
// OF CONSISTENT DIMENSIONALITY, using a provided PivotFunc algorithm
func Build(ds Datapoints, depth int, pivotDef PivotFunc) *Branch {
	if ds == nil {
		return nil
	}

	sz := len(ds)
	if sz <= 1 {
		return &Branch{ds[:1], 0, depth, nil, nil}
	}
	if ds.notDistinct() {
		return &Branch{ds, 0, depth, nil, nil}
	}

	if pivotDef == nil {
		pivotDef = LazyAverage
	}

	branch := Branch{
		Datapoints: ds,
		pivot:      0,
		depth:      depth,
		left:       nil,
		right:      nil,
	}

	dimensionality := len(branch.Datapoints[0].set)
	axis := depth % dimensionality
	branch.pivot = pivotDef(branch.Datapoints, axis)

	leftSet, rightSet := make(Datapoints, 0, sz), make(Datapoints, 0, sz)

	for i := range branch.Datapoints {
		if branch.Datapoints[i].set[axis] < branch.pivot {
			leftSet = append(leftSet, branch.Datapoints[i])
		} else {
			rightSet = append(rightSet, branch.Datapoints[i])
		}
	}

	branch.left = Build(leftSet, depth+1, pivotDef)
	branch.right = Build(rightSet, depth+1, pivotDef)
	return &branch
}

// MaxDepth returns the depth of the deepest leaf node from the input branch as 'root'
func (branch *Branch) MaxDepth() int {
	if branch == nil {
		return 0
	}
	return max(branch.depth, max(branch.left.MaxDepth(), branch.right.MaxDepth()))
}

// ANN will very rapidly return the **approximate nearest neighbour** Datapoint
// in a given k-d tree branch.
// If we consider the accuracy of ANN as the spatial distance, d, from the exact
// nearest neighbour to the Datapoint returned by ANN (where d=0 implies
// ANN(branch, target) = NN(branch, target)), the accuracy of ANN increases
// dramatically with the density of the branch passed to ANN.
// ANN achieves an extremely high degree of accuracy when the density of the points
// in each axis > 100,000; where density is defined as the number of leaves / (max-min).
func ANN(branch *Branch, target *Datapoint) *Datapoint {
	sz := len(branch.Datapoints)
	if sz == 1 {
		return branch.Datapoints[0]
	}
	if branch.Datapoints.notDistinct() {
		return branch.Datapoints[rand.Intn(sz)] // pick a pseudorandom point in the range.
	}

	dimensionality := len(branch.Datapoints[0].set)
	axis := branch.depth % dimensionality
	comparator := target.set[axis]

	if comparator < branch.pivot {
		branch = branch.left
	} else {
		branch = branch.right
	}

	return ANN(branch, target)
}

func areaN(branch *Branch, target *Datapoint, granularity int) Datapoints {
	if len(branch.Datapoints) <= granularity || branch.Datapoints.notDistinct() {
		return branch.Datapoints
	}

	dimensionality := len(branch.Datapoints[0].set)
	axis := branch.depth % dimensionality
	comparator := target.set[axis]

	if comparator < branch.pivot {
		branch = branch.left
	} else {
		branch = branch.right
	}

	bin := Datapoints{}
	bin = append(bin, areaN(branch, target, granularity)...)
	return bin
}

// NN returns the **exact** nearest-neighbouring Datapoint in the k-d tree branch.
// As the process is approximately 3X slower than ANN, unless you explicitly require
// the exact nearest nerighbour to the target, use ANN instead in most cases.
func NN(branch *Branch, target *Datapoint) *Datapoint {
	bin := areaN(branch, target, 10)
	best := bin[0]
	for i := 1; i < len(bin); i++ {
		if DistanceSq(target, bin[i]) < DistanceSq(target, best) {
			best = bin[i]
		}
	}
	return best
}

func inRange(xmin, xmax, lo, hi float64) bool {
	return xmin >= lo && xmax <= hi
}

// Range holds lower and upper range values
type Range struct {
	min, max float64
}

// RangeQuery returns all Datapoints in a specified bounded area
func RangeQuery(branch *Branch, bounds []Range) Datapoints {
	if branch == nil {
		return nil
	}

	dimensionality := len(branch.Datapoints[0].set)
	last := len(branch.Datapoints) - 1
	intersection := false

	if len(branch.Datapoints) <= 5 {
		for axis := 0; axis < dimensionality; axis++ {
			By(Comparator(axis)).Sort(branch.Datapoints)
			intersection = inRange(
				branch.Datapoints[0].set[axis],
				branch.Datapoints[last].set[axis],
				bounds[axis].min,
				bounds[axis].max,
			)
			if !intersection {
				break
			}
		}
	}
	if intersection {
		return branch.Datapoints
	}

	axis := branch.depth % dimensionality

	var rangeSet, leftSet, rightSet Datapoints
	if branch.pivot > bounds[axis].min { // continue tree traversal left
		leftSet = RangeQuery(branch.left, bounds)
	}
	if branch.pivot <= bounds[axis].max {
		rightSet = RangeQuery(branch.right, bounds)
	}

	rangeSet = append(rangeSet, leftSet...)
	rangeSet = append(rangeSet, rightSet...)
	return rangeSet
}

// MarshalJSON implements json.Marshaler interface
func (branch *Branch) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Depth":       branch.depth,
		"Cardinality": len(branch.Datapoints),
		"Datapoints":  branch.Datapoints,
		"Pivot":       branch.pivot,
		"leftChild":   branch.left,
		"rightChild":  branch.right,
	})
}

func buildDebug(ds Datapoints, depth int, pivotDef PivotFunc) *Branch {
	if ds == nil {
		return nil
	}
	sz := len(ds)
	fmt.Println(sz)
	fmt.Println(ds.PointsSetString())
	time.Sleep(250 * time.Millisecond)
	if sz <= 1 {
		return &Branch{ds[:1], 0, depth, nil, nil}
	}
	if ds.notDistinct() {
		return &Branch{ds, 0, depth, nil, nil}
	}

	if pivotDef == nil {
		pivotDef = LazyAverage
	}

	branch := Branch{
		Datapoints: ds,
		pivot:      0,
		depth:      depth,
		left:       nil,
		right:      nil,
	}

	dimensionality := len(branch.Datapoints[0].set)
	axis := depth % dimensionality
	branch.pivot = pivotDef(branch.Datapoints, axis)
	fmt.Println(`branch.pivot =`, branch.pivot)

	leftSet, rightSet := make(Datapoints, 0, sz), make(Datapoints, 0, sz)

	for i := range branch.Datapoints {
		if branch.Datapoints[i].set[axis] < branch.pivot {
			leftSet = append(leftSet, branch.Datapoints[i])
		} else {
			rightSet = append(rightSet, branch.Datapoints[i])
		}
	}

	branch.left = buildDebug(leftSet, depth+1, pivotDef)
	branch.right = buildDebug(rightSet, depth+1, pivotDef)
	return &branch
}

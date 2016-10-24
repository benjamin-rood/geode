package kdtree

// Branch is a Binary Tree Node
type Branch struct {
	Datapoints
	pivot       float64
	depth       int
	left, right *Branch
}

// Build constructs the kdtree from a set of assumed to be valid Datapoints
func Build(ds Datapoints, depth int) *Branch {
	branch := Branch{
		Datapoints: ds,
		pivot:      0,
		left:       nil,
		right:      nil,
	}
	sz := len(ds)
	if sz <= 1 {
		return &branch
	}

	dimensionality := len(ds[0].set)
	axis := depth % dimensionality

	first := branch.Datapoints[0].set[axis]
	last := branch.Datapoints[sz-1].set[axis]

	branch.pivot = (first + last) / 2

	var leftSet, rightSet Datapoints

	for i := range ds {
		switch ds[i].set[axis] < branch.pivot {
		case true:
			leftSet = append(leftSet, ds[i])
		case false:
			rightSet = append(rightSet, ds[i])
		}
	}

	branch.left = Build(leftSet, depth+1)
	branch.right = Build(rightSet, depth+1)
	return &branch
}

// ANN will very rapidly return the **approximate nearest neighbour** Datapoint
// in a given k-d tree branch.
// If we consider the accuracy of ANN as the spatial distance, d, from the exact
// nearest neighbour to the Datapoint returned by ANN (where d=0 implies
// ANN(branch, target) = NN(branch, target)), the accuracy of ANN increases
// dramatically with the density of the branch passed to ANN.
// ANN achieves an extremely high degree of accuracy when the density of the points
// in each axis > 100,000; where density is defined as the number of leaves / (max-min)).
func ANN(branch *Branch, target *Datapoint) *Datapoint {
	if len(branch.Datapoints) == 1 {
		return branch.Datapoints[0]
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
	if len(branch.Datapoints) <= granularity {
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
// As the process is approximately 3X slower than ANN, unless you explictly require
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

// RangeQuery returns all Datapoints in a specified bounded area
func RangeQuery(branch *Branch,
	bounds []struct {
		lo, hi float64
	}) Datapoints {
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
				bounds[axis].lo,
				bounds[axis].hi,
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
	if branch.pivot > bounds[axis].lo { // continue tree traversal left
		leftSet = RangeQuery(branch.left, bounds)
	}
	if branch.pivot <= bounds[axis].hi {
		rightSet = RangeQuery(branch.right, bounds)
	}

	rangeSet = append(rangeSet, leftSet...)
	rangeSet = append(rangeSet, rightSet...)
	return rangeSet
}

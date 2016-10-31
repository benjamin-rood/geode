package kdtree

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

func randomFloatInRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
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

const (
	void = `()`
	sc   = `;`
)

func branchToStrMatrix(branch *Branch) [][]string {
	height, width := branch.MaxDepth(), len(branch.Datapoints)
	matrix := make([][]string, width, width)
	for i := 0; i < height; i++ {
		matrix[i] = append(matrix[i], make([]string, height, height)...)
	}

	return matrix
}

type pivotList map[int]string

func (pL pivotList) MarshalJSON() ([]byte, error) {
	max := 0
	for key := range pL {
		if key > max {
			max = key
		}
	}
	representation := make(map[string]interface{})
	for key, value := range pL {
		if key == max {
			continue
		}
		rKey := fmt.Sprintf("depth=%v (pivots)", key)
		representation[rKey] = value
	}
	representation["leaves"] = pL[max]
	return json.Marshal(representation)
}

func breadthFirstSearchPivotsList(branch *Branch) pivotList {
	md := branch.MaxDepth() + 1
	pL := make(pivotList)
	queue := []Branch{} // LIFO
	b := *branch
	queue = append(queue, b) // add root to queue tail
	for len(queue) != 0 {
		b, queue = queue[0], queue[1:] // pop head
		if _, exists := pL[b.depth]; !exists {
			pL[b.depth] = "{"
		}
		pL[b.depth] += (fmt.Sprint(b.pivot) + ", ")
		if b.left != nil {
			queue = append(queue, *b.left) // add to tail
		}
		if b.right != nil {
			queue = append(queue, *b.right) // add to tail
		}
	}
	for key, value := range pL {
		pL[key] = (value[:len(value)-2] + "}")
	}

	dsQ := depthFirstSearchLeavesOnly(branch)
	pL[md] = dsQ.PointsSetString()
	return pL
}

func depthFirstSearchLeavesOnly(branch *Branch) Datapoints {
	var dsQ Datapoints   // LIFO queue
	stack := []*Branch{} // FIFO
	b := branch
	stack = append(stack, b)
	for len(stack) != 0 {
		b, stack = stack[len(stack)-1], stack[:len(stack)-1] // pop
		if b == nil {
			continue
		}
		if b.left == nil && b.right == nil {
			temp := Datapoints{b.Datapoints[0]}
			dsQ = append(temp, dsQ...) // append to front
		} else {
			if b.left != nil {
				stack = append(stack, b.left)
			}
			if b.right != nil {
				stack = append(stack, b.right)
			}
		}
	}
	return dsQ
}

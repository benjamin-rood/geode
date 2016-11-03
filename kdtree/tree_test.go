package kdtree

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"testing"
)

func Test_Tree_Branch_json_Marshaller_Interface(t *testing.T) {
	tree := Build(dps1, 0, Median)
	jsonTree, err := json.MarshalIndent(tree, "", "    ")
	if err != nil {
		t.Error(err)
	}
	want, err := ioutil.ReadFile("test_fixtures/dps1_median.json")
	if err != nil {
		t.Error(err)
	}
	if string(jsonTree) != string(want) {
		t.Error(`want: `, string(want), `
		got: `, string(jsonTree))
	}

	var dps2 Datapoints
	for i := 0; i < 20; i++ {
		dps2 = append(dps2, &Datapoint{nil, []float64{float64(rand.Intn(25)), float64(rand.Intn(25))}})
	}

	tree = Build(dps2, 0, Median)
	jsonTree, err = json.MarshalIndent(tree, "", "    ")
	if err != nil {
		t.Error(err)
	}
	want, err = ioutil.ReadFile("test_fixtures/dps2_median.json")
	if err != nil {
		t.Error(err)
	}
	if string(jsonTree) != string(want) {
		t.Error(`got: `, string(jsonTree))
	}

	tree = Build(dps3, 0, Median)
	jsonTree, err = json.MarshalIndent(tree, "", "    ")
	if err != nil {
		t.Error(err)
	}
	want, err = ioutil.ReadFile("test_fixtures/dps3_median.json")
	if err != nil {
		t.Error(err)
	}
	if string(jsonTree) != string(want) {
		t.Error(`got: `, string(jsonTree))
	}
}

func Test_Tree_Branch_MaxDepth(t *testing.T) {
	tree := Build(dps1, 0, nil)
	want := 3
	got := tree.MaxDepth()
	if got != want {
		t.Error(`want: `, want, `
		got: `, got)
	}
}

func Test_Tree_Branch_Build_Pivot_Mean(t *testing.T) {
	tree := Build(dps3, 0, Mean)
	want, err := ioutil.ReadFile("test_fixtures/branch_build_pivot_mean.json")
	if err != nil {
		t.Error(err)
	}
	pL := breadthFirstSearchPivotsList(tree)
	got, err := json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want: `, string(want), `
		got: `, string(got))
	}
}

func Test_Tree_Branch_Build_Pivot_Median(t *testing.T) {
	tree := Build(dps3, 0, Median)
	want, err := ioutil.ReadFile("test_fixtures/branch_build_pivot_median.json")
	if err != nil {
		t.Error(err)
	}
	pL := breadthFirstSearchPivotsList(tree)
	got, err := json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want: `, string(want), `
		got: `, string(got))
	}
}

func Test_Tree_Branch_Build_Pivot_LazyAverage(t *testing.T) {
	tree := Build(dps3, 0, LazyAverage)
	want, err := ioutil.ReadFile("test_fixtures/branch_build_pivot_lazyaverage.json")
	if err != nil {
		t.Error(err)
	}
	pL := breadthFirstSearchPivotsList(tree)
	got, err := json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want: `, string(want), `
		got: `, string(got))
	}
}

func Test_Tree_Branch_Build_NonDistinctDatapoints(t *testing.T) {
	tree := Build(nonDistinctDps, 0, LazyAverage)
	want, err := ioutil.ReadFile("test_fixtures/nonDistinctDatapoints_lazyaverage.json")
	if err != nil {
		t.Error(err)
	}
	pL := breadthFirstSearchPivotsList(tree)
	got, err := json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(LazyAverage): `, string(want), `
		got(LazyAverage): `, string(got))
	}

	tree = Build(nonDistinctDps, 0, Mean)
	want, err = ioutil.ReadFile("test_fixtures/nonDistinctDatapoints_mean.json")
	if err != nil {
		t.Error(err)
	}
	pL = breadthFirstSearchPivotsList(tree)
	got, err = json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(Mean): `, string(want), `
		got(Mean): `, string(got))
	}

	tree = Build(nonDistinctDps, 0, Median)
	want, err = ioutil.ReadFile("test_fixtures/nonDistinctDatapoints_median.json")
	if err != nil {
		t.Error(err)
	}
	pL = breadthFirstSearchPivotsList(tree)
	got, err = json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(Median): `, string(want), `
		got(Median): `, string(got))
	}
}

func Test_Tree_Branch_Build_SingleDimensionDatapoints(t *testing.T) {
	tree := Build(singleDimDps, 0, LazyAverage)
	want, err := ioutil.ReadFile("test_fixtures/singleDimDatapoints_lazyaverage.json")
	if err != nil {
		t.Error(err)
	}
	pL := breadthFirstSearchPivotsList(tree)
	got, err := json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(LazyAverage): `, string(want), `
		got(LazyAverage): `, string(got))
	}

	tree = Build(singleDimDps, 0, Mean)
	want, err = ioutil.ReadFile("test_fixtures/singleDimDatapoints_mean.json")
	if err != nil {
		t.Error(err)
	}
	pL = breadthFirstSearchPivotsList(tree)
	got, err = json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(Mean): `, string(want), `
		got(Mean): `, string(got))
	}

	tree = Build(singleDimDps, 0, Median)
	want, err = ioutil.ReadFile("test_fixtures/singleDimDatapoints_median.json")
	if err != nil {
		t.Error(err)
	}
	pL = breadthFirstSearchPivotsList(tree)
	got, err = json.MarshalIndent(pL, "", "    ")
	if err != nil {
		t.Error(err)
	}
	if string(got) != string(want) {
		t.Error(`want(Median): `, string(want), `
		got(Median): `, string(got))
	}
}

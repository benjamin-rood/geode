package kdtree

import (
	"encoding/json"
	"math/rand"
	"testing"
)

const (
	jsonDps1MedianStr = `{"Cardinality":5,"Datapoints":[{"data":null,"set":[1,2]},{"data":null,"set":[2,3]},{"data":null,"set":[3,4]},{"data":null,"set":[4,5]},{"data":null,"set":[5,6]}],"Depth":0,"Pivot":3,"leftChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[1,2]},{"data":null,"set":[2,3]}],"Depth":1,"Pivot":3,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[1,2]}],"Depth":2,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[2,3]}],"Depth":2,"Pivot":0,"leftChild":null,"rightChild":null}},"rightChild":{"Cardinality":3,"Datapoints":[{"data":null,"set":[3,4]},{"data":null,"set":[4,5]},{"data":null,"set":[5,6]}],"Depth":1,"Pivot":5,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[3,4]}],"Depth":2,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[4,5]},{"data":null,"set":[5,6]}],"Depth":2,"Pivot":5,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[4,5]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[5,6]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null}}}}`
	jsonDps2MedianStr = `{"Cardinality":11,"Datapoints":[{"data":null,"set":[0.06563701921747622,0.15651925473279124]},{"data":null,"set":[0.09696951891448456,0.30091186058528707]},{"data":null,"set":[0.20318687664732285,0.360871416856906]},{"data":null,"set":[0.21426387258237492,0.380657189299686]},{"data":null,"set":[0.28303415118044517,0.29310185733681576]},{"data":null,"set":[0.31805817433032985,0.4688898449024232]},{"data":null,"set":[0.4246374970712657,0.6868230728671094]},{"data":null,"set":[0.5152126285020654,0.8136399609900968]},{"data":null,"set":[0.6046602879796196,0.9405090880450124]},{"data":null,"set":[0.6645600532184904,0.4377141871869802]},{"data":null,"set":[0.6790846759202163,0.21855305259276428]}],"Depth":0,"Pivot":0.31805817433032985,"leftChild":{"Cardinality":5,"Datapoints":[{"data":null,"set":[0.06563701921747622,0.15651925473279124]},{"data":null,"set":[0.28303415118044517,0.29310185733681576]},{"data":null,"set":[0.09696951891448456,0.30091186058528707]},{"data":null,"set":[0.20318687664732285,0.360871416856906]},{"data":null,"set":[0.21426387258237492,0.380657189299686]}],"Depth":1,"Pivot":0.30091186058528707,"leftChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[0.06563701921747622,0.15651925473279124]},{"data":null,"set":[0.28303415118044517,0.29310185733681576]}],"Depth":2,"Pivot":0.28303415118044517,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.06563701921747622,0.15651925473279124]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.28303415118044517,0.29310185733681576]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null}},"rightChild":{"Cardinality":3,"Datapoints":[{"data":null,"set":[0.09696951891448456,0.30091186058528707]},{"data":null,"set":[0.20318687664732285,0.360871416856906]},{"data":null,"set":[0.21426387258237492,0.380657189299686]}],"Depth":2,"Pivot":0.20318687664732285,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.09696951891448456,0.30091186058528707]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[0.20318687664732285,0.360871416856906]},{"data":null,"set":[0.21426387258237492,0.380657189299686]}],"Depth":3,"Pivot":0.380657189299686,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.20318687664732285,0.360871416856906]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.21426387258237492,0.380657189299686]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null}}}},"rightChild":{"Cardinality":6,"Datapoints":[{"data":null,"set":[0.6790846759202163,0.21855305259276428]},{"data":null,"set":[0.6645600532184904,0.4377141871869802]},{"data":null,"set":[0.31805817433032985,0.4688898449024232]},{"data":null,"set":[0.4246374970712657,0.6868230728671094]},{"data":null,"set":[0.5152126285020654,0.8136399609900968]},{"data":null,"set":[0.6046602879796196,0.9405090880450124]}],"Depth":1,"Pivot":0.6868230728671094,"leftChild":{"Cardinality":3,"Datapoints":[{"data":null,"set":[0.31805817433032985,0.4688898449024232]},{"data":null,"set":[0.6645600532184904,0.4377141871869802]},{"data":null,"set":[0.6790846759202163,0.21855305259276428]}],"Depth":2,"Pivot":0.6645600532184904,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.31805817433032985,0.4688898449024232]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[0.6790846759202163,0.21855305259276428]},{"data":null,"set":[0.6645600532184904,0.4377141871869802]}],"Depth":3,"Pivot":0.4377141871869802,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.6790846759202163,0.21855305259276428]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.6645600532184904,0.4377141871869802]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null}}},"rightChild":{"Cardinality":3,"Datapoints":[{"data":null,"set":[0.4246374970712657,0.6868230728671094]},{"data":null,"set":[0.5152126285020654,0.8136399609900968]},{"data":null,"set":[0.6046602879796196,0.9405090880450124]}],"Depth":2,"Pivot":0.5152126285020654,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.4246374970712657,0.6868230728671094]}],"Depth":3,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":2,"Datapoints":[{"data":null,"set":[0.5152126285020654,0.8136399609900968]},{"data":null,"set":[0.6046602879796196,0.9405090880450124]}],"Depth":3,"Pivot":0.9405090880450124,"leftChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.5152126285020654,0.8136399609900968]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null},"rightChild":{"Cardinality":1,"Datapoints":[{"data":null,"set":[0.6046602879796196,0.9405090880450124]}],"Depth":4,"Pivot":0,"leftChild":null,"rightChild":null}}}}}`
)

var (
	dps1 = Datapoints{
		&Datapoint{nil, []float64{1, 2}},
		&Datapoint{nil, []float64{2, 3}},
		&Datapoint{nil, []float64{3, 4}},
		&Datapoint{nil, []float64{4, 5}},
		&Datapoint{nil, []float64{5, 6}},
	}
	dps2 = Datapoints{
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
		&Datapoint{nil, []float64{rand.Float64(), rand.Float64()}},
	}
)

func Test_Branch_json_Marshaller_Interface(t *testing.T) {
	tree := Build(dps1, 0, Median)
	jsonTree, _ := json.Marshal(tree)
	if string(jsonTree) != jsonDps1MedianStr {
		t.Error(string(jsonTree))
	}
	tree = Build(dps2, 0, Median)
	jsonTree, _ = json.Marshal(tree)
	if string(jsonTree) != jsonDps2MedianStr {
		t.Error(string(jsonTree))
	}
}
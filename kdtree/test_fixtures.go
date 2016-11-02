package kdtree

import (
	"encoding/json"
	"math"
	"math/big"
)

var (
	dps1 = Datapoints{
		&Datapoint{nil, []float64{1, 2}},
		&Datapoint{nil, []float64{2, 3}},
		&Datapoint{nil, []float64{3, 4}},
		&Datapoint{nil, []float64{4, 5}},
		&Datapoint{nil, []float64{5, 6}},
	}

	dps3 = Datapoints{
		&Datapoint{nil, []float64{1, 9}},
		&Datapoint{nil, []float64{2, 3}},
		&Datapoint{nil, []float64{4, 1}},
		&Datapoint{nil, []float64{3, 7}},
		&Datapoint{nil, []float64{5, 4}},
		&Datapoint{nil, []float64{6, 8}},
		&Datapoint{nil, []float64{7, 2}},
		&Datapoint{nil, []float64{8, 8}},
		&Datapoint{nil, []float64{7, 9}},
		&Datapoint{nil, []float64{9, 6}},
	}
)

type intAndFloat struct {
	i int
	f float64
}

func (inf intAndFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"i": inf.i,
		"f": inf.f,
	})
}

type embedTest struct {
	s string
	intAndFloat
}

func (emb embedTest) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"s":           emb.s,
		"intAndFloat": emb.intAndFloat,
	})
}

var (
	diffFptrs = func(f, g []*float64) bool {
		if len(f) != len(g) {
			return false
		}
		for i := range f {
			if f[i] == g[i] {
				return false
			}
		}
		return true
	}

	getF64ptrs = func(f []float64) (f64ptrs []*float64) {
		for i := range f {
			f64ptrs = append(f64ptrs, &f[i])
		}
		return
	}

	embedded = embedTest{"hello world", intAndFloat{0, math.E}}

	rational = big.NewRat(5, 4)

	str = "cassandra"

	mappy = map[string]interface{}{
		"a": int('a'),
		"b": string("banana"),
		"c": intAndFloat{0, math.Pi},
	}

	interfaces = []interface{}{&str, str, &rational, rational, &embedded, embedded, &mappy, mappy}

	constructorInputs = []struct {
		linked interface{}
		values []float64
	}{
		{
			&str,
			[]float64{6.0000125, 6.10000125, -1.3173, 1373},
		},
		{
			str,
			[]float64{6.0000125, 6.10000125, -1.3173, 1373},
		},
		{
			&rational,
			[]float64{1, 2, 3, 4, 5},
		},
		{
			rational,
			[]float64{1, 2, 3, 4, 5},
		},
		{
			&embedded,
			[]float64{100.3},
		},
		{
			embedded,
			[]float64{100.3},
		},
		{
			&mappy,
			[]float64{0.8050908121798804, 0.53238545404102},
		},
		{
			mappy,
			[]float64{0.8050908121798804, 0.53238545404102},
		},
	}
)

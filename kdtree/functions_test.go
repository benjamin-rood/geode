package kdtree

import (
	"math"
	"testing"
)

type dpPairsCmp struct {
	A, B *Datapoint
	want float64
}

var (
	pointPairs = []dpPairsCmp{
		{
			A:    &Datapoint{nil, []float64{0, 0}},
			B:    &Datapoint{nil, []float64{0, 1}},
			want: 0,
		},
		{
			A:    &Datapoint{nil, []float64{3, 3, 3}},
			B:    &Datapoint{nil, []float64{-2, 7.5, 0.125}},
			want: 0,
		},
	}
)

func Test_Func_Eucliean_Distance_Btwn_Datapoints(t *testing.T) {
	euclideanTests := pointPairs
	euclideanTests[0].want = math.Sqrt(0*0 + 1*1)
	euclideanTests[1].want = math.Sqrt(math.Pow((-2-3), 2) + math.Pow((7.5-3), 2) + math.Pow((0.125-3), 2))

	for _, euct := range euclideanTests {
		got := Distance(euct.A, euct.B)
		if got != euct.want {
			t.Error(`got: `, got, `
            want: `, euct.want)
		}
	}

}

func Test_Func_Squared_Distance_Btwn_Datapoints(t *testing.T) {
	squaredTests := pointPairs
	squaredTests[0].want = 0*0 + 1*1
	squaredTests[1].want = math.Pow((-2-3), 2) + math.Pow((7.5-3), 2) + math.Pow((0.125-3), 2)

	for _, sqart := range squaredTests {
		got := DistanceSq(sqart.A, sqart.B)
		if got != sqart.want {
			t.Error(`got: `, got, `
            want: `, sqart.want)
		}
	}
}

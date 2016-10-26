package kdtree

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"
)

type intAndFloat struct {
	i int
	f float64
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

	embedded = struct {
		s string
		intAndFloat
	}{"hello world", intAndFloat{0, math.E}}

	rational = big.NewRat(5, 4)

	str = "cassandra"

	mappy = map[string]interface{}{
		"a": int('a'),
		"b": string("banana"),
		"c": intAndFloat{0, math.Pi},
	}

	interfaces = []interface{}{&str, &rational, &embedded, &mappy}

	constructorInputs = []struct {
		linked interface{}
		values []float64
	}{
		{
			&str,
			[]float64{6.0000125, 6.10000125, -1.3173, 1373},
		},
		{
			&rational,
			[]float64{1, 2, 3, 4, 5},
		},
		{
			&embedded,
			[]float64{100.3},
		},
		{
			&mappy,
			[]float64{0.8050908121798804, 0.53238545404102},
		},
	}
)

func TestDatapointConstructor(t *testing.T) {
	constructorTests := []struct {
		dp         *Datapoint
		wantData   interface{}
		wantValues []float64
	}{
		{
			dp:         NewDatapoint(constructorInputs[0].linked, constructorInputs[0].values),
			wantData:   constructorInputs[0].linked,
			wantValues: constructorInputs[0].values,
		},
		{
			dp:         NewDatapoint(constructorInputs[1].linked, constructorInputs[1].values),
			wantData:   constructorInputs[1].linked,
			wantValues: constructorInputs[1].values,
		},
		{
			dp:         NewDatapoint(constructorInputs[2].linked, constructorInputs[2].values),
			wantData:   constructorInputs[2].linked,
			wantValues: constructorInputs[2].values,
		},
		{
			dp:         NewDatapoint(constructorInputs[3].linked, constructorInputs[3].values),
			wantData:   constructorInputs[3].linked,
			wantValues: constructorInputs[3].values,
		},
	}

	for _, ct := range constructorTests {
		if reflect.DeepEqual(ct.wantData, ct.dp.Data()) == false {
			t.Fail()
		}
		if reflect.DeepEqual(ct.wantValues, ct.dp.Set()) == false {
			t.Fail()
		}
	}
}

func TestDatapointSetCopy(t *testing.T) {

	var randomPoints Datapoints
	for i := uint(1); i <= 10; i++ {
		p := RandomDatapoint(i)
		randomPoints = append(randomPoints, p)
	}

	type dpst struct {
		dp             *Datapoint
		wantValues     []float64
		duplicateAddrs []*float64
	}

	var dpSetTests []dpst

	for i := range randomPoints {
		var values []float64
		copy(randomPoints[i].set, values)
		var f64ptrs []*float64
		for k := range randomPoints[i].set {
			values = append(values, randomPoints[i].set[k])
			f64ptrs = append(f64ptrs, &randomPoints[i].set[k])
		}
		dpSetTests = append(
			dpSetTests,
			dpst{
				randomPoints[i],
				values,
				f64ptrs,
			})
	}

	for i := range dpSetTests {
		gotSet := dpSetTests[i].dp.Set()
		// verify that the values were copied correctly
		if reflect.DeepEqual(gotSet, dpSetTests[i].wantValues) == false {
			t.Errorf(`wrong values returned by Set()
            got:       %v
            want:      %v`, gotSet, dpSetTests[i].wantValues)
		}
		// but we also need to make sure it's a true copy - i.e. not referring to the same underlying array/slice
		gotAddrs := getF64ptrs(gotSet)
		if diffFptrs(gotAddrs, dpSetTests[i].duplicateAddrs) == false {
			t.Errorf(`duplicate addresses returned by Set()
		    got:       %v
		    matching:  %v`, gotAddrs, dpSetTests[i].duplicateAddrs)
		}
	}
}

func TestDatapointLinkedDataIdentical(t *testing.T) {
	for k, ct := range constructorInputs {
		newDP := NewDatapoint(ct.linked, ct.values)
		if reflect.DeepEqual(newDP.Data(), interfaces[k]) == false {
			t.Error(`Datapoint.Data() does not refer to the same object
			got:      `, newDP.Data(), `
			want:     `, interfaces[k])
		}

		fromDp := reflect.Indirect(reflect.ValueOf(newDP.Data()))
		fromSrc := reflect.Indirect(reflect.ValueOf(interfaces[k]))
		if reflect.DeepEqual(fromDp, fromSrc) == false {
			t.Error(`Datapoint.Data() does not refer to an identical value
			got:      `, fromDp, `
			want:     `, fromSrc)
		}
	}
}

func TestDatapointDimensionality(t *testing.T) {
	var dimTests = []struct {
		dp   *Datapoint
		want int
	}{
		{RandomDatapoint(3), 3},
		{RandomDatapoint(0), 0},
	}

	for _, dt := range dimTests {
		if dt.want != dt.dp.Dimensionality() {
			t.Fail()
		}
	}
}

func TestDatapointStringer(t *testing.T) {
	dpStringStr := `{data: cassandra}, {set: [0:{6.0000125}, 1:{6.10000125}, 2:{-1.3173}, 3:{1373}]}`
	dpRationalStr := `{data: &{{false [5]} {false [4]}}}, {set: [0:{1}, 1:{2}, 2:{3}, 3:{4}, 4:{5}]}`
	var stringerTests = []struct {
		dp   *Datapoint
		want string
	}{
		{
			&Datapoint{
				&str,
				[]float64{6.0000125, 6.10000125, -1.3173, 1373},
			},
			dpStringStr,
		},
		{
			&Datapoint{
				&rational,
				[]float64{1, 2, 3, 4, 5},
			},
			dpRationalStr,
		},
	}

	for _, s := range stringerTests {
		got := fmt.Sprint(s.dp)
		if got != s.want {
			t.Error(`want: `, s.want, `
			got: `, got)
		}
	}

}

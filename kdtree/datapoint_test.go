package kdtree

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func Test_Datapoint_Constructor(t *testing.T) {
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

func Test_Datapoint_Set_Copy(t *testing.T) {

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

func Test_Datapoint_Linked_Data_Identical(t *testing.T) {
	for k, ct := range constructorInputs {
		newDP := NewDatapoint(ct.linked, ct.values)
		if reflect.DeepEqual(newDP.Data(), interfaces[k]) == false {
			t.Error(k, `: Datapoint.Data() does not refer to the same object
			got:      `, newDP.Data(), `
			want:     `, interfaces[k])
		}

		srcRv := reflect.ValueOf(interfaces[k])
		dpRv := reflect.ValueOf(newDP.Data())

		if srcRv.Kind() != dpRv.Kind() { //	if they aren't to the same kind, then we have already failed.
			t.Error(k, `: Datapoint.Data() does not refer to an identical reflect.Kind
			got:      `, dpRv, `
			want:     `, srcRv)
		}

		errStringGot := fmt.Sprint(`got:      `, reflect.TypeOf(newDP.Data()), ` `, dpRv.Interface())
		errStringWant := fmt.Sprint(`want:     `, reflect.TypeOf(interfaces[k]), ` `, srcRv.Interface())

		if reflect.DeepEqual(reflect.TypeOf(newDP.Data()), reflect.TypeOf(interfaces[k])) == false {
			t.Error(k, `: Datapoint.Data() does not reflect an identical reflect.Type 
			`, errStringGot, `
			`, errStringWant)
		}
	}
}

func Test_Datapoint_Dimensionality(t *testing.T) {
	var dimTests = []struct {
		dp   *Datapoint
		want int
	}{
		{RandomDatapoint(3), 3},
		{RandomDatapoint(0), 0},
	}

	for _, dt := range dimTests {
		got := dt.dp.Dimensionality()
		if got != dt.want {
			t.Error(`want: `, dt.want, `
			got: `, got)
		}
	}
}

func Test_Datapoint_Stringer(t *testing.T) {
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

func Test_Datapoint_MarshalJSON(t *testing.T) {
	for k, ct := range constructorInputs {
		newDP := NewDatapoint(ct.linked, ct.values)
		dpJSON, _ := json.Marshal(newDP)
		var want []byte
		switch k {
		case 0:
			want = []byte(`{"data":"cassandra","set":[6.0000125,6.10000125,-1.3173,1373]}`)
		case 1:
			want = []byte(`{"data":"cassandra","set":[6.0000125,6.10000125,-1.3173,1373]}`)
		case 2:
			want = []byte(`{"data":"5/4","set":[1,2,3,4,5]}`)
		case 3:
			want = []byte(`{"data":"5/4","set":[1,2,3,4,5]}`)
		case 4:
			want = []byte(`{"data":{"intAndFloat":{"f":2.718281828459045,"i":0},"s":"hello world"},"set":[100.3]}`)
		case 5:
			want = []byte(`{"data":{"intAndFloat":{"f":2.718281828459045,"i":0},"s":"hello world"},"set":[100.3]}`)
		case 6:
			want = []byte(`{"data":{"a":97,"b":"banana","c":{"f":3.141592653589793,"i":0}},"set":[0.8050908121798804,0.53238545404102]}`)
		case 7:
			want = []byte(`{"data":{"a":97,"b":"banana","c":{"f":3.141592653589793,"i":0}},"set":[0.8050908121798804,0.53238545404102]}`)
		}
		if reflect.DeepEqual(dpJSON, want) == false {
			t.Error(`want: `, string(want), `
			got: `, string(dpJSON))
		}
	}
}

type char rune

func (c *char) ToDatapoint() *Datapoint {
	return &Datapoint{c, []float64{float64(*c) / 10.0, float64(*c) / 100.0}}
}

type myString string

func (str *myString) ToDatapoint() *Datapoint {
	var d Datapoint
	d.data = str
	var f []float64
	for _, s := range *str {
		f = append(f, math.Pi/float64(s))
	}
	d.set = f
	return &d
}

func (emb *embedTest) ToDatapoint() *Datapoint {
	var d Datapoint
	d.data = emb
	var f []float64
	inf := emb.intAndFloat
	for _, s := range emb.s {
		f = append(f, float64(s))
	}
	f = append(f, float64(inf.i))
	f = append(f, inf.f)
	d.set = f
	return &d
}

func Test_Datapoints_Import_with_Importable_Interface(t *testing.T) {
	A := char('a')
	B := char('b')
	S := myString("Hello!")
	E := embedTest{"Aloha!", intAndFloat{i: 1, f: math.Log2E}}
	importTests := []struct {
		imp  Importable
		want *Datapoint
	}{
		{
			imp:  &A,
			want: &Datapoint{&A, []float64{9.7, 0.97}},
		},
		{
			imp:  &B,
			want: &Datapoint{&B, []float64{9.8, 0.98}},
		},
		{
			imp:  &S,
			want: &Datapoint{&S, []float64{0.04363323129985824, 0.031104877758314782, 0.02908882086657216, 0.02908882086657216, 0.028302636518826967, 0.09519977738150888}},
		},
		{
			imp:  &E,
			want: &Datapoint{&E, []float64{65, 108, 111, 104, 97, 33, 1, 1.4426950408889634}},
		},
	}

	var dps Datapoints

	for k, it := range importTests {
		dps.Import(it.imp)
		if len(dps) != k+1 {
			t.Fail()
		}
		if reflect.DeepEqual(dps[k], it.want) == false {
			t.Error(`got: `, dps[k], `
		want: `, it.want)
		}
		if reflect.DeepEqual(reflect.TypeOf(dps[k].Data()), reflect.TypeOf(it.want.Data())) == false {
			t.Fail()
		}

		dpsRv := reflect.ValueOf(dps[k].Data())
		itwRv := reflect.ValueOf(it.want.Data())
		if reflect.DeepEqual(dpsRv, itwRv) == false {
			t.Errorf("\ngot:\t%v\nwant:\t%v\n", dpsRv, itwRv)
		}
		if reflect.DeepEqual(reflect.Indirect(dpsRv), reflect.Indirect(itwRv)) == false {
			t.Errorf("\ngot:\t%v\nwant:\t%v\n", reflect.Indirect(dpsRv), reflect.Indirect(itwRv))
		}
		if reflect.Indirect(dpsRv).CanSet() == false {
			t.Fail()
		}
	}
}

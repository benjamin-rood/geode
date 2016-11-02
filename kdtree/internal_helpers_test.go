package kdtree

import "testing"

func Test_Helpers_SumFloats(t *testing.T) {
	sumTests := []struct {
		set  []float64
		want float64
	}{
		{
			set:  []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			want: 45,
		},
		{
			set:  []float64{},
			want: 0,
		},
		{
			set:  []float64{0, -1, 5, 5000, 5000, 5000, 1},
			want: 15005,
		},
		{
			set: []float64{
				0.7289138247988962,
				0.13066099928359792,
				0.2977389733902786,
				0.38367281153431915,
				0.038781596430323854,
				0.0012844394790036428,
				0.5535391136630047,
				0.4396184083526607,
			},
			want: 2.574210166932085,
		},
		{
			set: []float64{
				-9.125,
				-2.5,
				11.75,
				0.5,
			},
			want: 0.625,
		},
	}

	for _, st := range sumTests {
		got := sum(st.set)
		if got != st.want {
			t.Errorf("sum%v == %v, want %v\n", st.set, got, st.want)
		}
	}
}

func Test_Helpers_datapoint_setString(t *testing.T) {
	d := &Datapoint{nil, []float64{0.6227283173637045, 0.3696928436398219}}
	want := `(0.6227283173637045, 0.3696928436398219)`
	got := d.setString()
	if got != want {
		t.Error(`want: `, string(want), `
		got: `, string(got))
	}
}

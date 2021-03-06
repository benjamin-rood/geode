package kdtree

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// Datapoint stores a set of floating-point values and a pointer to any other
// structure or type which you may wish to associate with the Datapoint.
type Datapoint struct {
	data interface{} // ideally a pointer to some other associated thing
	set  []float64
}

// Datapoints is a slice multiple of pointers to individual Datapoints
type Datapoints []*Datapoint

// NewDatapoint is an explicit constructor as an alternative to manual declaration
func NewDatapoint(data interface{}, points []float64) *Datapoint {
	if points == nil {
		points = []float64{}
	}
	f := make([]float64, len(points), len(points))
	copy(f, points)
	d := Datapoint{
		data: data,
		set:  f,
	}
	return &d
}

// Data returns the interface value of the object that the Datapoint is linked with.
func (d *Datapoint) Data() interface{} {
	return d.data
}

// Set returns a copy of the slice of floating-point values
func (d *Datapoint) Set() []float64 {
	var export = make([]float64, len(d.set), len(d.set))
	copy(export, d.set)
	return export
}

// Dimensionality returns spatial dimensions the Datapoint fits over.
func (d *Datapoint) Dimensionality() int {
	return len(d.set)
}

// RandomDatapoint will produce a 'free' PRNG Datapoint in n dimensions
// where all values in the set lie in [0,1).
// Useful for testing or adding noise to a dataset.
func RandomDatapoint(n uint) *Datapoint {
	return RandomDatapointInRange(n, 0, 1)
}

// RandomDatapointInRange will produce a 'free' PRNG Datapoint in n dimensions
// where all values in the set lie in [min,max).
// Useful for testing or adding noise to a dataset.
func RandomDatapointInRange(n uint, min, max float64) *Datapoint {
	f := make([]float64, n, n)
	for i := range f {
		f[i] = randomFloatInRange(min, max)
	}
	d := Datapoint{
		data: nil,
		set:  f,
	}
	return &d
}

func (d *Datapoint) setString() string {
	if d == nil {
		return ""
	}
	pointString := "("
	for _, f := range d.set {
		pointString += fmt.Sprint(f, ", ")
	}
	pointString = pointString[:len(pointString)-2]
	pointString += ")"
	return pointString
}

// String returns a formatted string presentation of the Datapoint object,
// implementing Stringer interface
// Bug – present, currently assumes data interface{} is a pointer and never a concrete type.
func (d *Datapoint) String() string {
	var present string
	present += fmt.Sprintf("{data: %v}, ", reflect.Indirect(reflect.ValueOf(d.data)))
	present += "{set: ["
	for i := range d.set {
		present += fmt.Sprintf("%d:{%v}, ", i, d.set[i])
	}
	present = present[0 : len(present)-2]
	present += "]}"
	return present
}

// By is the function signature required to wrap a given Less method as closure
type By func(p, q *Datapoint) bool

type datapointSorter struct {
	Datapoints
	by By // closure used in the Less method.
}

// Sort acts as interface implementation wrapper on a collection of Datapoints,
// called by functions with the By signature
func (by By) Sort(d Datapoints) {
	ds := &datapointSorter{
		Datapoints: d,
		by:         by,
	}
	sort.Sort(ds)
}

func (s *datapointSorter) Len() int {
	return len(s.Datapoints)
}

func (s *datapointSorter) Swap(i, j int) {
	s.Datapoints[i], s.Datapoints[j] = s.Datapoints[j], s.Datapoints[i]
}

func (s *datapointSorter) Less(i, j int) bool {
	return s.by(s.Datapoints[i], s.Datapoints[j])
}

// Comparator returns a dynamic "By" function on the specified plane,
// which gets passed to the Sort implementation's Less method.
func Comparator(plane int) By {
	return func(p, q *Datapoint) bool {
		return p.set[plane] < q.set[plane]
	}
}

// EqualTo provides a direct equality comparison between two Datapoints
func (d *Datapoint) EqualTo(q *Datapoint) bool {
	if len(d.set) != len(q.set) {
		return false
	}
	for i := range d.set {
		if d.set[i] != q.set[i] {
			return false
		}
	}
	return true
}

// EqualTo provides an equality comparison between each Datapoint in a set of Datapoints.
func (ds Datapoints) EqualTo(qs Datapoints) bool {
	if len(ds) != len(qs) {
		return false
	}
	for i := range ds {
		if !ds[i].EqualTo(qs[i]) {
			return false
		}
	}
	return true
}

// Import uses the Importable interface to cleanly append a single Datapoint to a the end of a set (slice) of Datapoints
func (ds *Datapoints) Import(I Importable) {
	*ds = append(*ds, I.ToDatapoint())
}

// PointsSetString returns a concatenated presentation of each Datapoint set as a single set presentation.
// e.g. `{(1,2) (3,4) (5,6) (7,8)}`
// possibly should be internal only?
func (ds *Datapoints) PointsSetString() string {
	pss := `{`
	for _, d := range *ds {
		if d == nil {
			continue
		}
		pss += (d.setString() + " ")
	}

	pss = pss[:len(pss)-1] + `}`
	return pss
}

func (ds Datapoints) notDistinct() bool {
	sz := len(ds)
	if sz <= 1 {
		return false
	}

	for i := 0; i < len(ds)-1; i++ {
		if !ds[i].EqualTo(ds[i+1]) {
			return false
		}
	}
	return true
}

// MarshalJSON implements encoding/json Marshaler interface
func (d *Datapoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"data": d.data,
		"set":  d.set,
	})
}

// TODO: Implement encoding/json Unmarshaler interface method
// func (d *Datapoint) UnmarshalJSON([]byte) error {
// }

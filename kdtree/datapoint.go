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

// Importable is the interface implemented by types who can be directly converted into a valid Datapoint.
type Importable interface {
	ToDatapoint() *Datapoint
}

// Exportable is the interface implemented by types which can be take a Datapoint and use the set of floating-point values to update the calling object's data members.
type Exportable interface {
	FromDatapoint(*Datapoint)
}

// Import uses the Importable interface to cleanly append a single Datapoint to a the end of a set (slice) of Datapoints
func (ds *Datapoints) Import(I Importable) {
	*ds = append(*ds, I.ToDatapoint())
}

// Convert uses the Importable interface to cleanly produce a kdtree
// from a slice of some type which has implented ToDataPoint(), and
// according to the pivot algorithm (PivotFunc).
func Convert(c []Importable, pivotDef PivotFunc) (*Branch, error) {
	var points = make(Datapoints, len(c), len(c))
	// basedim := len(c[0].ToDatapoint().set)
	for i := range c {
		points[i] = c[i].ToDatapoint()
		// TODO:
		// Implement build error types.
		// if points[i].Dimensionality() != basedim {
		// 	return nil, DimClashError
		// }
	}

	if pivotDef == nil {
		pivotDef = LazyAverage
	}

	b := Build(points, 0, pivotDef)
	return b, nil
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

package kdtree

// Importable is the interface implemented by types who can be directly converted into a valid Datapoint.
type Importable interface {
	ToDatapoint() *Datapoint
}

// Exportable is the interface implemented by types which can be take a Datapoint and use the set of floating-point values to update the calling object's data members.
type Exportable interface {
	FromDatapoint(*Datapoint)
}

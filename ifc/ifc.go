package ifc

// Gridder specifies and interface for getting and setting values in a grid
// Coordinates referenced outside of the grids bounds are expected to result in a noop
// Swap is expected to efficiently swap all grid values between receiver and argument
type Gridder interface {
	Set(i, j int, val float64)
	Get(i, j int) (val float64)
	Swap(other Gridder)
}

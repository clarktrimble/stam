// Package vektor isolates the mysteries of Swap for further investigation
package vektor

import (
	"fmt"

	"github.com/clarktrimble/stam/ifc"
)

// Vektor implements the Gridder interface
// in the c stlyle as seen in stam
// https://www.josstam.com/_files/ugd/cf1fd6_9989229efbd34a26ba5ccd913721a2ac.pdf
//
// Basically a 1D vector stores the 2D values
// and two grid coordinates are translated to a single index
//
// Perhaps the golang offers more performant or idomatic approaches
// - multi dimension slice, prolly better?
// - array (sans slice), doubtful?
// In the meantime, we have some isolation via the interface
type Vektor struct {
	dim  int
	vals *[]float64
}

// New creates a vektor, given:
//
//	dim: width and height of grid
func New(dim int) *Vektor {

	vals := make([]float64, dim*dim)
	return &Vektor{
		dim:  dim,
		vals: &vals,
	}
}

// Get gets the value for a given cell
func (vk *Vektor) Get(i, j int) float64 {

	n := vk.index(i, j)
	if n < 0 {
		return 0
	}

	return (*vk.vals)[n]
}

// Set sets a value for the given cell
func (vk *Vektor) Set(i, j int, val float64) {

	n := vk.index(i, j)
	if n < 0 {
		return
	}

	(*vk.vals)[n] = val
}

// Swap swaps the value slices of two vektors
// for this to work:
// - pass in pointers to the vektors (of course)
// - respective vektors hold pointers to value slices
// "type vektorAlt []float64" with swapping in layer could be cleaner?
// Todo: yes!
func (vk *Vektor) Swap(other ifc.Gridder) {

	ov, ok := other.(*Vektor)
	if !ok {
		panic(fmt.Sprintf("somehow asked to swap a non-vector: %#v", other))
	}

	if vk.dim != ov.dim {
		panic("will not swap vectors of differing dimentions")
	}

	tmp := vk.vals
	vk.vals = ov.vals
	ov.vals = tmp
}

// unexported

func (vk *Vektor) index(i, j int) (n int) {

	n = i + vk.dim*j

	if n >= vk.dim*vk.dim {
		n = -1
	}
	return
}

// notes

// why are these so slow ???
/*
func (vk *Vektor) Zero() {

	for n := 0; n < vk.Dim*vk.Dim; n++ {
		(*vk.vals)[n] = 0
	}
}

func (vk *Vektor) MinMax() (min, max float64) {

	for _, density := range *vk.vals {
		if max < density {
			max = density
		}
		if min > density {
			min = density
		}
	}

	return
}
*/

// what's actually performant?

//  multidim!! ?
// https://stackoverflow.com/questions/39804861/what-is-a-concise-way-to-create-a-2d-slice-in-go
/*
matrix := make([][]int, n)
rows := make([]int, n*m)
for i := 0; i < n; i++ {
    matrix[i] = rows[i*m : (i+1)*m]
}
*/
//  generic ??
/*
func Make2D[T any](n, m int) [][]T {
    matrix := make([][]T, n)
    rows := make([]T, n*m)
    for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
        endRow := startRow + m
        matrix[i] = rows[startRow:endRow:endRow]
    }
    return matrix
}
a := Make2D[uint8](dy, dx)
*/

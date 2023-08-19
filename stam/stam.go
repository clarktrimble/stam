// Package stam simulates dye in a 2-D incompressible fluid
//
// as presented by Joe Stam in Real-Time Fluid Dynamics for Games in 2003
// https://www.josstam.com/_files/ugd/cf1fd6_9989229efbd34a26ba5ccd913721a2ac.pdf
//
// see "begin stam" comment below for code transliterated from the paper
package stam

import (
	"github.com/clarktrimble/stam/ifc"
)

// Fluid is a fluid model
type Fluid struct {
	size int
	dt   float64
	Visc float64
	Diff float64
	d0   ifc.Gridder
	dd   ifc.Gridder
	u0   ifc.Gridder
	uu   ifc.Gridder
	v0   ifc.Gridder
	vv   ifc.Gridder
}

// NewFluid creates a fluid model, given:
//
//		gridSize: width and height of grid in terms of cells
//	           note: a boundary layer wraps around the grid so that it ranges from 0,0 to size+1,size+1
//		visc:     viscocity of fluid
//		diff:     diffusivity of dye in fluid
//		dt:       change in time per step
//		factory:  funcion returning an implementation of the Gridder interface
func NewFluid(gridSize int, visc, diff, dt float64, factory func(size int) ifc.Gridder) (fluid Fluid) {

	count := (gridSize + 2) * (gridSize + 2)

	fluid = Fluid{
		size: gridSize,
		Visc: visc,
		Diff: diff,
		dt:   dt,
		d0:   factory(count),
		u0:   factory(count),
		v0:   factory(count),
		dd:   factory(count),
		uu:   factory(count),
		vv:   factory(count),
	}

	return
}

// Density gets dye density for a given cell
func (fluid Fluid) Density(i, j int) float64 {
	return fluid.dd.Get(i, j)
}

// Velocity gets fluid velocity for a given cell
func (fluid Fluid) Velocity(i, j int) (u, v float64) {
	return fluid.uu.Get(i, j), fluid.vv.Get(i, j)
}

// Min finds the minimum dye density
func (fl Fluid) Min() (min float64) {

	min = 99999 // hax!?

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			dd := fl.dd.Get(i, j)
			if min > dd {
				min = dd
			}
		}
	}

	return
}

// Max finds the maximum dye density
func (fl Fluid) Max() (max float64) {

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			dd := fl.dd.Get(i, j)
			if max < dd {
				max = dd
			}
		}
	}

	return
}

// AddVelocity adds velocity to a square of cells
func (fluid Fluid) AddVelocity(i, j, n int, u, v float64) {

	// add velocity to u0, v0 where it will be picked up by addSource

	fluid.zero(fluid.u0)
	fluid.zero(fluid.v0)

	for l := -n; l <= n; l++ {
		for m := -n; m <= n; m++ {
			fluid.u0.Set(i+l, j+m, u)
			fluid.v0.Set(i+l, j+m, v)
		}
	}
}

// AddDensity adds dye to a square of cells
func (fluid Fluid) AddDensity(i, j, n int, d float64) {

	// add density to d0 where it will be picked up by addSource

	fluid.zero(fluid.d0)

	for l := -n; l <= n; l++ {
		for m := -n; m <= n; m++ {
			fluid.d0.Set(i+l, j+m, d)
		}
	}
}

// Level subtracts a small amount of dye from all cells
func (fl *Fluid) Level(min float64) {

	threshold := 0.0001

	if min < threshold {
		return
	}

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {

			result := fl.dd.Get(i, j) - 4*min

			if result < threshold {
				result = 0
			}
			fl.dd.Set(i, j, result)
		}
	}

	fl.setBnd(0, fl.dd)
}

func (fl *Fluid) zero(xx ifc.Gridder) {

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {

			xx.Set(i, j, 0)
		}
	}

	fl.setBnd(0, xx)
}

// end non-stam
// begin stam

// Step calculates the next state of the fluid
func (fl *Fluid) Step() {

	// uu and vv contain current velocities
	// dd contains current densities

	//  u0 and v0 are set with "source", add to uu and vv
	fl.addSource(fl.uu, fl.u0)
	fl.addSource(fl.vv, fl.v0)

	// put uu into u0, vv into v0
	fl.uu.Swap(fl.u0)
	fl.vv.Swap(fl.v0)

	// diffuse for a new uu and vv
	fl.diffuse(1, fl.Visc, fl.uu, fl.u0)
	fl.diffuse(2, fl.Visc, fl.vv, fl.v0)

	// project, you know, for mass conservation and swap
	fl.project(fl.uu, fl.vv, fl.u0, fl.v0)
	fl.uu.Swap(fl.u0)
	fl.vv.Swap(fl.v0)

	// advect and project for mass conservation again
	fl.advect(1, fl.uu, fl.u0, fl.u0, fl.v0)
	fl.advect(2, fl.vv, fl.v0, fl.u0, fl.v0)
	fl.project(fl.uu, fl.vv, fl.u0, fl.v0)

	//

	// d0 is set to "source", add to dd and swap
	fl.addSource(fl.dd, fl.d0)
	fl.dd.Swap(fl.d0)

	// diffuse, swap, advect
	fl.diffuse(0, fl.Diff, fl.dd, fl.d0)
	fl.dd.Swap(fl.d0)
	fl.advect(0, fl.dd, fl.d0, fl.uu, fl.vv)

	// now uu and vv contain newly computed velocities
	// and dd contains newly computed densities
	return
}

func (fl *Fluid) addSource(dst, src ifc.Gridder) {

	// Todo: consider combining with AddVelocity/Density
	// seems better, but would be less stam'ish

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			val := dst.Get(i, j) + fl.dt*src.Get(i, j)
			dst.Set(i, j, val)
		}
	}
}

func (fl *Fluid) setBnd(bnd int, xx ifc.Gridder) {

	// usually set boundary to same as adjacent cell, but:
	// if bnd -> 1 then negate column boundaries
	// if bnd -> 2 then negate row boundaries
	// ( negation is used to cancel velocity )

	for i := 1; i <= fl.size; i++ {

		// set 0 and N+1 columns
		if bnd == 1 {
			xx.Set(0, i, -xx.Get(1, i))
			xx.Set(fl.size+1, i, -xx.Get(fl.size, i))
		} else {
			xx.Set(0, i, xx.Get(1, i))
			xx.Set(fl.size+1, i, xx.Get(fl.size, i))
		}

		// set 0 and N+1 rows
		if bnd == 2 {
			xx.Set(i, 0, -xx.Get(i, 1))
			xx.Set(i, fl.size+1, -xx.Get(i, fl.size))
		} else {
			xx.Set(i, 0, xx.Get(i, 1))
			xx.Set(i, fl.size+1, xx.Get(i, fl.size))
		}
	}

	// average corners

	result := 0.5 * (xx.Get(1, 0) + xx.Get(0, 1))
	xx.Set(0, 0, result)
	result = 0.5 * (xx.Get(1, fl.size+1) + xx.Get(0, fl.size))
	xx.Set(0, fl.size+1, result)
	result = 0.5 * (xx.Get(fl.size, 0) + xx.Get(fl.size+1, 1))
	xx.Set(fl.size+1, 0, result)
	result = 0.5 * (xx.Get(fl.size, fl.size+1) + xx.Get(fl.size+1, fl.size))
	xx.Set(fl.size+1, fl.size+1, result)
}

func (fl *Fluid) diffuse(bnd int, diff float64, xx, x0 ifc.Gridder) {

	a := fl.dt * diff * float64(fl.size*fl.size)

	// Todo: well named constant for 20
	for k := 0; k < 20; k++ {
		for i := 1; i <= fl.size; i++ {
			for j := 1; j <= fl.size; j++ {

				num := x0.Get(i, j) + a*(xx.Get(i-1, j)+xx.Get(i+1, j)+xx.Get(i, j-1)+xx.Get(i, j+1))
				result := num / (1 + 4*a)

				xx.Set(i, j, result)
			}
		}
		fl.setBnd(bnd, xx)
	}

	return
}

func (fl *Fluid) advect(bnd int, dd, d0, uu, vv ifc.Gridder) {
	var (
		i0, j0, i1, j1       int
		x, y, s0, t0, s1, t1 float64
	)

	nf := float64(fl.size)
	dt0 := fl.dt * nf

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			x = float64(i) - dt0*uu.Get(i, j)
			y = float64(j) - dt0*vv.Get(i, j)

			if x < 0.5 {
				x = 0.5
			}
			if x > nf+0.5 {
				x = nf + 0.5
			}
			i0 = int(x)
			// Todo: ^^^ same as in C??
			i1 = i0 + 1

			if y < 0.5 {
				y = 0.5
			}
			if y > nf+0.5 {
				y = nf + 0.5
			}
			j0 = int(y)
			// Todo: ^^^ same as in C??
			j1 = j0 + 1

			s1 = x - float64(i0)
			s0 = 1 - s1
			t1 = y - float64(j0)
			t0 = 1 - t1

			result := s0*(t0*d0.Get(i0, j0)+t1*d0.Get(i0, j1)) + s1*(t0*d0.Get(i1, j0)+t1*d0.Get(i1, j1))
			dd.Set(i, j, result)
		}
	}
	fl.setBnd(bnd, dd)
}

func (fl *Fluid) project(u, v, p, div ifc.Gridder) {
	// Todo: p and div are mysteriously named
	h := 1.0 / float64(fl.size)

	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			result := -.05 * h * (u.Get(i+1, j) - u.Get(i-1, j) + v.Get(i, j+1) - v.Get(i, j-1))
			div.Set(i, j, result)
			p.Set(i, j, 0)
		}
	}
	fl.setBnd(0, div)
	fl.setBnd(0, p)
	for k := 0; k < 20; k++ {
		for i := 1; i <= fl.size; i++ {
			for j := 1; j <= fl.size; j++ {
				result := (div.Get(i, j) + p.Get(i-1, j) + p.Get(i+1, j) + p.Get(i, j-1) + p.Get(i, j+1)) / 4
				p.Set(i, j, result)
			}
		}
		fl.setBnd(0, p)
	}
	for i := 1; i <= fl.size; i++ {
		for j := 1; j <= fl.size; j++ {
			result := u.Get(i, j) - 0.5*(p.Get(i+1, j)-p.Get(i-1, j))/h
			u.Set(i, j, result)
			result = v.Get(i, j) - 0.5*(p.Get(i, j+1)-p.Get(i, j-1))/h
			v.Set(i, j, result)
		}
	}
	fl.setBnd(1, u)
	fl.setBnd(2, v)
}

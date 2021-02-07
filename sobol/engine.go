package sobol

import (
	"fmt"
	"math"
)

// findRightmostZeroBit returns index from the right of the first zero bit of n.
func findRightmostZeroBit(n uint32) uint32 {
	c := uint32(0)
	for n&(1<<c) != 0 {
		c++
	}
	return c + 1 // starts from 1
}

func initDirectionNumbers(dim uint32) [][]uint32 {
	v := make([][]uint32, dim)
	for i := uint32(0); i < dim; i++ {
		v[i] = make([]uint32, maxBit)
	}

	// First row of sobol state is all '1'.
	for m := 0; m < maxBit; m++ {
		v[0][m] = 1 << (32 - m) // all m's = 1
	}

	// Remaining rows of sobol state (row 2 through dim, indexed by [1:dim])
	for j := uint32(1); j < dim; j++ {
		v[j] = make([]uint32, maxBit+1)

		// Read in parameters from file
		// Skip 1000 lines from the top as Joe&Kuo's C++ program do.
		dn := directionNumbers[1000+j]
		for i := uint32(1); i <= dn.S; i++ {
			v[j][i] = dn.M[i-1] << (32 - i)
		}
		for i := dn.S + 1; i <= maxBit; i++ {
			v[j][i] = v[j][i-dn.S] ^ (v[j][i-dn.S] >> dn.S)
			for k := uint32(1); k <= dn.S-1; k++ {
				v[j][i] ^= ((dn.A >> (dn.S - 1 - k)) & 1) * v[j][i-k]
			}
		}

	}
	return v
}

// Engine is Sobol's quasirandom number generator.
type Engine struct {
	dim uint32     // dimensions
	n   uint32     // the number of generate times
	v   [][]uint32 // direction numbers
	x   [][]uint32
}

// NewEngine returns Sobol's quasirandom number generator.
func NewEngine(dimension uint32) *Engine {
	if dimension > maxDim {
		panic(fmt.Errorf("maximum supported dimensionality is %d", maxDim))
	}

	v := initDirectionNumbers(dimension)
	x := make([][]uint32, dimension+1)
	for i := uint32(0); i <= dimension; i++ {
		// Pre-allocate memory to sample 512 points
		x[i] = make([]uint32, 0, 512)
		x[i] = append(x[i], 0)
	}

	return &Engine{
		dim: dimension,
		n:   0,
		v:   v,
		x:   x,
	}
}

// Draw samples from Sobol sequence.
func (e *Engine) Draw() []float64 {
	e.n++
	points := make([]float64, e.dim)

	for j := uint32(0); j < e.dim; j++ {
		c := findRightmostZeroBit(e.n - 1)
		e.x[j] = append(e.x[j], e.x[j][e.n-1]^e.v[j][c])
		points[j] = float64(e.x[j][e.n]) / math.Pow(2.0, 32)
	}
	return points
}
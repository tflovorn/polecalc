package polecalc

import (
	"math"
)

// Return a square mesh coordinate corresponding to the index i.
// i=0 corresponds to (-pi, -pi); i=L-1 is (pi-step, -pi);
// i=L is (-pi, -pi+step); i=L^2-1 is (pi-step, pi-step)
func SquareAt(i uint64, L uint32) Vector2 {
	if i < 0 || i >= uint64(math.Pow(float64(L), 2.0)) {
		// panic here instead of returning an error since we will call
		// this function pretty often - presumably only returning one
		// variable is better for performance
		panic("invalid index for square mesh")
	}
	start, stop := -math.Pi, math.Pi
	length := stop - start
	step := length / float64(L)
	// transform 1d index to 2d coordinate indices
	ny := uint32(math.Floor(float64(i) / float64(L)))
	nx := uint32(i) - ny * L
	// get coordinates
	x := start + float64(nx) * step
	y := start + float64(ny) * step
	return Vector2{x, y}
}

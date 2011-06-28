package polecalc

import "math"

// Return two channels:
// The first contains points on a square grid with boundaries (-Pi, Pi) x (-Pi, Pi).
// (this is the first Brillouin zone of a square lattice)
// When all points have been consumed, the value true is passed on the second channel.
func Square(pointsPerSide uint32) chan Vector2 {
	cmesh := make(chan Vector2)
	go helpSquare(cmesh, pointsPerSide)
	return cmesh
}

// Do the work of generating the square mesh.
func helpSquare(cmesh chan Vector2, pointsPerSide uint32) {
	length := 2 * math.Pi
	step := length / float64(pointsPerSide)
	begin := -math.Pi
	end := math.Pi - step
	x, y := begin, begin
	for y <= end {
		for x <= end {
			cmesh <- Vector2{x, y}
			x += step
		}
		y += step
		x = begin
	}
	close(cmesh)
}

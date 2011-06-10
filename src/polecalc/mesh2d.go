package polecalc

import "math"

// Return two channels:
// The first contains points on a square grid with boundaries (-Pi, Pi) x (-Pi, Pi).
// (this is the first Brillouin zone of a square lattice)
// When all points have been consumed, the value true is passed on the second channel.
func Square(pointsPerSide uint32) (chan []float64, chan bool) {
	cmesh := make(chan []float64)
	done := make(chan bool)
	go helpSquare(cmesh, done, pointsPerSide)
	return cmesh, done
}

// Do the work of generating the square mesh.
func helpSquare(cmesh chan []float64, done chan bool, pointsPerSide uint32) {
	length := 2 * math.Pi
	step := length / float64(pointsPerSide)
	x, y := -math.Pi, -math.Pi
	for y < math.Pi {
		for x < math.Pi {
			cmesh <- []float64{x, y}
			x += step
		}
		y += step
		x = -math.Pi
	}
	for {
		done <- true
	}
}

package polecalc

import (
	"math"
	"testing"
)

// Does Square stay within (-Pi, Pi) x (-Pi, Pi)?
// --- todo: make pointsPerSide arbitrary (although not too big) ---
// Slow with large pointsPerSide - interesting performance test.
func TestSquareBounds(t *testing.T) {
	var pointsPerSide uint32 = 128
	cmesh := Square(pointsPerSide)
	done := make(chan bool)
	go func() {
		for point, ok := <-cmesh; ok; point, ok = <-cmesh {
			x, y := point.X, point.Y
			if x > math.Pi || x < -math.Pi {
				t.Fatalf("x out of bounds (x=%f, y=%f)", x, y)
			}
			if y > math.Pi || y < -math.Pi {
				t.Fatalf("y out of bounds (x=%f, y=%f)", x, y)
			}
		}
		done <- true
	}()
	<-done // wait on one worker
}

// Does Square produce the expected number of points?
func TestSquarePointNumber(t *testing.T) {
	var pointsPerSide uint32 = 128
	var count uint64 = 0
	done := make(chan bool)
	cmesh := Square(pointsPerSide)
	go func() {
		for _, ok := <-cmesh; ok; _, ok = <-cmesh {
			count++
		}
		done <- true
	}()
	<-done // wait on one worker
	points64 := uint64(pointsPerSide)
	expectedCount := points64 * points64
	if count != expectedCount {
		t.Fatalf("point total (%d) != expected points (%d)", count, expectedCount)
	}
}

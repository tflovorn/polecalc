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
	cmesh, done := Square(pointsPerSide)
	go func() {
		for {
			point := <-cmesh
			x, y := point[0], point[1]
			if x > math.Pi || x < -math.Pi {
				t.Fatalf("x out of bounds (x=%f, y=%f)", x, y)
			}
			if  y > math.Pi || y < -math.Pi {
				t.Fatalf("y out of bounds (x=%f, y=%f)", x, y)
			}
		}
	}()
	<-done
}

// Does Square produce the expected number of points?
func TestSquarePointNumber(t *testing.T) {
	var pointsPerSide uint32 = 128
	var count uint64 = 0
	cmesh, done := Square(pointsPerSide)
	go func() {
		for {
			<-cmesh
			count++
		}
	}()
	<-done
	points64 := uint64(pointsPerSide)
	expectedCount := points64 * points64
	if count != expectedCount {
		t.Fatalf("point total (%d) != expected points (%d)", count, expectedCount)
	}
}

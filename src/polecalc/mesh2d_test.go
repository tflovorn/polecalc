package polecalc

import (
	"math"
	"testing"
)

// Does Square stay within (-Pi, Pi) x (-Pi, Pi)?
func TestBounds(t *testing.T) {
	var pointsPerSide uint32 = 128
	cmesh, done := Square(pointsPerSide)
	go func() {
		for {
			point := <-cmesh
			x, y := point[0], point[1]
			if x > math.Pi || x < -math.Pi || y > math.Pi || y < -math.Pi {
				t.FailNow()
			}
		}
	}()
	<-done
}

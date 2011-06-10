package polecalc

import (
	"testing"
	"math"
)

// Average of Sin(kx) * Sin(ky) over Brillouin zone should be 0.
// Check if the average found is under machine epsilon.
func TestSinSin(t *testing.T) {
	worker := func(cmesh chan []float64, accum chan float64) {
		for {
			k := <-cmesh
			kx, ky := k[0], k[1]
			accum <- math.Sin(kx) * math.Sin(ky)
		}
	}
	if avg := Average(128, worker, 4); avg > MachEpsFloat64() {
		t.FailNow()
	}
}

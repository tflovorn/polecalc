package polecalc

import (
	"testing"
	"math"
)

// Average of Sin(kx) * Sin(ky) over Brillouin zone should be 0.
// Check if the average found is under machine epsilon.
func TestSinSin(t *testing.T) {
	worker := func(k []float64) float64 {
		return math.Sin(k[0]) * math.Sin(k[1])
	}
	if avg := Average(128, worker, 4); avg > MachEpsFloat64() {
		t.FailNow()
	}
}

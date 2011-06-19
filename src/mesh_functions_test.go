package polecalc

import (
	"testing"
	"math"
)

// Average of Sin(kx) * Sin(ky) over Brillouin zone should be 0.
// Check if the average found is under machine epsilon.
// Also check if the minimum is close to -1 (arbitrary tolerance 1e-9)
func TestSinSin(t *testing.T) {
	worker := func(k []float64) float64 {
		return math.Sin(k[0]) * math.Sin(k[1])
	}
	if avg := Average(128, worker, 4); math.IsNaN(avg) || avg > MachEpsFloat64() {
		t.Fatalf("average of sin(kx)*sin(ky) incorrect (got %f)", avg)
	}
	if min := Minimum(128, worker, 4); math.IsNaN(min) || (min+1) > 1e-9 {
		t.Fatalf("minimum of sin(kx)*sin(ky) incorrect (got %f)", min)
	}
}

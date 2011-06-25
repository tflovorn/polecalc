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

// Check if DeltaBinner is working for 2 non-q-dependent delta functions
func TestTwoDeltas(t *testing.T) {
	var pointsPerSide uint32 = 64
	deltaPoints := []float64{5.5, 10.5}
	deltaTerms := func(q []float64) ([]float64, []float64) {
		omegas := deltaPoints
		coeffs := []float64{1.0, 1.0}
		return omegas, coeffs
	}
	binner := NewDeltaBinner(deltaTerms, 0.0, 15.0, 64)
	result := DeltaBin(pointsPerSide, binner, 1)
	expected := 1.0
	for _, point := range deltaPoints {
		index := binner.BinVarToIndex(point)
		if math.Fabs(result[index]-expected) > MachEpsFloat64() {
			t.Fatalf("incorrect delta sum (%f)", result[index])
		}
	}
}

package polecalc

import (
	"testing"
	"math"
)

// Does the cubic spline integral produce the expected result when 
// interpolating over only 3 points?
func TestCubicSplineIntegral3Points(t *testing.T) {
	someCubic := makeCubic(1.0, 1.0, 1.0, 1.0)
	// constants chosen to have xs be integer values
	n := 3 // number of points to interpolate	
	start, stop := -10.0, 10.0
	step := (stop - start) / float64(n-1)
	xs, ys := make([]float64, n), make([]float64, n)
	for i, _ := range xs {
		xs[i] = start + float64(i)*step
		ys[i] = someCubic(xs[i])
	}
	xl, xr := -9.0, 1.0
	val, err := SplineIntegral(xs, ys, xl, xr)
	if err != nil {
		t.Fatal(err)
	}
	// expect spline coefficients as follows:
	// a = [0.05 -0.5], b = [0 1.5], c = [86 101], d = [-909 1]
	valKnown := -3747.025
	// 1e-15 is somewhat arbitrary
	// (couldn't get results to match to better accuracy)
	if math.Fabs((val-valKnown)/valKnown) > 1e-15 {
		t.Fatalf("failed to reproduce integral (got %f, expected %f)", val, valKnown)
	}
}

// Does the principal value integral return the correct result for a removable
// singularity?  (simulate this singularity by integrating a constant)
func TestPrincipalValueRemovable(t *testing.T) {
	constant := func(x float64) float64 {
		return 1.0
	}
	eps := 1e-9
	a, b := 0.0, 5.0
	w := (a + b) / 2.0
	integral, err := PvIntegral(constant, a, b, w, eps, uint(256))
	if err != nil {
		t.Fatal(err)
	}
	expected := b - a
	tolerance := 1e-6
	if math.Fabs(integral-expected) > tolerance {
		t.Fatalf("pv integral gave incorrect value (got %f, expected %f)", integral, expected)
	}
}

package polecalc

import (
	"testing"
	"math"
	"fmt"
)

// Does the cubic spline produce the expected result when interpolating over
// only 3 points?
func TestCubicSpline3Points(t *testing.T) {
	// constants chosen to have xs be integer values
	n := 3 // number of points to interpolate	
	start, stop := -10.0, 10.0
	someCubic := makeCubic(1.0, 1.0, 1.0, 1.0)
	step := (stop - start) / float64(n-1)
	xs, ys := make([]float64, n), make([]float64, n)
	for i, _ := range xs {
		xs[i] = start + float64(i)*step
		ys[i] = someCubic(xs[i])
	}
	spline, err := NewCubicSpline(xs, ys)
	if err != nil {
		t.Fatal(err)
	}
	x := 1.0
	y := spline.At(x)
	// expect spline coefficients as follows:
	// a = [0.05 -0.5], b = [0 1.5], c = [86 101], d = [-909 1]
	yKnown := 103.45
	if math.Fabs(y-yKnown) > MachEpsFloat64() {
		fmt.Printf("a:%v b:%v c:%v d:%v\n", spline.a, spline.b, spline.c, spline.d)
		fmt.Printf("xs = %v; ys = %v\n", xs, ys)
		t.Fatalf("failed to reproduce interpolation at known value (at %f got %f, expected %f)", x, y, yKnown)
	}
}

// Does the cubic spline error become small when using many points?
func TestCubicSplineManyPoints(t *testing.T) {
	accuracy := 1e-6
	// constants chosen to have xs be integer values
	n := 10001 // number of points to interpolate
	start, stop := -10.0, 10.0
	someCubic := makeCubic(1.0, 1.0, 1.0, 1.0)
	step := (stop - start) / float64(n-1)
	xs, ys := make([]float64, n), make([]float64, n)
	for i, _ := range xs {
		xs[i] = start + float64(i)*step
		ys[i] = someCubic(xs[i])
	}
	spline, err := NewCubicSpline(xs, ys)
	if err != nil {
		t.Fatal(err)
	}
	for i, xi := range xs {
		if i == n-1 {
			// don't jump out of the interpolation range
			continue
		}
		x := xi + step/2
		y := spline.At(x)
		yKnown := someCubic(x)
		if math.Fabs((y-yKnown)/yKnown) > accuracy {
			t.Fatalf("failed to interpolate to expected accuracy (at %f got %f, expected %f)", x, y, yKnown)
		}
	}
}

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
	val, err := SplineIntegral(xs, ys, xl, xr, uint(n))
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


// Build a cubic functions with given coefficients
func makeCubic(a, b, c, d float64) func(x float64) float64 {
	cubic := func(x float64) float64 {
		return a*math.Pow(x, 3.0) + b*math.Pow(x, 2.0) + c*x + d
	}
	return cubic
}

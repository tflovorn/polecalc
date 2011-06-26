package polecalc

import (
	"testing"
	"math"
)

// Simple 3x3 case for tridiagonal matrix equation
// [4 1 0  [x1    [1
//  1 4 1   x2  =  1
//  0 1 4]  x3]    1]
// the solution is given by: x = [3/14 1/7 3/14]
func TestTridiagonal3(t *testing.T) {
	a := []float64{0, 1, 1}
	b := []float64{4, 4, 4}
	c := []float64{1, 1, 0}
	d := []float64{1, 1, 1}
	x := TridiagonalSolve(a, b, c, d)
	expected := []float64{3.0 / 14.0, 1.0 / 7.0, 3.0 / 14.0}
	for i, val := range expected {
		if math.Fabs(val-x[i]) > MachEpsFloat64() {
			t.Fatalf("tridiagonal solution failed: got (%v), expected (%v)", x, expected)
		}
	}
}

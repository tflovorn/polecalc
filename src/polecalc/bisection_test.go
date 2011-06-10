package polecalc

import (
	"math"
	"testing"
	"testing/quick"
)

// Does SolveBisection correctly solve a simple linear function?
func TestLinear(t *testing.T) {
	makeLinear := func(root float64) func(float64) float64 {
		return func(x float64) float64 { return root - x }
	}
	checker := func(root float64) bool {
		scale := math.Fabs(root)
		epsilon := scale / 1e-5
		val, err := SolveBisection(makeLinear(root), root - scale, root + scale, epsilon)
		if err != nil {
			return false
		}
		return val == root
	}
	if !(checker(5) && checker(100) && checker(1e-5)) {
		t.FailNow()
	}
}

// Does BisectionFullPrecision correctly solve arbitrary simple linear functions?
func TestLinearFullPrecision(t *testing.T) {
	macheps := math.Pow(2.0, -53.0)
	makeLinear := func(root float64) func(float64) float64 {
		return func(x float64) float64 { return root - x }
	}
	checker := func(root float64) bool {
		scale := math.Fabs(root)
		val, err := BisectionFullPrecision(makeLinear(root), root - scale, root + scale)
		if err != nil {
			return false
		}
		return val - root <= macheps
	}
	if err := quick.Check(checker, nil); err != nil {
		t.Error(err)
	}
}

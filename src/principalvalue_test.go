package polecalc

import (
	"testing"
	"math"
)

func TestPrincipalValueGSLConstant(t *testing.T) {
	constant := func(x float64) float64 {
		return 1.0
	}
	epsabs := 1e-9
	epsrel := 1e-9
	limit := uint16(2048)
	a, b := 0.0, 5.0
	c := (a + b) / 2.0
	integral := PvIntegralGSL(constant, a, b, c, epsabs, epsrel, limit)
	expected := math.Log(math.Fabs((b - c) / (a - c)))
	if math.Fabs(integral-expected) > epsabs {
		t.Fatalf("tolerance exceeded")
	}
}

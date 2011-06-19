package polecalc

import (
	"testing"
	"math"
)

// Is EpsilonMin set such that the minimum of Epsilon is 0?
func TestEpsilonMin(t *testing.T) {
	env, err := EnvironmentFromFile("environment_test.json")
	env.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	worker := func(k []float64) float64 {
		return Epsilon(*env, k)
	}
	min := Minimum(env.GridLength, worker, env.NumProcs)
	if math.Fabs(min) > MachEpsFloat64() {
		t.Fatalf("minimum of Epsilon too large (%f)", min)
	}
}

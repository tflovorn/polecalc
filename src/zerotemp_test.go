package polecalc

import (
	"testing"
	"reflect"
	"flag"
	"math"
	"fmt"
)

var cached *bool = flag.Bool("gc_cache", false, "used cached Environment for TestGc0")

func TestKnownZeroTempSystem(t *testing.T) {
	envStr := "{\"GridLength\":8,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.05777149373506872,\"Mu\":-0.18330570279347036,\"F0\":0.12945949461029926,\"EpsilonMin\":-1.8}"
	expectedEnv, err := EnvironmentFromString(envStr)
	if err != nil {
		t.Fatal(err)
	}
	tolerances := []float64{1e-6, 1e-6, 1e-6}
	system := NewZeroTempSystem(tolerances)
	env, err := EnvironmentFromFile("zerotemp_test.json")
	if err != nil {
		t.Fatal(err)
	}
	env.Initialize()
	solution, err := system.Solve(*env)
	if err != nil {
		t.Fatal(err)
	}
	solvedEnv := solution.(Environment)
	if !reflect.DeepEqual(solvedEnv, *expectedEnv) {
		t.Fatalf("unknown solution to zero-temp system: got\n%s, expected\n%s", (&solvedEnv).String(), expectedEnv.String())
	}
}

func TestGc0(t *testing.T) {
	cacheFileName := "zerotemp_test_gc0_cache.json"
	flag.Parse()
	tolerances := []float64{1e-6, 1e-6, 1e-6}
	system := NewZeroTempSystem(tolerances)
	env, err := EnvironmentFromFile("zerotemp_test_gc0.json")
	if err != nil {
		t.Fatal(err)
	}
	env.Initialize()
	var solvedEnv Environment // not sure if the seperate declaration is needed
	if !(*cached) {
		solution, err := system.Solve(*env)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnv = solution.(Environment)
		solvedEnv.WriteToFile(cacheFileName)
	} else {
		cacheEnv, err := EnvironmentFromFile(cacheFileName)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnv = *cacheEnv
	}
	k := Vector2{0.0 * math.Pi, 0.0 * math.Pi}
	poles, err := ZeroTempGreenPolePoint(solvedEnv, k)
	if err != nil {
		t.Fatal(err)
	}
	if len(poles) != 4 {
		t.Fatal("did not get expected number of poles")
	}
	expected := []string{"-3.8862440783169987", "3.86016952787911", "3.9111002349865496", "7.560058812941273"}
	for i, p := range poles {
		if fmt.Sprintf("%v", p) != expected[i] {
			t.Fatal("did not get expected pole value")
		}
	}
	split := 0.01
	poleCurve := func(x float64) Vector2 {
		val := 0.5*math.Pi + split*(2*x-1)
		return Vector2{val, val}
	}
	ZeroTempPlotPoleCurve(solvedEnv, poleCurve, 2048, "zerotemp.testignore.polecurve")
}

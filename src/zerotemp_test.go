package polecalc

import (
	"testing"
	"reflect"
)

func TestKnownZeroTempSystem(t *testing.T) {
	envStr := "{\"GridLength\":8,\"ImG0Bins\":0,\"NumProcs\":1,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"J\":0,\"A\":0,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"Lambda\":0,\"D1\":0.05777149373506872,\"Mu\":-0.18330570279347036,\"F0\":0.12945949461029926,\"EpsilonMin\":-1.8}"
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
	//	println(solvedEnv.ZeroTempErrors())
}

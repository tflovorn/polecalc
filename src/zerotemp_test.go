package polecalc

import "testing"

func TestKnownZeroTempSystem(t *testing.T) {
	knownResult := "{\"GridLength\":16,\"NumProcs\":1,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.052859796549526133,\"Mu\":-0.21741415401597314,\"F0\":0.13102201547770487,\"EpsilonMin\":-1.8}"
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
	if (&solvedEnv).String() != knownResult {
		t.Fatal("unknown solution to zero-temp system: got ", (&solvedEnv).String())
	}
}

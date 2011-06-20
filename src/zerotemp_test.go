package polecalc

import "testing"

func TestKnownZeroTempSystem(t *testing.T) {
	knownResult := "{\"GridLength\":64,\"NumProcs\":1,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.053352856524156486,\"Mu\":-0.23959545685146372,\"F0\":0.1326082268107227,\"EpsilonMin\":-1.8}"
	tolerances := []float64{1e-6, 1e-6, 1e-6}
	system := NewZeroTempSystem(tolerances)
	env, err := EnvironmentFromFile("environment_test.json")
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
		t.Fatal("unknown solution to zero-temp system")
	}
}

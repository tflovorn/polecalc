package polecalc

import "testing"

func TestKnownZeroTempSystem(t *testing.T) {
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
	println((&solvedEnv).String())
}

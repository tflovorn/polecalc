package polecalc

import (
	"reflect"
	"testing"
)

// Does Environment correctly convert to and from JSON?
func TestEnvironmentConvert(t *testing.T) {
	envStr := "{\"GridLength\":64,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.1,\"Mu\":0.1,\"F0\":0.1,\"EpsilonMin\":-1.8}"
	expectedEnv, err := EnvironmentFromString(envStr)
	if err != nil {
		t.Fatal(err)
	}
	expectedEnv.Initialize()
	env, err := EnvironmentFromFile("environment_test.json")
	if err != nil {
		t.Fatal(err)
	}
	env.Initialize()
	if !reflect.DeepEqual(env, expectedEnv) {
		println(env.String())
		t.Fatal("Environment does not match known value")
	}
}

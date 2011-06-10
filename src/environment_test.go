package polecalc

import (
	"testing"
	"fmt"
)

// Does Environment correctly convert to and from JSON?
func TestEnvironmentConvert(t *testing.T) {
	expectedEnv := "{\"GridLength\":64,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.1,\"Mu\":0.1,\"F0\":0.1}"
	env, err := EnvironmentFromFile("environment_test.json")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	env.Initialize()
	if env.String() != expectedEnv {
		t.FailNow()
	}
}

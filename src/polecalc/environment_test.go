package polecalc

import (
	"testing"
	"fmt"
)

// Does Environment correctly convert to and from JSON?
func TestEnvironmentConvert(t *testing.T) {
	env, err := EnvironmentFromFile("environment_test.json")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	env.Initialize()
	fmt.Println(env)
}

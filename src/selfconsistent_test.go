package polecalc

import (
	"testing"
	"os"
	"fmt"
)

type LinearEquation struct {
	root float64
}

func (eq LinearEquation) AbsError(args interface{}) float64 {
	return eq.root - *args.(*float64)
}

func (eq LinearEquation) SetArguments(x float64, args *interface{}) {
	*(*args).(*float64) = x
}

func (eq LinearEquation) Range(args interface{}) (float64, float64, os.Error) {
	return eq.root - 2 * eq.root, eq.root + 2 * eq.root, nil
}

// Does Solve() correctly find the root of an equation?
func TestSolve(t *testing.T) {
	root := 10.0
	eq := LinearEquation{root}
	guess := 0.0
	solution, err := Solve(eq, &guess)
	if err != nil {
		t.Fatal(err)
	}
	if solution == nil {
		t.Fatalf("got nil solution")
	}
	if guess - root > MachEpsFloat64() {
		t.Fatalf("solution not found to expected precision")
	}
}

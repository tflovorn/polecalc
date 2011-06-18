package polecalc

import (
	"testing"
	"os"
	"math"
)

type LinearEquation struct {
	root float64
	myVar string
}

func (eq LinearEquation) AbsError(args interface{}) float64 {
	vars := args.(map[string]float64)
	return eq.root - vars[eq.myVar]
}

func (eq LinearEquation) SetArguments(x float64, args *interface{}) {
	vars := (*args).(map[string]float64)
	vars[eq.myVar] = x
}

func (eq LinearEquation) Range(args interface{}) (float64, float64, os.Error) {
	return eq.root - 2 * eq.root, eq.root + 2 * eq.root, nil
}

// Does Solve() correctly find the root of a linear equation?
func TestSolve(t *testing.T) {
	root := 10.0
	eq := LinearEquation{root, "uno"}
	guess := make(map[string]float64)
	guess["uno"] = 0.0
	solution, err := Solve(eq, guess)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if solution == nil {
		t.Fatalf("got nil solution")
	}
	vals := solution.(map[string]float64)
	if math.Fabs(vals["uno"] - root) > MachEpsFloat64() {
		t.Fatalf("solution not found to expected precision")
	}
}

// Does SelfConsistentSystem.Solve() correctly solve a system of non-coupled
// linear equations?
func TestSystemSolve(t *testing.T) {
	root1, root2 := 10.0, 20.0
	eq1, eq2 := LinearEquation{root1, "uno"}, LinearEquation{root2, "dos"}
	tol1, tol2 := 1e-9, 1e-9
	system := &SelfConsistentSystem{[]SelfConsistentEquation{eq1, eq2}, []float64{tol1, tol2}}
	args := make(map[string]float64)
	args["uno"] = 0.0
	args["dos"] = 0.0
	solution, err := system.Solve(args)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if solution == nil {
		t.Fatalf("got nil solution")
	}
	vals := solution.(map[string]float64)
	if math.Fabs(vals["uno"] - root1) > MachEpsFloat64() {
		t.Fatalf("solution uno not found to expected precision")
	}
	if math.Fabs(vals["dos"] - root2) > MachEpsFloat64() {
		t.Fatalf("solution dos not found to expected precision")
	}
}

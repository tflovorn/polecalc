package polecalc

import "os"

type Arguments interface{}

// One-parameter scalar self-consistent equation
type SelfConsistentEquation interface {
	// return the absolute error associated with this equation under the given Arguments
	AbsError(args Arguments) float64
	// set the appropriate variable in args to value
	SetArguments(value float64, args *Arguments)
	// range of possible values for SetArguments
	Range(args Arguments) (float64, float64, os.Error)
}

// Return an Arguments which solves eq to the given tolerance
func Solve(eq SelfConsistentEquation, args Arguments, tolerance float64) (Arguments, os.Error) {
	eqError := func(value float64) float64 {
		eq.SetArguments(value, &args)
		return eq.AbsError(args)
	}
	// if eq only has one root, left and right must bracket it
	// --- todo: need to verify this ---
	left, right, err := eq.Range(args)
	if err != nil {
		return args, err
	}
	solution, err := BisectionFullPrecision(eqError, left, right)
	if err != nil {
		return args, err
	}
	eq.SetArguments(solution, &args)
	return args, nil
}

// A group of self-consistent equations which may be coupled and must all be
// solved for the group to be considered solved.
type SelfConsistentSystem struct {
	Equations  []SelfConsistentEquation
	Tolerances []float64
}

// Solve the self-consistent system, returning the resulting Arguments
func (system *SelfConsistentSystem) Solve(args Arguments) (Arguments, os.Error) {
	i := 0
	for !system.IsSolved(args) {
		// this should never be true: if it is, failed to iterate
		if i >= len(system.Equations) {
			panic("self-consistent system overran bounds")
		}
		// set args to the value that solves the equation
		newEnv, err := Solve(system.Equations[i], args, system.Tolerances[i])
		if err != nil {
			return args, err
		}
		args = newEnv
		// check if we need to iterate
		if !system.solvedUpTo(args, i) {
			// previous equations have been disturbed; restart
			i = 0
		} else {
			// ok to continue to next equation
			i++
		}
	}
	return args, nil
}

// Check if the first (maxIndex + 1) equations are solved
func (system *SelfConsistentSystem) solvedUpTo(args Arguments, maxIndex int) bool {
	for i, eq := range system.Equations {
		if i > maxIndex {
			break
		}
		if eq.AbsError(args) > system.Tolerances[i] {
			return false
		}
	}
	return true
}

// Are all the self-consistent equations solved?
func (system *SelfConsistentSystem) IsSolved(args Arguments) bool {
	if len(system.Equations) == 0 {
		return true
	}
	return system.solvedUpTo(args, len(system.Equations)-1)
}

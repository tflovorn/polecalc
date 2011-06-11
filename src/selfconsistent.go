package polecalc

import "os"

// One-parameter scalar self-consistent equation
type SelfConsistentEquation interface {
	// return the absolute error associated with this equation under the given Environment
	AbsError(env Environment) float64
	// set the appropriate variable in env to value
	SetEnvironment(value float64, env *Environment)
	// range of possible values for SetEnvironment
	Range(env Environment) (float64, float64, os.Error)
}

// Return an Environment which solves eq to the given tolerance
func Solve(eq SelfConsistentEquation, env Environment, tolerance float64) (Environment, os.Error) {
	eqError := func(value float64) float64 {
		eq.SetEnvironment(value, &env)
		return eq.AbsError(env)
	}
	// if eq only has one root, left and right must bracket it
	// --- todo: need to verify this ---
	left, right, err := eq.Range(env)
	if err != nil {
		return env, err
	}
	solution, err := BisectionFullPrecision(eqError, left, right)
	if err != nil {
		return env, err
	}
	eq.SetEnvironment(solution, &env)
	return env, nil
}

// A group of self-consistent equations which may be coupled and must all be
// solved for the group to be considered solved.
type SelfConsistentSystem struct {
	Equations  []SelfConsistentEquation
	Tolerances []float64
}

// Solve the self-consistent system, returning the resulting Environment
func (system *SelfConsistentSystem) Solve(env Environment) (Environment, os.Error) {
	i := 0
	for !system.IsSolved(env) {
		// this should never be true: if it is, failed to iterate
		if i >= len(system.Equations) {
			panic("self-consistent system overran bounds")
		}
		// set env to the value that solves the equation
		newEnv, err := Solve(system.Equations[i], env, system.Tolerances[i])
		if err != nil {
			return env, err
		}
		env = newEnv
		// check if we need to iterate
		if !system.solvedUpTo(env, i) {
			// previous equations have been disturbed; restart
			i = 0
		} else {
			// ok to continue to next equation
			i++
		}
	}
	return env, nil
}

// Check if the first (maxIndex + 1) equations are solved
func (system *SelfConsistentSystem) solvedUpTo(env Environment, maxIndex int) bool {
	for i, eq := range system.Equations {
		if i > maxIndex {
			break
		}
		if eq.AbsError(env) > system.Tolerances[i] {
			return false
		}
	}
	return true
}

// Are all the self-consistent equations solved?
func (system *SelfConsistentSystem) IsSolved(env Environment) bool {
	if len(system.Equations) == 0 {
		return true
	}
	return system.solvedUpTo(env, len(system.Equations)-1)
}

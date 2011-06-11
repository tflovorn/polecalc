package polecalc

// Error returned when the self-consistent equation cannot be solved
type SelfConsistentError struct {
	ResponsibleEnv Environment
	Reason         string
}

func (e *SelfConsistentError) String() string {
	return e.Reason
}

// One-parameter scalar self-consistent equation
type SelfConsistentEquation interface {
	// return the absolute error associated with this equation under the given Environment
	AbsError(env Environment) float64
	// set the appropriate variable in env to value
	SetEnvironment(value float64, env *Environment)
}

// Return an Environment which solves eq to the given tolerance
func Solve(eq SelfConsistentEquation, env Environment, tolerance float64) (Environment, *SelfConsistentError) {
	error := func(value float64) float64 {
		eq.SetEnvironment(value, &env)
		return eq.AbsError(env)
	}
	// -- find solution --
	error(0)        // dummy
	solution := 0.0 // dummy
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
func (system *SelfConsistentSystem) Solve(env Environment) (Environment, *SelfConsistentError) {
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

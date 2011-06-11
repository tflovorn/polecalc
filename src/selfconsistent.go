package polecalc

// Error returned when the self-consistent equation cannot be solved
type SelfConsistentError struct {
	Reason string
}

func (e *SelfConsistentError) String() string {
	return e.Reason
}

// One-parameter scalar self-consistent equation
type SelfConsistentEquation interface {
	AbsError(value float64, env Environment) float64
}

// Return an Environment which solves eq to the given tolerance
func Solve(eq SelfConsistentEquation, env Environment, tolerance float64) (*Environment, *SelfConsistentError) {
	return &env, nil
}

// A group of self-consistent equations which may be coupled and must all be
// solved for the group to be considered solved.
type SelfConsistentSystem struct {
	equations []SelfConsistentEquation
	tolerances []float64
}

func (system *SelfConsistentSystem) Solve(env Environment) (*Environment, *SelfConsistentError) {
	var i uint = 0
	for !system.IsSolved(env) {
		// -- panic(or return error?) if i >= len(system.equations) --

		// -- solve the equation at index i --

		// check if we need to iterate
		if !system.solvedUpTo(env, i) {
			// previous equations have been disturbed; restart
			i = 0
		} else {
			// ok to continue to next equation
			i++
		}
	}
	return &env, nil
}

func (system *SelfConsistentSystem) solvedUpTo(env Environment, index uint) bool {
	return false
}

func (system *SelfConsistentSystem) IsSolved(env Environment) bool {
	if len(system.equations) == 0 {
		return true
	}
	return system.solvedUpTo(env, uint(len(system.equations) - 1))
}

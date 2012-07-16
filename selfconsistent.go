package polecalc

import "math"

// One-parameter scalar self-consistent equation
type SelfConsistentEquation interface {
	// return the absolute error associated with this equation under the given interface{}
	AbsError(args interface{}) float64
	// set the appropriate variable in args to value
	SetArguments(value float64, args interface{}) interface{}
	// range of possible values for SetArguments
	Range(args interface{}) (float64, float64, error)
}

// Return an interface{} which solves eq to tolerance of BisectionFullPrecision
func Solve(eq SelfConsistentEquation, args interface{}) (interface{}, error) {
	eqError := func(value float64) float64 {
		args = eq.SetArguments(value, args)
		return eq.AbsError(args)
	}
	leftEdge, rightEdge, err := eq.Range(args)
	if err != nil {
		return args, err
	}
	left, right, err := FindBracket(eqError, leftEdge, rightEdge)
	if err != nil {
		return args, err
	}
	solution, err := BisectionFullPrecision(eqError, left, right)
	if err != nil {
		return args, err
	}
	args = eq.SetArguments(solution, args)
	return args, nil
}

// Return a slice of interface{}'s which solve eq
func MultiSolve(eq SelfConsistentEquation, args interface{}) ([]interface{}, error) {
	eqError := func(value float64) float64 {
		args = eq.SetArguments(value, args)
		return eq.AbsError(args)
	}
	leftEdge, rightEdge, err := eq.Range(args)
	if err != nil {
		return nil, err
	}
	brackets, err := MultiBracket(eqError, leftEdge, rightEdge)
	if err != nil {
		return nil, err
	}
	solutions := []interface{}{}
	for _, bracket := range brackets {
		left, right := bracket[0], bracket[1]
		solution, err := BisectionFullPrecision(eqError, left, right)
		if err != nil {
			return solutions, err
		}
		solutions = append(solutions, eq.SetArguments(solution, args))
	}
	return solutions, nil
}

// A group of self-consistent equations which may be coupled and must all be
// solved for the group to be considered solved.
type SelfConsistentSystem struct {
	Equations  []SelfConsistentEquation
	Tolerances []float64
}

// Solve the self-consistent system, returning the resulting interface{}
func (system *SelfConsistentSystem) Solve(args interface{}) (interface{}, error) {
	i := 0
	for !system.IsSolved(args) {
		// this should never be true: if it is, failed to iterate
		if i >= len(system.Equations) {
			panic("self-consistent system overran bounds")
		}
		// set args to the value that solves the equation
		newEnv, err := Solve(system.Equations[i], args)
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
func (system *SelfConsistentSystem) solvedUpTo(args interface{}, maxIndex int) bool {
	for i, eq := range system.Equations {
		if i > maxIndex {
			break
		}
		if math.Abs(eq.AbsError(args)) > system.Tolerances[i] {
			return false
		}
	}
	return true
}

// Are all the self-consistent equations solved?
func (system *SelfConsistentSystem) IsSolved(args interface{}) bool {
	if len(system.Equations) == 0 {
		return true
	}
	return system.solvedUpTo(args, len(system.Equations)-1)
}

package polecalc

import "os"

// (should probably set these constants through a configuration method)
// Number of steps to take in the first attempt to find a bracket.
const InitialBracketNumber uint = 32
// If using more steps than this to find a bracket, stop.
const MaxBracketNumber uint = 8192 // 8 iterations from 32 (32 * 2^8)

// Find all pairs of points which bracket roots of f between left and right.
func MultiBracket(f Func1D, left, right float64, bracketNum uint) ([][]float64, os.Error) {
	return bracketHelper(f, left, right, InitialBracketNumber, -1)
}

// Find a pair of points which bracket a root of f between left and right.
func FindBracket(f Func1D, left, right float64) (float64, float64, os.Error) {
	bracket, err := bracketHelper(f, left, right, InitialBracketNumber, 1)
	if err != nil {
		return 0.0, 0.0, err
	}
	bl, br := bracket[0][0], bracket[0][1]
	return bl, br, err
}

// Use a number of divisions equal to bracketNum to find a root.
// If maxBrackets <= 0, get as many brackets as possible.
func bracketHelper(f Func1D, left, right float64, bracketNum uint, maxBrackets int) ([][]float64, os.Error) {
	if left == right {
		return nil, os.NewError("bracket error: must give two distinct points to find bracket")
	}
	if left > right {
		left, right = right, left
	}
	xs := MakeRange(left, right, bracketNum)
	scale := xs[1]-xs[0]
	brackets := [][]float64{}
	for i, _ := range xs {
		// only get as many brackets as requested
		if maxBrackets > 0 && len(brackets) >= maxBrackets {
			return brackets, nil
		}
		// don't check [endpoint, endpoint+scale] bracket
		if i >= len(xs)-1 {
			break
		}
		// check function values
		fa, fb := f(xs[i]), f(xs[i+1])
		if FuzzyEqual(fb, 0.0) {
			brackets = append(brackets, []float64{xs[i], xs[i+1] + scale})
		} else if FuzzyEqual(fa, 0.0) {
			brackets = append(brackets, []float64{xs[i] - scale, xs[i+1]})
		}
		if !sameSign(fa, fb) {
			brackets = append(brackets, []float64{xs[i], xs[i+1]})
		}
	}
	// overshot bounds if without finding bracket if we get here
	if bracketNum >= MaxBracketNumber {
		// too many divisions
		return nil, os.NewError("cannot find bracket")
	}
	// try again with smaller divisions
	return bracketHelper(f, left, right, bracketNum*2, maxBrackets)
}

// If x and y don't have the same sign, we know they bracket a root.
func sameSign(x, y float64) bool {
	if x >= 0 && y >= 0 {
		return true
	} else if x <= 0 && y <= 0 {
		return true
	}
	return false
}

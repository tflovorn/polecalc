package polecalc

import "os"

// Lazy hack - this really should be a parameter in Environment but using that
// from selfconsistent would require all args passed to it to know this value
// -- need an interface for selfconsistent args --
const InitialBracketNumber uint = 32
const MaxBracketNumber uint = 8192 // 8 iterations from 32 (32 * 2^8)

func FindBracket(f Func1D, left, right float64) (float64, float64, os.Error) {
	if left == right {
		return 0.0, 0.0, os.NewError("must give two distinct points to find bracket")
	}
	if left > right {
		left, right = right, left
	}
	return bracketHelper(f, left, right, InitialBracketNumber)
}

func bracketHelper(f Func1D, left, right float64, bracketNum uint) (float64, float64, os.Error) {
	scale := (right - left) / float64(bracketNum)
	a, b := left, left+scale
	for sameSign(f(a), f(b)) {
		if f(b) == 0.0 {
			return a, b + scale, nil
		} else if f(a) == 0.0 {
			return a - scale, b, nil
		}
		a, b = b, b+scale
		if a >= right || b > right {
			// overshot bounds, need to bail or iterate
			if bracketNum >= MaxBracketNumber {
				return 0.0, 0.0, os.NewError("cannot find bracket")
			} else {
				return bracketHelper(f, left, right, bracketNum*2)
			}
		}
	}
	return a, b, nil
}

func sameSign(x, y float64) bool {
	if x >= 0 && y >= 0 {
		return true
	} else if x <= 0 && y <= 0 {
		return true
	}
	return false
}

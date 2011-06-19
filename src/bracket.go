package polecalc

import "os"

// Lazy hack - this really should be a parameter in Environment but using that
// from selfconsistent would require all args passed to it to know this value
// -- need an interface for selfconsistent args --
func BracketNum() uint {
	return 64
}

func FindBracket(f Func1D, left, right float64) (float64, float64, os.Error) {
	sameSign := func(x, y float64) bool {
		if x >= 0 && y >= 0 {
			return true
		} else if x <= 0 && y <= 0 {
			return true
		}
		return false
	}
	if left == right {
		return 0.0, 0.0, os.NewError("must give two distinct points to find bracket")
	}
	if left > right {
		left, right = right, left
	}
	bracketNum := BracketNum()
	scale := (right - left) / float64(bracketNum)
	a, b := left, left+scale
	for sameSign(f(a), f(b)) {
		if f(b) == 0.0 {
			//println("f(b) is 0")
			return a, b + scale, nil
		} else if f(a) == 0.0 {
			//println("f(a) is 0")
			return a - scale, b, nil
		}
		a, b = b, b+scale
		if a >= right || b > right {
			return 0.0, 0.0, os.NewError("cannot find bracket")
		}
	}
	//println(f(a), f(b), a, b)
	return a, b, nil
}

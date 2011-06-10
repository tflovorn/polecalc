package polecalc

import "math"

// Error returned when a root cannot be found
type RootError struct {
	Reason string
}

func (e *RootError) String() string {
	return e.Reason
}

// One-dimensional scalar function - root can be found by bisection
type Func1D func(float64) float64

// Find the root of f in the interval (left, right) to precision epsilon using the bisection method
func SolveBisection(f Func1D, left, right, epsilon float64) (float64, *RootError) {
	if left > right {
		left, right = right, left
	}
	var error *RootError
	for math.Fabs(right - left) > 2*epsilon {
		left, right, error = BisectionIterate(f, left, right)
		if error != nil {
			return (left + right) / 2.0, error
		}
	}
	return (left + right) / 2.0, nil
}

// Provide the next iteration of the bisection method for f on the interval (left, right)
func BisectionIterate(f Func1D, left, right float64) (float64, float64, *RootError) {
	if f(left) * f(right) > 0 {
		return left, right, &RootError{"arguments do not bracket a root"}
	}
	midpoint := (left + right) / 2.0
	fl, fm, fr := f(left), f(midpoint), f(right)
	if fl == 0 {
		// left side is the root
		return left, left, nil
	} else if fr == 0 {
		// right side is the root
		return right, right, nil
	} else if (fl > 0 && fm < 0) || (fl < 0 && fm > 0) {
		// root is in the left half
		return left, midpoint, nil
	} else if (fm > 0 && fr < 0) || (fm < 0 && fr > 0) {
		// root is in the right half
		return midpoint, right, nil
	}
	// midpoint must be the root
	return midpoint, midpoint, nil
}

// Solve f for the root in interval (a, b) up to machine precision using bisection
// Cribbed from implementation on Wikipedia page 'Bisection method'
func BisectionFullPrecision(f Func1D, a, b float64) (float64, *RootError) {
	fa, fb := f(a), f(b)
	if !((fa >= 0 && fb <= 0) || (fa <= 0 && fb >= 0)) {
		// no root bracketed
		return 0, &RootError{"arguments do not bracket a root"}
	}
	var lo, hi float64
	if fa <= 0 {
		lo, hi = a, b
	} else {
		lo, hi = b, a
	}
	mid := lo + (hi - lo) / 2.0
	for mid != lo && mid != hi {
		if f(mid) <= 0 {
			lo = mid
		} else {
			hi = mid
		}
		mid = lo + (hi - lo) / 2.0
	}
	return mid, nil
}

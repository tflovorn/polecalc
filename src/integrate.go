package polecalc

import (
	"os"
	"fmt"
)

// Integrate the cubic spline interpolation of y from x = left to x = right.  
// xs is an ordered slice of equally spaced x values.
// ys is a slice of the corresponding y values.
// Assume left >= xs[0] and right <= xs[len(xs)-1].
func SplineIntegral(xs, ys []float64, left, right float64) (float64, os.Error) {
	// if the arguments are reversed, we need a minus sign later
	sign := 1.0
	if left > right {
		sign = -1.0
		left, right = right, left
	}
	// make the spline
	s, err := NewCubicSpline(xs, ys)
	if err != nil {
		return 0.0, err
	}
	xMin, xMax := s.Range()
	// can't integrate if left or right is outside of interpolation range
	eps := SplineExtrapolationDistance
	if (xMin-left > eps) || (right-xMax > eps) {
		return 0.0, os.NewError("SplineIntegral error: integral arguments out of bounds")
	}
	// k and q are the first and last indices for integration
	k, q := s.indexOf(left), s.indexOf(right)
	first := s.antiDeriv(k, xs[k+1]) - s.antiDeriv(k, left)
	last := s.antiDeriv(q, right) - s.antiDeriv(q, xs[q])
	// add upp all the middle segment integrals
	middle, compensate := 0.0, 0.0
	for i := k + 1; i < q; i++ {
		integral := s.antiDeriv(i, xs[i+1]) - s.antiDeriv(i, xs[i])
		middle, compensate = KahanSum(integral, middle, compensate)
	}
	return sign * (first + middle + last), nil
}

// Principal value integral of f(x) from x = a to x = b.  Assume there is a 
// pole at x = w and do cubic spline integrals in the appropriate spots, 
// staying a distance eps away from the pole.  Use n points for the cubic
// spline on each side of the pole.
func PvIntegral(f Func1DError, a, b, w, eps float64, n uint) (float64, os.Error) {
	// can't integrate if a boundary is on top of the pole
	if a == w || b == w {
		return 0.0, fmt.Errorf("PvIntegral error: pole (%f) equals a boundary (%f, %f)", w, a, b)
	}
	// if the pole is within eps of a boundary, reduce eps
	drop := 10.0
	for w > a && w-a < eps {
		eps = (w - a) / drop
	}
	for b > w && b-w < eps {
		eps = (b - w) / drop
	}
	// if the bounds were given out of order, we need a minus sign later
	sign := 1.0
	if a > b {
		sign = -1.0
		a, b = b, a
	}
	// pole is fully inside integration region
	if a <= w && w <= b {
		// avoid the interval [wl, wr]
		wl, wr := w-eps, w+eps
		// left and right sets of x points
		xls, xrs := MakeRange(a, wl, n), MakeRange(wr, b, n)
		// left and right sets of y points
		yls, yrs := make([]float64, n), make([]float64, n)
		for i := uint(0); i < n; i++ {
			yl, err := f(xls[i])
			if err != nil {
				return 0.0, err
			}
			yr, err := f(xrs[i])
			if err != nil {
				return 0.0, err
			}
			yls[i], yrs[i] = yl, yr
		}
		// do the integrals
		// if the pole is very close to the boundary, integral ~ 0
		var leftInt, rightInt float64
		if wl < a || FuzzyEqual(w-a, 0.0) {
			leftInt = 0.0
		} else {
			var err os.Error
			leftInt, err = SplineIntegral(xls, yls, a, wl)
			if err != nil {
				return 0.0, err
			}
		}
		if wr > b || FuzzyEqual(b-w, 0.0) {
			rightInt = 0.0
		} else {
			var err os.Error
			rightInt, err = SplineIntegral(xrs, yrs, wr, b)
			if err != nil {
				return 0.0, err
			}
		}
		return sign * (leftInt + rightInt), nil
	}
	// pole is fully outside the integration region
	// x values take the entire range
	xs := MakeRange(a, b, 2*n)
	// associated y values
	ys := make([]float64, 2*n)
	for i := uint(0); i < 2*n; i++ {
		y, err := f(xs[i])
		if err != nil {
			return 0.0, err
		}
		ys[i] = y
	}
	// only one integral to do
	integral, err := SplineIntegral(xs, ys, a, b)
	if err != nil {
		return 0.0, err
	}
	return sign * integral, nil
}

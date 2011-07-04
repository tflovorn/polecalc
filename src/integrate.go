package polecalc

import (
	"os"
	"math"
)

// Integrate the cubic spline interpolation of y from x = left to x = right.  
// xs is an ordered slice of equally spaced x values.
// ys is a slice of the corresponding y values.
// Assume left < right, left >= xs[0], and right <= xs[len(xs)-1].
func SplineIntegral(xs, ys []float64, left, right float64) (float64, os.Error) {
	// if the arguments are reversed, we need a minus sign later
	sign := 1.0
	if left > right {
		sign = -1.0
		left, right = right, left
	}
	s, err := NewCubicSpline(xs, ys)
	if err != nil {
		return 0.0, err
	}
	xMin, xMax := s.Range()
	if left < xMin || right > xMax {
		return 0.0, os.NewError("integral arguments out of bounds")
	}
	k, q := s.indexOf(left), s.indexOf(right)
	first := s.antiDeriv(k, xs[k+1]) - s.antiDeriv(k, left)
	middle, compensate := 0.0, 0.0
	for i := k + 1; i < q; i++ {
		integral := s.antiDeriv(i, xs[i+1]) - s.antiDeriv(i, xs[i])
		middle, compensate = KahanSum(integral, middle, compensate)
	}
	last := s.antiDeriv(q, right) - s.antiDeriv(q, xs[q])
	return sign*(first + middle + last), nil
}

// Principal value integral of f(x) from x = a to x = b.  Assume there is a 
// pole at x = w and do cubic spline integrals in the appropriate spots, 
// staying a distance eps away from the pole.  Use n points for the cubic
// spline on each side of the pole.
func PvIntegral(f Func1D, a, b, w, eps float64, n uint) (float64, os.Error) {
	// can't integrate if a boundary is on top of the pole
	if math.Fabs(a - w) < eps || math.Fabs(b - w) < eps {
		return 0.0, os.NewError("pv integral error: boundary too close to pole")
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
		wl, wr := w - eps, w + eps
		// left and right sets of x points
		xls, xrs := MakeRange(a, wl, n), MakeRange(wr, b, n)
		// left and right sets of y points
		yls, yrs := make([]float64, n), make([]float64, n)
		for i := uint(0); i < n; i++ {
			yls[i], yrs[i] = f(xls[i]), f(xrs[i])
		}
		// do the integrals
		leftInt, err := SplineIntegral(xls, yls, a, wl)
		if err != nil {
			return 0.0, err
		}
		rightInt, err := SplineIntegral(xrs, yrs, wr, b)
		if err != nil {
			return 0.0, err
		}
		return sign*(leftInt + rightInt), nil
	}
	// pole is fully outside the integration region
	// x values take the entire range
	xs := MakeRange(a, b, 2*n)
	// associated y values
	ys := make([]float64, 2*n)
	for i := uint(0); i < 2*n; i++ {
		ys[i] = f(xs[i])
	}
	// only one integral to do
	integral, err := SplineIntegral(xs, ys, a, b)
	if err != nil {
		return 0.0, err
	}
	return sign*integral, nil
}

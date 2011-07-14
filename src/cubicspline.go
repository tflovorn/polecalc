// Cubic spline interpolation and integration of discrete points
// Based on description of algorithm found at:
// http://web.archive.org/web/20090408054627/http://online.redwoods.cc.ca.us/instruct/darnold/laproj/Fall98/SkyMeg/Proj.PDF
package polecalc

import (
	"os"
	"math"
)

const SplineExtrapolationDistance = 1e-14

// Cubic spline is defined by the set of functions s_i given on the intervals
// between the values of xs: s_i(x) is defined on [xs[i],xs[i+1]).
// s_i(x) = a[i](x-xs[i])^3 + b[i](x-xs[i])^2 + c[i](x-xs[i]) + d[i]
type CubicSpline struct {
	a, b, c, d []float64 // length n - 1
	xs         []float64 // length n
}

// Return a pointer to a cubic spline interpolating y = f(x).
// xs is an ordered slice of equally spaced x values.
// ys is a slice of the corresponding y values.
func NewCubicSpline(xs, ys []float64) (*CubicSpline, os.Error) {
	// xs and ys must have the same length
	if len(xs) != len(ys) {
		return nil, os.NewError("input slices must be the same length")
	}
	// must have at least three points
	if len(xs) < 3 {
		return nil, os.NewError("not enough points for cubic spline")
	}
	// xs must be ordered
	if !inAscendingOrder(xs) {
		return nil, os.NewError("xs must be in ascending order")
	}
	spline := new(CubicSpline)
	spline.xs = xs
	spline.a, spline.b, spline.c, spline.d = splineCoeffs(xs, ys)
	return spline, nil
}

// Value of the interpolated function S(x) at x
// Will panic if x is outside the interpolation range
func (s *CubicSpline) At(x float64) float64 {
	xMin, xMax := s.Range()
	eps := SplineExtrapolationDistance
	if (xMin-x > eps) || (x-xMax > eps) {
		panic("accessing cubic spline out of bounds")
	}
	i := s.indexOf(x)
	return s.splineAt(i, x)
}

// Individual spline functions si(x) at index i, position x
// Assumes i > 0 and x is in the appropriate range for si
func (s *CubicSpline) splineAt(i int, x float64) float64 {
	dx := x - s.xs[i]
	return s.a[i]*math.Pow(dx, 3.0) + s.b[i]*math.Pow(dx, 2.0) + s.c[i]*dx + s.d[i]
}

// Antiderivative of the spline functions (with integration constant = 0)
// Makes the same assumptions as splineAt.
func (s *CubicSpline) antiDeriv(i int, x float64) float64 {
	dx := x - s.xs[i]
	return s.a[i]*math.Pow(dx, 4.0)/4 + s.b[i]*math.Pow(dx, 3.0)/3 + s.c[i]*math.Pow(dx, 2.0)/2 + s.d[i]*x
}

// Interpolation range of the spline
func (s *CubicSpline) Range() (float64, float64) {
	n := len(s.xs)
	return s.xs[0], s.xs[n-1]
}

// Return the index i such that xs[i] <= x < xs[i+1]
// Assume x is within the bounds of the spline
// i will be between 0 and n-2 where n is len(s.xs)
func (s *CubicSpline) indexOf(x float64) int {
	xMin, xMax := s.Range()
	// -1 to accomodate having one less interpolating function than the
	// number of points
	numSplines := len(s.xs) - 1
	step := (xMax - xMin) / float64(numSplines)
	i := int(math.Floor((x - xMin) / step))
	// final point does not jump to a new spline
	if i == numSplines && FuzzierEqual(x, xMax) {
		return i - 1
	}
	return i
}

// Find the cubic spline coefficients corresponding to the given points
func splineCoeffs(xs []float64, ys []float64) ([]float64, []float64, []float64, []float64) {
	n := len(xs)
	h := xs[1] - xs[0]
	M := solveNaturalSplineEqn(h, ys)
	a, b, c, d := make([]float64, n-1), make([]float64, n-1), make([]float64, n-1), make([]float64, n-1)
	for i, _ := range a {
		a[i] = (M[i+1] - M[i]) / (6 * h)
		b[i] = M[i] / 2
		c[i] = (ys[i+1]-ys[i])/h - h*(M[i+1]+2*M[i])/6
		d[i] = ys[i]
	}
	return a, b, c, d
}

// Solve the tridiagonal matrix equation for M, a slice of second derivative 
// values used in calculating the interpolating function coefficients.
func solveNaturalSplineEqn(h float64, ys []float64) []float64 {
	M := TridiagonalSolve(splineTriDiagInit(h, ys))
	// natural spline condition
	M = PadLeftWith0(M)
	M = append(M, 0)
	return M
}

// Initialize the tridiagonal matrix
func splineTriDiagInit(h float64, ys []float64) ([]float64, []float64, []float64, []float64) {
	n := len(ys)
	a, b, c, d := make([]float64, n-1), make([]float64, n-2), make([]float64, n-1), make([]float64, n-2)
	for i, _ := range a {
		a[i] = 1
		c[i] = 1
	}
	for i, _ := range b {
		b[i] = 4
		d[i] = (6 / (h * h)) * (ys[i] - 2*ys[i+1] + ys[i+2])
	}
	return a, b, c, d
}

// check if xs is in ascending order
func inAscendingOrder(xs []float64) bool {
	for i, val := range xs {
		if i != 0 && xs[i-1] > val {
			return false
		}
	}
	return true
}

// Solve a tridiagonal matrix equation
// (in order to support cubic spline interpolation)
// Based on Matlab implementation on Wikipedia
// (http://en.wikipedia.org/wiki/Tridiagonal_matrix_algorithm#Implementation_in_Matlab)
package polecalc

// a: sub-diagonal (below main), input with length n - 1
// b: main diagonal, input with length n
// c: sup-diagonal (above main), input with length n - 1
// d: right side of equation, input with length n
func TridiagonalSolve(a, b, c, d []float64) []float64 {
	n := len(b)
	x := make([]float64, n)
	// algorithm expects these to be the same length as b
	a = PadLeftWith0(a)
	c = append(c, 0)
	// forward sweep: modify coefficients
	// first row
	c[0] = c[0] / b[0]
	d[0] = d[0] / b[0]
	// remaining rows
	for i := 1; i < n; i++ {
		id := 1 / (b[i] - c[i-1]*a[i])
		c[i] = c[i] * id
		d[i] = (d[i] - d[i-1]*a[i]) * id
	}
	// back substitute: get solution
	// final row
	x[n-1] = d[n-1]
	// remaining rows
	for i := n - 2; i >= 0; i-- {
		x[i] = d[i] - c[i]*x[i+1]
	}
	return x
}

// Return a new slice which is expanded by adding a 0 on the 0th index
func PadLeftWith0(xs []float64) []float64 {
	r := make([]float64, len(xs)+1)
	for i, _ := range r {
		if i == 0 {
			r[i] = 0
		} else {
			r[i] = xs[i-1]
		}
	}
	return r
}

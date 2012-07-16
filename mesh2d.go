package polecalc

import "math"

// Return coordinate from the square mesh of lenght L corresponding to the 
// index i.
// i=0 corresponds to (-pi, -pi); i=L-1 is (pi-step, -pi);
// i=L is (-pi, -pi+step); i=L^2-1 is (pi-step, pi-step)
func SquareAt(i uint64, L uint32) Vector2 {
	if i < 0 || i >= uint64(math.Pow(float64(L), 2.0)) {
		// panic here instead of returning an error since we will call
		// this function pretty often - presumably only returning one
		// variable is better for performance
		panic("invalid index for square mesh")
	}
	start, stop := -math.Pi, math.Pi
	length := stop - start
	step := length / float64(L)
	// transform 1d index to 2d coordinate indices
	ny := uint32(math.Floor(float64(i) / float64(L)))
	nx := uint32(i) - ny*L
	// get coordinates
	x := start + float64(nx)*step
	y := start + float64(ny)*step
	return Vector2{x, y}
}

type Callback func(k Vector2) error
type Acceptor func(k Vector2) bool

func CallOnAccepted(L uint32, callback Callback, acceptor Acceptor) error {
	sqrtN := uint64(L)
	N := sqrtN * sqrtN
	for i := uint64(0); i < N; i++ {
		k := SquareAt(i, L)
		if acceptor(k) {
			err := callback(k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// call callback on all points in the square mesh of length L
func CallOnPlane(L uint32, callback Callback) error {
	acceptor := func(k Vector2) bool {
		return true
	}
	return CallOnAccepted(L, callback, acceptor)
}

// call callback on the third quadrant only
func CallOnThirdQuad(L uint32, callback Callback) error {
	acceptor := func(k Vector2) bool {
		return k.X <= 0 && k.Y <= 0
	}
	return CallOnAccepted(L, callback, acceptor)
}

// call callback on the given curve
type CurveGenerator func(float64) Vector2

func CallOnCurve(curve CurveGenerator, numPoints uint, callback Callback) error {
	xValues := MakeRange(0.0, 1.0, numPoints)
	for _, x := range xValues {
		k := curve(x)
		err := callback(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// Scan k values along lines of high symmetry. Call callback at each point.
// k values form a cycle: (0, 0) -> (pi, 0) -> (pi, pi) -> (0, 0)
func CallOnSymmetryLines(numPoints uint, callback Callback) error {
	xAxis := func(x float64) Vector2 {
		return Vector2{x * math.Pi, 0.0}
	}
	yAxisShifted := func(y float64) Vector2 {
		return Vector2{math.Pi, y * math.Pi}
	}
	equal_xy := func(x float64) Vector2 {
		return Vector2{math.Pi - x*math.Pi, math.Pi - x*math.Pi}
	}
	curves := []CurveGenerator{xAxis, yAxisShifted, equal_xy}
	for _, curve := range curves {
		err := CallOnCurve(curve, numPoints, callback)
		if err != nil {
			return err
		}
	}
	return nil
}

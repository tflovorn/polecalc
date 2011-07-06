package polecalc

import (
	"math"
	"testing"
)

// Does SquareAt produce the expected points?
func TestSquareAtKnown(t *testing.T) {
	L, L64 := uint32(32), uint64(32) // points per side
	step := 2 * math.Pi / float64(L)
	points := []Vector2{SquareAt(0, L), SquareAt(L64-1, L), SquareAt(L64, L), SquareAt(L64*L64-1,L)}
	expected := []Vector2{Vector2{-math.Pi, -math.Pi}, Vector2{math.Pi-step, -math.Pi}, Vector2{-math.Pi,-math.Pi+step}, Vector2{math.Pi-step, math.Pi-step}}
	for i, p := range points {
		e := expected[i]
		if !FuzzyEqual(p.X, e.X) || !FuzzyEqual(p.Y, e.Y) {
			t.Fatalf("square point %d is not as expected (got %v, expected %v", i, p, e)
		}
	}
}

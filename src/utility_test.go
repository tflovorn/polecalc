package polecalc

import (
	"testing"
	"math"
)

func TestMakeRange(t *testing.T) {
	rng := MakeRange(0.0, 10.0, 11)
	for i := 0; i < 11; i++ {
		if math.Fabs(rng[i]-float64(i)) > MachEpsFloat64() {
			t.Fatal("incorrect range produced, got %v", rng)
		}
	}
}

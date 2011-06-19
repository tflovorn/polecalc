package polecalc

import "math"

func ZeroTempDelta(env Environment, k []float64) float64 {
	sx := math.Sin(k[0])
	sy := math.Sin(k[1])
	return 4 * env.F0 * (env.T0 + env.Tz) * (sx + float64(env.Alpha) * sy)
}

func ZeroTempPairEnergy(env Environment, k []float64) float64 {
	xi := Xi(env, k)
	delta := ZeroTempDelta(env, k)
	return math.Sqrt(xi * xi + delta * delta)
}

func ZeroTempFermi(energy float64) float64 {
	if (energy <= 0.0) {
		return 1.0
	}
	return 0.0
}

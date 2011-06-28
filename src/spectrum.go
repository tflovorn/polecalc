package polecalc

import "math"

// Single-hole hopping energy.  Minimum must be 0.
func Epsilon(env Environment, k Vector2) float64 {
	return EpsilonBar(env, k) - env.EpsilonMin
}

// Single-hole hopping energy without fixed minimum.
func EpsilonBar(env Environment, k Vector2) float64 {
	sx, sy := math.Sin(k.X), math.Sin(k.Y)
	return 2*env.Th()*((sx+sy)*(sx+sy)-1) + 4*(env.D1*env.T0-env.Thp)*sx*sy
}

// Find the minimum of EpsilonBar() to help in calculating Epsilon()
func EpsilonMin(env Environment) float64 {
	worker := func(k Vector2) float64 {
		return EpsilonBar(env, k)
	}
	return Minimum(env.GridLength, worker, env.NumProcs)
}

// Effective hopping energy (epsilon - mu).  Minimum is -mu.
func Xi(env Environment, k Vector2) float64 {
	return Epsilon(env, k) - env.Mu
}

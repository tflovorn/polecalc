package polecalc

import "math"

func EpsilonMin(env Environment) float64 {
	worker := func(k []float64) float64 {
		return EpsilonBar(env, k)
	}
	return Minimum(env.GridLength, worker, env.NumProcs)
}

func Epsilon(env Environment, k []float64) float64 {
	return EpsilonBar(env, k) - env.EpsilonMin
}

func EpsilonBar(env Environment, k []float64) float64 {
	sx, sy := math.Sin(k[0]), math.Sin(k[1])
	return 2*env.Th()*((sx+sy)*(sx+sy)-1) + 4*(env.D1*env.T0-env.Thp)*sx*sy
}

func Xi(env Environment, k []float64) float64 {
	return Epsilon(env, k) - env.Mu
}

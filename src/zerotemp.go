package polecalc

import (
	"math"
	"os"
)

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

type ZeroTempD1Equation struct{}

func (eq ZeroTempD1Equation) AbsError(args interface{}) float64 {
	env := args.(Environment)
	worker := func(k []float64) float64 {
		sx := math.Sin(k[0])
		sy := math.Sin(k[1])
		return -0.5 * (1 - Xi(env, k)) / ZeroTempPairEnergy(env, k) * sx * sy
	}
	return env.D1 - Average(env.GridLength, worker, env.NumProcs)
}

func (eq ZeroTempD1Equation) SetArguments(D1 float64, args *interface{}) {
	env := (*args).(Environment)
	env.D1 = D1
}

func (eq ZeroTempD1Equation) Range(args interface{}) (float64, float64, os.Error) {
	return 0.0, 1.0, nil
}

type ZeroTempMuEquation struct{}

func (eq ZeroTempMuEquation) AbsError(args interface{}) float64 {
	env := args.(Environment)
	worker := func(k []float64) float64 {
		return 0.5 * (1 - Xi(env, k)) / ZeroTempPairEnergy(env, k)
	}
	return env.X - Average(env.GridLength, worker, env.NumProcs)
}

func (eq ZeroTempMuEquation) SetArguments(Mu float64, args *interface{}) {
	env := (*args).(Environment)
	env.Mu = Mu
}

func (eq ZeroTempMuEquation) Range(args interface{}) (float64, float64, os.Error) {
	env := args.(Environment)
	return -2 * env.T0, 2 * env.T0, nil
}

type ZeroTempF0Equation struct{}

func (eq ZeroTempF0Equation) AbsError(args interface{}) float64 {
	env := args.(Environment)
	worker := func(k []float64) float64 {
		sinPart := math.Sin(k[0]) + float64(env.Alpha) * math.Sin(k[1])
		return sinPart * sinPart / ZeroTempPairEnergy(env, k)
	}
	return env.D1 - Average(env.GridLength, worker, env.NumProcs)
}

func (eq ZeroTempF0Equation) SetArguments(F0 float64, args *interface{}) {
	env := (*args).(Environment)
	env.F0 = F0
}

func (eq ZeroTempF0Equation) Range(args interface{}) (float64, float64, os.Error) {
	return 0.0, 1.0, nil
}

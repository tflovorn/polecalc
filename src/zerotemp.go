package polecalc

import (
	"math"
	"os"
)

// Returns the system of equations needed to solve the system at T = 0
func NewZeroTempSystem(tolerances []float64) *SelfConsistentSystem {
	eqD1 := ZeroTempD1Equation{}
	eqMu := ZeroTempMuEquation{}
	eqF0 := ZeroTempF0Equation{}
	equations := []SelfConsistentEquation{eqD1, eqMu, eqF0}
	system := &SelfConsistentSystem{equations, tolerances}
	return system
}

// --- D1 equation ---

// D1 = -1/(2N) \sum_k (1 - xi(k)/E(k)) * sin(kx) * sin(ky)
func ZeroTempD1AbsError(env Environment) float64 {
	worker := func(k []float64) float64 {
		sx, sy := math.Sin(k[0]), math.Sin(k[1])
		return -0.5 * (1 - Xi(env, k)/ZeroTempPairEnergy(env, k)) * sx * sy
	}
	return env.D1 - Average(env.GridLength, worker, env.NumProcs)
}

type ZeroTempD1Equation struct{}

func (eq ZeroTempD1Equation) AbsError(args interface{}) float64 {
	return ZeroTempD1AbsError(args.(Environment))
}

func (eq ZeroTempD1Equation) SetArguments(D1 float64, args interface{}) interface{} {
	env := args.(Environment)
	env.D1 = D1
	// Epsilon depends on D1 so we may have changed the minimum
	env.EpsilonMin = EpsilonMin(env)
	return env
}

func (eq ZeroTempD1Equation) Range(args interface{}) (float64, float64, os.Error) {
	return 0.0, 1.0, nil
}

// --- mu equation ---

// x = 1/(2N) \sum_k (1 - xi(k)/E(k))
func ZeroTempMuAbsError(env Environment) float64 {
	worker := func(k []float64) float64 {
		return 0.5 * (1 - Xi(env, k)/ZeroTempPairEnergy(env, k))
	}
	return env.X - Average(env.GridLength, worker, env.NumProcs)
}

type ZeroTempMuEquation struct{}

func (eq ZeroTempMuEquation) AbsError(args interface{}) float64 {
	return ZeroTempMuAbsError(args.(Environment))
}

func (eq ZeroTempMuEquation) SetArguments(Mu float64, args interface{}) interface{} {
	env := args.(Environment)
	env.Mu = Mu
	return env
}

// mu < 0 is enforced since for mu >= 0 terms with 1 / PairEnergy() can blow up
// Factor of -2 is arbitrary, may need to be enlarged for some Environments
func (eq ZeroTempMuEquation) Range(args interface{}) (float64, float64, os.Error) {
	env := args.(Environment)
	return -2 * env.T0, -MachEpsFloat64(), nil
}

// --- F0 equation ---

// 1/(t0+tz) = 1/N \sum_k (sin(kx) + alpha*sin(ky))^2 / E(k)
func ZeroTempF0AbsError(env Environment) float64 {
	worker := func(k []float64) float64 {
		sinPart := math.Sin(k[0]) + float64(env.Alpha)*math.Sin(k[1])
		return sinPart * sinPart / ZeroTempPairEnergy(env, k)
	}
	return 1/(env.T0+env.Tz) - Average(env.GridLength, worker, env.NumProcs)
}

type ZeroTempF0Equation struct{}

func (eq ZeroTempF0Equation) AbsError(args interface{}) float64 {
	return ZeroTempF0AbsError(args.(Environment))
}

func (eq ZeroTempF0Equation) SetArguments(F0 float64, args interface{}) interface{} {
	env := args.(Environment)
	env.F0 = F0
	return env
}

func (eq ZeroTempF0Equation) Range(args interface{}) (float64, float64, os.Error) {
	return 0.0, 1.0, nil
}

// --- energy scales and related functions ---

// Holon (pair?) gap energy.
func ZeroTempDelta(env Environment, k []float64) float64 {
	sx, sy := math.Sin(k[0]), math.Sin(k[1])
	return 4 * env.F0 * (env.T0 + env.Tz) * (sx + float64(env.Alpha)*sy)
}

// Energy of a pair of holes.
func ZeroTempPairEnergy(env Environment, k []float64) float64 {
	xi := Xi(env, k)
	delta := ZeroTempDelta(env, k)
	return math.Sqrt(xi*xi + delta*delta)
}

// Energy of a singlet (?)
func ZeroTempOmega(env Environment, k []float64) float64 {
	phi := func(k []float64) float64 {
		// where are J and A specified?
		return 4 * env.J * env.A * (math.Sin(k[0]) + math.Sin(k[1]))
	}(k)
	return math.Sqrt(env.Lambda*env.Lambda - phi*phi)
}

// Fermi distribution at T = 0 is H(-x), where H is the Heaviside step function.
// H(0) is taken to be 1.
func ZeroTempFermi(energy float64) float64 {
	if energy <= 0.0 {
		return 1.0
	}
	return 0.0
}

// --- Green's function for the physical electron ---

// imaginary part - values for all omega are calculated simultaneously, so
// return two slices of floats.  first is omega values, second is coefficients
func ZeroTempImG0(env Environment, k []float64) ([]float64, []float64) {
	maxWorker := func(k []float64) float64 {
		return ZeroTempPairEnergy(env, k)
	}
	pairEnergyMax := Maximum(env.GridLength, maxWorker, env.NumProcs)
	maxAbsOmega := env.Lambda + pairEnergyMax
	omegaMin, omegaMax := -maxAbsOmega, maxAbsOmega
	plusMinus := func(x, y float64) (float64, float64) {
		return x + y, x - y
	}
	deltaTerms := func(q []float64) ([]float64, []float64) {
		omega_q := ZeroTempOmega(env, q)
		E_h := ZeroTempPairEnergy(env, q)
		omegas := []float64{omega_q - E_h, omega_q + E_h, -omega_q - E_h, -omega_q + E_h}
		lambda_p, lambda_m := plusMinus(1, env.Lambda/ZeroTempOmega(env, q))
		xi_p, xi_m := plusMinus(1, Xi(env, q)/ZeroTempPairEnergy(env, q))
		// 0 here is really bose function of omega_q = ZeroTempOmega(q)
		// since omega_q > 0 and mu < 0, bose function result is 0
		f_p, f_m := plusMinus(0, ZeroTempFermi(ZeroTempPairEnergy(env, []float64{q[0] - k[0], q[1] - k[1]})))
		coeffs := []float64{0.25 * lambda_p * xi_p * f_p, 0.25 * lambda_p * xi_m * (f_m + 1), -0.25 * lambda_m * xi_m * (f_m + 1), -0.25 * lambda_m * xi_m * f_p}
		return omegas, coeffs
	}
	binner := NewDeltaBinner(deltaTerms, omegaMin, omegaMax, env.ImG0Bins)
	result := DeltaBin(env.GridLength, binner, env.NumProcs)
	omegas := binner.BinVarValues()
	return omegas, result
}

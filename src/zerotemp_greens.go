package polecalc

import (
	"fmt"
	"math"
	"os"
	"reflect"
)

// --- noninteracting Green's function for the physical electron ---

// -- imaginary part of noninteracting Green's function --
func deltaTermsGc0(env Environment, k Vector2, q Vector2) ([]float64, []float64) {
	omega_q := ZeroTempOmega(env, q)
	E_h := ZeroTempPairEnergy(env, q)
	lambda_p, lambda_m := plusMinus(1, env.Lambda()/ZeroTempOmega(env, q))
	// 0 in f_p & f_m is really bose function of omega_q = ZeroTempOmega(q)
	// since omega_q > 0 and mu < 0, bose function result is 0
	f_p, f_m := plusMinus(0, ZeroTempFermi(ZeroTempPairEnergy(env, q.Sub(k))))
	c := -0.25 * math.Pi
	if env.Superconducting {
		xi_p, xi_m := plusMinus(1, Xi(env, q)/ZeroTempPairEnergy(env, q))
		omegas := []float64{omega_q - E_h, omega_q + E_h, -omega_q - E_h, -omega_q + E_h}

		coeffs := []float64{c * lambda_p * xi_p * f_p, c * lambda_p * xi_m * (f_m + 1), -c * lambda_m * xi_m * (f_m + 1), -c * lambda_m * xi_m * f_p}
		return omegas, coeffs
	}
	omegas := []float64{omega_q - E_h, -omega_q - E_h}
	coeffs := []float64{c * lambda_p * f_p, -c * lambda_m * (f_m + 1)}
	return omegas, coeffs
}

// values for all omega are calculated simultaneously, so return two slices of 
// floats.  first is omega values, second is coefficients
func ZeroTempImGc0(env Environment, k Vector2) ([]float64, []float64) {
	var omegaMin, omegaMax float64
	if env.Superconducting {
		pairWorker := func(k Vector2) float64 {
			return ZeroTempPairEnergy(env, k)
		}
		pairEnergyMax := Maximum(env.GridLength, pairWorker)
		maxAbsOmega := env.Lambda() + pairEnergyMax
		omegaMin, omegaMax = -maxAbsOmega, maxAbsOmega
	} else {
		xiWorker := func(k Vector2) float64 {
			return Xi(env, k)
		}
		xiMax := Maximum(env.GridLength, xiWorker)
		maxAbsOmega := env.Lambda() + xiMax
		omegaMin, omegaMax = -maxAbsOmega, maxAbsOmega
	}
	deltaTerms := func(q Vector2) ([]float64, []float64) {
		return deltaTermsGc0(env, k, q)
	}
	binner := NewDeltaBinner(deltaTerms, omegaMin, omegaMax, env.ImGc0Bins)
	result := DeltaBin(env.GridLength, binner)
	omegas := binner.BinVarValues()
	return omegas, result
}

func cachedEnvIndex(env Environment) int {
	for i, cacheEnv := range imGc0CacheEnv {
		if reflect.DeepEqual(env, cacheEnv) {
			return i
		}
	}
	return -1
}

func cachedImGc0(env Environment, k Vector2) (*CubicSpline, bool) {
	i := cachedEnvIndex(env)
	if i == -1 {
		return nil, false
	}
	cacheX, ok := imGc0CacheK[i][k.X]
	if !ok {
		return nil, false
	}
	spline, ok := cacheX[k.Y]
	return spline, ok
}

func addToCacheImGc0(env Environment, k Vector2, spl *CubicSpline) {
	i := cachedEnvIndex(env)
	if i == -1 {
		// env not encountered yet
		imGc0CacheEnv = append(imGc0CacheEnv, env)
		i = len(imGc0CacheEnv) - 1
		imGc0CacheK[i] = make(map[float64]map[float64]*CubicSpline)
	}
	if xCache, ok := imGc0CacheK[i][k.X]; ok {
		// have seen this X before
		xCache[k.Y] = spl
	} else {
		// new X
		imGc0CacheK[i][k.X] = make(map[float64]*CubicSpline)
		imGc0CacheK[i][k.X][k.Y] = spl
	}
}

func getFromCacheImGc0(env Environment, k Vector2) (*CubicSpline, os.Error) {
	var imPart *CubicSpline
	if cache, ok := cachedImGc0(env, k); ok {
		imPart = cache
	} else {
		var err os.Error
		imPartOmegaVals, imPartFuncVals := ZeroTempImGc0(env, k)
		imPart, err = NewCubicSpline(imPartOmegaVals, imPartFuncVals)
		if err != nil {
			return nil, err
		}
		addToCacheImGc0(env, k, imPart)
	}
	return imPart, nil
}

// implementing this the lazy way for now by interpolating ImGc0(k)
// could also calculate ImGc0(k,omega) directly
func ZeroTempImGc0Point(env Environment, k Vector2, omega float64) (float64, os.Error) {
	imPart, err := getFromCacheImGc0(env, k)
	if err != nil {
		return 0.0, err
	}
	omegaMin, omegaMax := imPart.Range()
	if omegaMin <= omega && omega <= omegaMax {
		return imPart.At(omega), nil
	}
	return 0.0, nil
}

// -- real part of noninteracting Green's function --
func ZeroTempReGc0(env Environment, k Vector2, omega float64) (float64, os.Error) {
	imPart, err := getFromCacheImGc0(env, k)
	if err != nil {
		return 0.0, err
	}
	omegaMin, omegaMax := imPart.Range()
	// assume that Im(Gc0) is smooth near omegaPrime, so that spline
	// interpolation is good enough
	integrand := func(omegaPrime float64) float64 {
		if omegaMin <= omega && omega <= omegaMax {
			return (1 / math.Pi) * imPart.At(omegaPrime) / (omegaPrime - omega)
		}
		return 0.0
	}
	integral, err := PvIntegral(integrand, omegaMin, omegaMax, omega, env.ReGc0dw, env.ReGc0Points)
	if err != nil {
		return 0.0, err
	}
	return integral, nil
}

// --- full Green's function poles ---
// find solutions to Re[1/Gc0(k,omega)] - ElectronEnergy(k) = 0
// ==> ((ReGc0)^2 + (ImGc0)^2)*ElectronEnergy - ReGc0 = 0

type ZeroTempGreenPoleEq struct {
	K Vector2
}

type ZeroTempGreenArgs struct {
	Env   Environment
	Omega float64
}

func (eq ZeroTempGreenPoleEq) AbsError(args interface{}) float64 {
	greenArgs := args.(ZeroTempGreenArgs)
	env, omega := greenArgs.Env, greenArgs.Omega
	ReGc0, err := ZeroTempReGc0(env, eq.K, omega)
	if err != nil {
		panic("error encountered searching for ReGc0: " + err.String())
	}
	/*
		ImGc0, err := ZeroTempImGc0Point(env, eq.K, omega)
		if err != nil {
			panic("error encountered searching for ImGc0: " + err.String())
		}
	*/
	epsilon_k := ZeroTempElectronEnergy(env, eq.K)
	//return (ReGc0*ReGc0+ImGc0*ImGc0)*epsilon_k - ReGc0
	return 1.0 - epsilon_k*ReGc0
}

func (eq ZeroTempGreenPoleEq) SetArguments(omega float64, args interface{}) interface{} {
	env := args.(ZeroTempGreenArgs).Env
	return ZeroTempGreenArgs{env, omega}
}

func (eq ZeroTempGreenPoleEq) Range(args interface{}) (float64, float64, os.Error) {
	env := args.(ZeroTempGreenArgs).Env
	return -10.0 * env.T, 10.0 * env.T, nil
}

// return all pole omegas at a given k
func ZeroTempGreenPolePoint(env Environment, k Vector2) ([]float64, os.Error) {
	// find brackets for all the poles
	// lazy for now - only look at one solution
	eq := ZeroTempGreenPoleEq{k}
	initArgs := ZeroTempGreenArgs{env, 0.0}
	solvedArgs, err := MultiSolve(eq, initArgs)
	if err != nil {
		return nil, err
	}
	solutions := []float64{}
	for _, args := range solvedArgs {
		omega := args.(ZeroTempGreenArgs).Omega
		solutions = append(solutions, omega)
	}
	return solutions, nil
}

// real part of the full Green's function
func FullReGc(env Environment, k Vector2, omega float64) (float64, os.Error) {
	ReGc0, err := ZeroTempReGc0(env, k, omega)
	if err != nil {
		return 0.0, err
	}
	ImGc0, err := ZeroTempImGc0Point(env, k, omega)
	if err != nil {
		return 0.0, err
	}
	mag := ReGc0*ReGc0 + ImGc0*ImGc0
	epsilon_k := ZeroTempElectronEnergy(env, k)
	numer := mag * (ReGc0 - mag*epsilon_k)
	denom := math.Pow(ReGc0-mag*epsilon_k, 2.0) + ImGc0*ImGc0
	return numer / denom, nil
}

// --- plotting helper functions ---

type GreenPole struct {
	K     Vector2
	Omega float64
}

func (gp GreenPole) String() string {
	return fmt.Sprintf("k: %v; omega: %f", gp.K, gp.Omega)
}

func capturePoles(env Environment, k Vector2, poles []GreenPole) ([]GreenPole, os.Error) {
	kPoles, err := ZeroTempGreenPolePoint(env, k)
	if err != nil {
		if err.String() == ErrorNoBracket {
			println("bracket error at k = ", k.String())
			return poles, nil
		}
		return poles, err
	} else {
		for _, p := range kPoles {
			println("got pole k=", k.String(), p)
			poles = append(poles, GreenPole{k, p})
		}
	}
	return poles, err
}

// scan the k space looking for poles; return all those found
func ZeroTempGreenPolePlane(env Environment, pointsPerSide uint32, minimal bool) ([]GreenPole, os.Error) {
	poles := []GreenPole{}
	callback := func(k Vector2) os.Error {
		var err os.Error
		poles, err = capturePoles(env, k, poles)
		return err
	}
	err := CallOnThirdQuad(pointsPerSide, callback)
	return poles, err
}

// Scan k values given along poleCurve, which takes a value from 0 to 1 and 
// returns a vector in k space.  Return all poles found.
func ZeroTempGreenPoleCurve(env Environment, poleCurve CurveGenerator, numPoints uint) ([]GreenPole, os.Error) {
	poles := []GreenPole{}
	callback := func(k Vector2) os.Error {
		var err os.Error
		poles, err = capturePoles(env, k, poles)
		return err
	}
	err := CallOnCurve(poleCurve, numPoints, callback)
	return poles, err
}

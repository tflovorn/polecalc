package polecalc

import (
	"fmt"
	"math"
)

var imGc0Cache = NewListCache()

// --- noninteracting Green's function for the physical electron ---

// -- imaginary part of noninteracting Green's function --
func deltaTermsGc0(env Environment, k Vector2, q Vector2) ([]float64, []float64) {
	omega_q := ZeroTempOmega(env, q)
	E_h := ZeroTempPairEnergy(env, q.Sub(k))
	lambda_p, lambda_m := plusMinus(1, env.Lambda()/omega_q)
	// 0 in f_p & f_m is really bose function of omega_q = ZeroTempOmega(q)
	// since omega_q > 0 and mu < 0, bose function result is 0
	f_p, f_m := plusMinus(0, ZeroTempFermi(E_h))
	if env.Superconducting {
		c := -0.25 * math.Pi
		xi := Xi(env, q.Sub(k))
		xi_p, xi_m := plusMinus(1, xi/E_h)
		omegas := []float64{omega_q - E_h, omega_q + E_h, -omega_q - E_h, -omega_q + E_h}

		coeffs := []float64{c * lambda_p * xi_p * f_p, c * lambda_p * xi_m * (f_m + 1), -c * lambda_m * xi_p * (f_m + 1), -c * lambda_m * xi_m * f_p}
		return omegas, coeffs
	}
	// non-superconducting if we get here
	c := -0.5 * math.Pi
	omegas := []float64{omega_q - E_h, -omega_q - E_h}
	coeffs := []float64{c * lambda_p * f_p, -c * lambda_m * (f_m + 1)}
	return omegas, coeffs
}

// values for all omega are calculated simultaneously, so return two slices of 
// floats.  first is omega values, second is coefficients
func ZeroTempImGc0(env Environment, k Vector2) ([]float64, []float64) {
	var omegaMin, omegaMax float64
	if env.Superconducting {
		pairWorker := func(q Vector2) float64 {
			return ZeroTempPairEnergy(env, q.Sub(k))
		}
		pairEnergyMax := Maximum(env.GridLength, pairWorker)
		maxAbsOmega := env.Lambda() + pairEnergyMax
		omegaMin, omegaMax = -maxAbsOmega-1.0, maxAbsOmega+1.0
	} else {
		xiWorker := func(q Vector2) float64 {
			return Xi(env, q.Sub(k))
		}
		xiMax := Maximum(env.GridLength, xiWorker)
		maxAbsOmega := env.Lambda() + xiMax
		omegaMin, omegaMax = -maxAbsOmega-1.0, maxAbsOmega+1.0
	}
	deltaTerms := func(q Vector2) ([]float64, []float64) {
		return deltaTermsGc0(env, k, q)
	}
	binner := NewDeltaBinner(deltaTerms, omegaMin, omegaMax, env.ImGc0Bins)
	result := DeltaBin(env.GridLength, binner)
	omegas := binner.BinVarValues()
	return omegas, result
}

func cachedImGc0(env Environment, k Vector2) (*CubicSpline, bool) {
	kCacheInterface, ok := imGc0Cache.Get(env)
	if !ok {
		return nil, false
	}
	kCache := kCacheInterface.(VectorCache)
	spline, ok := kCache.Get(k)
	if ok {
		return spline.(*CubicSpline), ok
	}
	return nil, ok
}

func addToCacheImGc0(env Environment, k Vector2, spl *CubicSpline) {
	if !imGc0Cache.Contains(env) {
		// env not encountered yet
		kCache := *NewVectorCache()
		kCache.Set(k, spl)
		imGc0Cache.Set(env, kCache)
	} else {
		kCacheInterface, _ := imGc0Cache.Get(env)
		kCache := kCacheInterface.(VectorCache)
		kCache.Set(k, spl)
	}
}

func getFromCacheImGc0(env Environment, k Vector2) (*CubicSpline, error) {
	var imPart *CubicSpline
	if cache, ok := cachedImGc0(env, k); ok {
		imPart = cache
	} else {
		var err error
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
func ZeroTempImGc0Point(env Environment, k Vector2, omega float64) (float64, error) {
	imPart, err := getFromCacheImGc0(env, k)
	if err != nil {
		return 0.0, err
	}
	omegaMin, omegaMax := imPart.Range()
	if omegaMin <= omega && omega <= omegaMax {
		return imPart.At(omega)
	}
	return 0.0, nil
}

// -- real part of noninteracting Green's function --
func ZeroTempReGc0(env Environment, k Vector2, omega float64) (float64, error) {
	imPart, err := getFromCacheImGc0(env, k)
	if err != nil {
		return 0.0, err
	}
	omegaMin, omegaMax := imPart.Range()
	// assume that Im(Gc0) is smooth near omegaPrime, so that spline
	// interpolation is good enough
	integrand := func(omegaPrime float64) (float64, error) {
		if omegaMin <= omega && omega <= omegaMax {
			im, err := imPart.At(omegaPrime)
			if err != nil {
				return 0.0, err
			}
			return (1 / math.Pi) * im / (omegaPrime - omega), nil
		}
		return 0.0, nil
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
		panic("error encountered searching for ReGc0: " + err.Error())
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

func (eq ZeroTempGreenPoleEq) Range(args interface{}) (float64, float64, error) {
	env := args.(ZeroTempGreenArgs).Env
	return -10.0 * env.T, 10.0 * env.T, nil
}

// return all pole omegas at a given k
func ZeroTempGreenPolePoint(env Environment, k Vector2) ([]float64, error) {
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
func FullReGc(env Environment, k Vector2, omega float64) (float64, error) {
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

func capturePoles(env Environment, k Vector2, poles []GreenPole) ([]GreenPole, error) {
	kPoles, err := ZeroTempGreenPolePoint(env, k)
	if err != nil {
		if err.Error() == ErrorNoBracket {
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
func ZeroTempGreenPolePlane(env Environment, pointsPerSide uint32, minimal bool) ([]GreenPole, error) {
	poles := []GreenPole{}
	callback := func(k Vector2) error {
		var err error
		poles, err = capturePoles(env, k, poles)
		return err
	}
	err := CallOnThirdQuad(pointsPerSide, callback)
	return poles, err
}

// Scan k values given along poleCurve, which takes a value from 0 to 1 and 
// returns a vector in k space.  Return all poles found.
func ZeroTempGreenPoleCurve(env Environment, poleCurve CurveGenerator, numPoints uint) ([]GreenPole, error) {
	poles := []GreenPole{}
	callback := func(k Vector2) error {
		var err error
		poles, err = capturePoles(env, k, poles)
		return err
	}
	err := CallOnCurve(poleCurve, numPoints, callback)
	return poles, err
}

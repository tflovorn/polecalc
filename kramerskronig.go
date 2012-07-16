package polecalc

import "math"

// Return the real part of f = realPart + i * imagPart.
// f must be analytic in the upper half-plane (i.e. causal).
// left and right are the limits of integration (approximation to -+ infinity).
func RealFromImaginary(imagPart Func1DError, left, right, eps float64, numPoints uint) Func1DError {
	realPart := func(omega float64) (float64, error) {
		integrand := func(omegaPrime float64) (float64, error) {
			im, err := imagPart(omegaPrime)
			if err != nil {
				return 0.0, err
			}
			return (1 / math.Pi) * im / (omegaPrime - omega), nil
		}
		return PvIntegral(integrand, left, right, omega, eps, numPoints)
	}
	return realPart
}

// Kramers-Kronig equation for the imaginary part is identical to that for the
// real part except for a minus sign.
func ImaginaryFromReal(realPart Func1DError, left, right, eps float64, numPoints uint) Func1DError {
	imaginaryPart := func(omega float64) (float64, error) {
		almost := RealFromImaginary(realPart, left, right, eps, numPoints)
		almostIm, err := almost(omega)
		return -almostIm, err
	}
	return imaginaryPart
}

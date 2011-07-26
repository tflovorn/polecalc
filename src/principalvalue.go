// Interface to GSL principal value integral functionality
// (see http://www.gnu.org/software/gsl/manual/html_node/QAWC-adaptive-integration-for-Cauchy-principal-values.html)
// Callback passing through cgo follows the model at:
// http://stackoverflow.com/questions/6125683/call-go-functions-from-c/6147097#6147097
package polecalc

// #cgo LDFLAGS: -lgsl
// #include <stdlib.h>
// #include <gsl/gsl_integration.h>
// extern double goEvaluate(double, void*);
//
// static double goPvIntegral(double a, double b, double c, double epsabs, double epsrel, size_t limit, void* userdata) {
// 	gsl_integration_workspace *w = gsl_integration_workspace_alloc(limit);
//	double result, error;
//	gsl_function F;
//	F.function = &goEvaluate;
//	F.params = userdata;
//	gsl_integration_qawc(&F, a, b, c, epsabs, epsrel, limit, w, &result, &error);
//	gsl_integration_workspace_free(w);
//	return result;
// }
import "C"
import "unsafe"

//export goEvaluate
func goEvaluate(x C.double, userdata unsafe.Pointer) C.double {
	req := (*Func1D)(userdata)
	f := *req
	xfloat := float64(x)
	val := f(xfloat)
	return C.double(val)
}

func PvIntegralGSL(f Func1D, a, b, c, epsabs, epsrel float64, limit uint16) float64 {
	req := unsafe.Pointer(&f)
	return float64(C.goPvIntegral(C.double(a), C.double(b), C.double(c), C.double(epsabs), C.double(epsrel), C.size_t(limit), req))
}

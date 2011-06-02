package main

import "fmt"
import "./bisection"

func linear(x float64) float64 {
	return 1.0 - x
}

func quadratic(x float64) float64 {
	return x * x - 5.0
}

func main() {
	root, error := goroots.SolveBisection(linear, 0.0, 2.0, 1e-9)
	if error == nil {
		fmt.Println(root)
	} else {
		fmt.Println("error finding root in linear function: " + error.String())
	}
	root, error = goroots.SolveBisection(quadratic, 1.0, 3.0, 1e-9)
	if error == nil {
		fmt.Println(root)
	} else {
		fmt.Println("error finding root in quadratic function: " + error.String())
	}
}

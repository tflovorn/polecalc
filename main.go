package main

import "fmt"
import "./bisection"

func linear(x float64) float64 {
	return 1.0 - x
}

func main() {
	root, error := goroots.SolveBisection(linear, 0.0, 2.0, 1e-9)
	if error == nil {
		fmt.Println(root)
	} else {
		fmt.Println("error finding root: " + error.String())
	}
}

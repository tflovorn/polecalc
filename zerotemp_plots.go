package polecalc

import (
	"fmt"
	"math"
)

func ZeroTempPlotGc(env Environment, k Vector2, numOmega uint, outputPath string) error {
	imOmegas, imCalcValues := ZeroTempImGc0(env, k)
	imSpline, err := NewCubicSpline(imOmegas, imCalcValues)
	if err != nil {
		return err
	}
	imOmegaMin, imOmegaMax := imSpline.Range()
	omegas := MakeRange(imOmegaMin-1.0, imOmegaMax+1.0, numOmega)
	realValues := make([]float64, numOmega)
	imValues := make([]float64, numOmega)
	fullReValues := make([]float64, numOmega)
	for i := 0; i < int(numOmega); i++ {
		if omegas[i] < imOmegaMin || omegas[i] > imOmegaMax {
			imValues[i] = 0.0
		} else {
			im, err := imSpline.At(omegas[i])
			if err != nil {
				return err
			}
			imValues[i] = im
		}
		re, err := ZeroTempReGc0(env, k, omegas[i])
		if err != nil {
			return err
		}
		realValues[i] = re
		fullRe, err := FullReGc(env, k, omegas[i])
		if err != nil {
			return err
		}
		fullReValues[i] = fullRe
	}
	reGraph := NewGraph()
	imGraph := NewGraph()
	fullReGraph := NewGraph()
	rePath := outputPath + "_re"
	imPath := outputPath + "_im"
	fullRePath := outputPath + "_fullRe"
	reGraph.SetGraphParameters(map[string]interface{}{"graph_filepath": rePath})
	imGraph.SetGraphParameters(map[string]interface{}{"graph_filepath": imPath})
	fullReGraph.SetGraphParameters(map[string]interface{}{"graph_filepath": fullRePath})
	reData := make([][]float64, len(omegas))
	imData := make([][]float64, len(omegas))
	fullReData := make([][]float64, len(omegas))
	for i, _ := range reData {
		reData[i] = []float64{omegas[i], realValues[i]}
		imData[i] = []float64{omegas[i], imValues[i]}
		if !math.IsNaN(fullReValues[i]) {
			fullReData[i] = []float64{omegas[i], fullReValues[i]}
		} else {
			fullReData[i] = []float64{omegas[i], 0.0}
		}
	}
	reGraph.AddSeries(map[string]string{"label": "re_gc0"}, reData)
	imGraph.AddSeries(map[string]string{"label": "im_gc0"}, imData)
	fullReGraph.AddSeries(map[string]string{"label": "fullRe_gc0"}, fullReData)
	err = MakePlot(reGraph, rePath)
	if err != nil {
		return err
	}
	err = MakePlot(imGraph, imPath)
	if err != nil {
		return err
	}
	err = MakePlot(fullReGraph, fullRePath)
	if err != nil {
		return err
	}
	return nil
}

// Plot im/re gc0 and re gc along lines of high symmetry in k space.
func PlotGcSymmetryLines(env Environment, kPoints, numOmega uint, outputPath string) error {
	callback := func(k Vector2) error {
		extra := fmt.Sprintf("_kx_%f_ky_%f", k.X, k.Y)
		fullPath := outputPath + extra
		err := ZeroTempPlotGc(env, k, numOmega, fullPath)
		return err
	}
	err := CallOnSymmetryLines(kPoints, callback)
	return err
}

// Plot the existence of poles throughout the k plane.
func ZeroTempPlotPolePlane(env Environment, outputPath string, sideLength uint32) error {
	polePlane, err := ZeroTempGreenPolePlane(env, sideLength, true)
	if err != nil {
		return err
	}
	graphPoleData(polePlane, outputPath, &Vector2{32.0, 32.0})
	return nil
}

// Plot the line of poles specified by poleCurve, which takes a float value from
// 0 to 1 and returns a Vector2 corresponding to that value
func ZeroTempPlotPoleCurve(env Environment, poleCurve func(float64) Vector2, numPoints uint, outputPath string) error {
	polePoints, err := ZeroTempGreenPoleCurve(env, poleCurve, numPoints)
	if err != nil {
		return err
	}
	graphPoleData(polePoints, outputPath, nil)
	return nil
}

func graphPoleData(poles []GreenPole, outputPath string, dims *Vector2) {
	poleData := [][]float64{}
	for _, gp := range poles {
		k := gp.K
		poleData = append(poleData, []float64{k.X, k.Y})
	}
	poleGraph := NewGraph()
	params := make(map[string]interface{})
	if dims != nil {
		params["dimensions"] = []float64{dims.X, dims.Y}
	}
	params["graph_filepath"] = outputPath
	poleGraph.SetGraphParameters(params)
	poleGraph.AddSeries(map[string]string{"label": "poles", "style": "k."}, poleData)
	MakePlot(poleGraph, outputPath)
}

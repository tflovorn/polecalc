package polecalc

import "os"

func ZeroTempPlotGc0(env Environment, k Vector2, numOmega uint, outputPath string) os.Error {
	imOmegas, imCalcValues := ZeroTempImGc0(env, k)
	imSpline, err := NewCubicSpline(imOmegas, imCalcValues)
	if err != nil {
		return err
	}
	imOmegaMin, imOmegaMax := imSpline.Range()
	omegas := MakeRange(imOmegaMin-1.0, imOmegaMax+1.0, numOmega)
	realValues := make([]float64, numOmega)
	imValues := make([]float64, numOmega)
	for i := 0; i < int(numOmega); i++ {
		if omegas[i] < imOmegaMin || omegas[i] > imOmegaMax {
			imValues[i] = 0.0
		} else {
			imValues[i] = imSpline.At(omegas[i])
		}
		g, err := ZeroTempReGc0(env, k, omegas[i])
		if err != nil {
			return err
		}
		realValues[i] = g
	}
	reGraph := NewGraph()
	imGraph := NewGraph()
	rePath := outputPath + "_re"
	imPath := outputPath + "_im"
	reGraph.SetGraphParameters(map[string]interface{}{"graph_filepath": rePath})
	imGraph.SetGraphParameters(map[string]interface{}{"graph_filepath": imPath})
	reData := make([][]float64, len(omegas))
	imData := make([][]float64, len(omegas))
	for i, _ := range reData {
		reData[i] = []float64{omegas[i], realValues[i]}
		imData[i] = []float64{omegas[i], imValues[i]}
	}
	reGraph.AddSeries(map[string]string{"label": "re_gc0"}, reData)
	imGraph.AddSeries(map[string]string{"label": "im_gc0"}, imData)
	MakePlot(reGraph, rePath)
	MakePlot(imGraph, imPath)
	return nil
}

func ZeroTempPlotPolePlane(env Environment, outputPath string, sideLength uint32) os.Error {
	polePlane, err := ZeroTempGreenPolePlane(env, sideLength, true)
	if err != nil {
		return err
	}
	graphPoleData(polePlane, outputPath, &Vector2{32.0, 32.0})
	return nil
}

// Plot the line of poles specified by poleCurve, which takes a float value from
// 0 to 1 and returns a Vector2 corresponding to that value
func ZeroTempPlotPoleCurve(env Environment, poleCurve func(float64) Vector2, numPoints uint, outputPath string) os.Error {
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

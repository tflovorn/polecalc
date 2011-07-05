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
	reGraph.SetGraphParameters(map[string]string{"graph_filepath": rePath})
	imGraph.SetGraphParameters(map[string]string{"graph_filepath": imPath})
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

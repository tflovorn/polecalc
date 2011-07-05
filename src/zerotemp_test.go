package polecalc

import (
	"testing"
	"reflect"
	"flag"
	"math"
)

var cached *bool = flag.Bool("gc_cache", false, "used cached Environment for TestGc0")

func TestKnownZeroTempSystem(t *testing.T) {
	envStr := "{\"GridLength\":8,\"NumProcs\":1,\"InitD1\":0.1,\"InitMu\":0.1,\"InitF0\":0.1,\"Alpha\":-1,\"T0\":1,\"Tz\":0.1,\"Thp\":0.1,\"X\":0.1,\"D1\":0.05777149373506872,\"Mu\":-0.18330570279347036,\"F0\":0.12945949461029926,\"EpsilonMin\":-1.8}"
	expectedEnv, err := EnvironmentFromString(envStr)
	if err != nil {
		t.Fatal(err)
	}
	tolerances := []float64{1e-6, 1e-6, 1e-6}
	system := NewZeroTempSystem(tolerances)
	env, err := EnvironmentFromFile("zerotemp_test.json")
	if err != nil {
		t.Fatal(err)
	}
	env.Initialize()
	solution, err := system.Solve(*env)
	if err != nil {
		t.Fatal(err)
	}
	solvedEnv := solution.(Environment)
	if !reflect.DeepEqual(solvedEnv, *expectedEnv) {
		t.Fatalf("unknown solution to zero-temp system: got\n%s, expected\n%s", (&solvedEnv).String(), expectedEnv.String())
	}
}

func TestGc0(t *testing.T) {
	cacheFileName := "zerotemp_test_gc0_cache.json"
	flag.Parse()
	tolerances := []float64{1e-6, 1e-6, 1e-6}
	system := NewZeroTempSystem(tolerances)
	env, err := EnvironmentFromFile("zerotemp_test_gc0.json")
	if err != nil {
		t.Fatal(err)
	}
	env.Initialize()
	var solvedEnv Environment // not sure if the seperate declaration is needed
	if !(*cached) {
		solution, err := system.Solve(*env)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnv = solution.(Environment)
		solvedEnv.WriteToFile(cacheFileName)
	} else {
		cacheEnv, err := EnvironmentFromFile(cacheFileName)
		if err != nil {
			t.Fatal(err)
		}
		solvedEnv = *cacheEnv
	}
	numOmega := uint(2048)
	k := Vector2{math.Pi / 2.0, math.Pi / 2.0}
	k = Vector2{0.1, 0.1}
	imOmegas, imCalcValues := ZeroTempImGc0(solvedEnv, k)
	imSpline, err := NewCubicSpline(imOmegas, imCalcValues)
	if err != nil {
		t.Fatal(err)
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
		g, err := ZeroTempReGc0(solvedEnv, k, omegas[i])
		if err != nil {
			t.Fatal(err)
		}
		realValues[i] = g
	}
	reGraph := NewGraph()
	imGraph := NewGraph()
	rePath := "zerotemp_test_re_gc0.testignore"
	imPath := "zerotemp_test_im_gc0.testignore"
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
	poles, err := ZeroTempGreenPolePoint(solvedEnv, k)
	if err != nil {
		t.Fatal(err)
	}
	println(poles[0])
	println(imSpline.At(poles[0]))
	println(ZeroTempGap(solvedEnv, k))
}

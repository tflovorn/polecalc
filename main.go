package main

import (
	"fmt"
	"os"
	"./mpljson"
)

func main() {
	graph := mpljson.NewGraph()
	filepath, _ := os.Getwd()
	filepath = filepath + "/test_go_graph"
	graph.SetGraphParameters(map[string]string {"xlabel":"X_{test}$", "graph_filepath":filepath})
	graph.AddSeries(map[string]string {"label":"red", "style":"r-"}, [][]float64{{1.0, 2.0}, {3.0, 4.0}})
	err := mpljson.MakePlot(graph, "lol.json")
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

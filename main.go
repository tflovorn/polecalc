package main

import (
	"fmt"
	"./mpljson"
)

func main() {
	graph := mpljson.NewGraph()
	graph.SetGraphParameters(map[string]string {"xlabel":"X_{test}$", "graph_filepath":"lol"})
	graph.AddSeries(map[string]string {"label":"red", "style":"r-"}, [][]float64{{1.0, 2.0}, {3.0, 4.0}})
	err := mpljson.MakePlot(graph, "lol.json")
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

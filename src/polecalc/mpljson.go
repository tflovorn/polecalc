package polecalc

import (
	"json"
	"os"
	"exec"
)

const SERIES_KEY = "series"
const DATA_KEY = "data"

// A map compatible with the JSON Object type
type jsonObject map[string]interface{}

// Intermediate representation for graph data.
type Graph struct {
	graphParameters map[string]string
	seriesParameters []map[string]string
	seriesData [][][]float64
}

func NewGraph() *Graph {
	graph := new(Graph)
	graph.graphParameters = make(map[string]string)
	graph.seriesParameters = make([]map[string]string, 0)
	graph.seriesData = make([][][]float64, 0)
	return graph
}

// Set the global graph parameters as given in params.
func (graph *Graph) SetGraphParameters(params map[string]string) {
	for key, value := range params {
		graph.graphParameters[key] = value
	}
}

// Add a new data series to the graph.
func (graph *Graph) AddSeries(params map[string]string, data [][]float64) {
	graph.seriesParameters = append(graph.seriesParameters, params)
	graph.seriesData = append(graph.seriesData, data)
}

// Implements interface json.Marshaler
func (graph *Graph) MarshalJSON() ([]byte, os.Error) {
	jsonGraph := jsonObject {}
	// add global graph parameters
	for key, value := range graph.graphParameters {
		jsonGraph[key] = value
	}
	// add parameters and data for each series
	jsonGraph[SERIES_KEY] = []jsonObject {}
	for i, someSeriesParams := range graph.seriesParameters {
		newSeriesParams := jsonObject {}
		for key, value := range someSeriesParams {
			newSeriesParams[key] = value
		}
		newSeriesParams[DATA_KEY] = graph.seriesData[i]
		jsonGraph[SERIES_KEY] = append(jsonGraph[SERIES_KEY].([]jsonObject), newSeriesParams)
	}
	marshalled, err := json.Marshal(jsonGraph)
	return marshalled, err
}

// Constructs a plot from graph_data using matplotlib.
// graph_data must be a list or dictionary containing objects representable
// in JSON.  Blocks until Python script is finished.
func MakePlot(graphData interface{}, jsonFilePath string) os.Error {
	if err := WriteToJSONFile(graphData, jsonFilePath); err != nil {
		return err
	}
	wd, _ := os.Getwd()
	cmd, err := exec.Run("/usr/bin/env", []string{"/usr/bin/env", "python", wd + "/grapher.py", jsonFilePath}, os.Environ(), wd, exec.DevNull, exec.PassThrough, exec.PassThrough)
	if err != nil {
		return err
	}
	// Wait(0) -> wait for cmd to finish
	if _, err := cmd.Wait(0); err != nil {
		return err
	}
	if err := cmd.Close(); err != nil {
		return err
	}
	return nil
}

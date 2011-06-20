package polecalc

import (
	"json"
	"os"
	"math"
)

// Write Marshal-able object to a new file at filePath
func WriteToJSONFile(object interface{}, filePath string) os.Error {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return err
	}
	jsonFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err := jsonFile.Write(marshalled); err != nil {
		return err
	}
	if err := jsonFile.Close(); err != nil {
		return err
	}
	return nil
}

// Machine epsilon (upper bound on error due to rounding) for float64
func MachEpsFloat64() float64 {
	return math.Pow(2.0, -53.0)
}

package polecalc

import (
	"json"
	"os"
	"math"
	"strings"
)

// One-dimensional scalar function - used for bisection and bracket
type Func1D func(float64) float64

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

// Convert string to byte slice
func StringToBytes(str string) ([]byte, os.Error) {
	reader := strings.NewReader(str)
	bytes := make([]byte, len(str))
	for seen := 0; seen < len(str); {
		n, err := reader.Read(bytes)
		if err != nil {
			return nil, err
		}
		seen += n
	}
	return bytes, nil
}

package polecalc

import (
	"json"
	"os"
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

package prepareimport

import (
	"encoding/json"
	"os"
)

// SaveToFile saves the AdminObject as a JSON object at filePath
func SaveToFile(toSave AdminObject, filePath string) error {
	jsonData, err := json.Marshal(toSave)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

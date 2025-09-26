package prepareimport

import (
	"encoding/json"
	"opensignauxfaibles/lib/base"
	"os"
)

// SaveToFile saves the AdminObject as a JSON object at filePath
func SaveToFile(toSave base.AdminBatch, filePath string) error {
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

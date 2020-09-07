package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/globalsign/mgo/bson"
)

func parseJSONObject(filename string) (object bson.M, err error) {
	text, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(text, &object) // transform json string into bson.M
	}
	return object, err
}

func parseJSONArray(filename string) (array []bson.M, err error) {
	text, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(text, &array) // transform json string into bson.M
	}
	return array, err
}

// listJSONSchemaFiles retourne la liste des fichiers JSON Schema présents dans le répertoire validation.
func listJSONSchemaFiles() ([]string, error) {
	var files []string
	rootDir := "validation"
	err := filepath.Walk(rootDir, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && strings.Contains(filePath, ".schema.json") {
			files = append(files, filePath)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// LoadJSONSchemaFiles cherche les Schemas JSON pour GetRawDataValidationPipeline.
func LoadJSONSchemaFiles() (jsonSchema map[string]bson.M, err error) {
	jsonSchema = make(map[string]bson.M)
	files, err := listJSONSchemaFiles()
	if err != nil {
		return nil, err
	}
	for _, filePath := range files {
		dataType := strings.Replace(filepath.Base(filePath), ".schema.json", "", 1)
		jsonSchema[dataType], err = parseJSONObject(filePath)
		if err != nil {
			return nil, err
		}
	}
	return jsonSchema, nil
}

// GetRawDataValidationPipeline produit un pipeline pour retourner la listes des documents invalides depuis RawData.
func GetRawDataValidationPipeline(jsonSchema map[string]bson.M) (pipeline []bson.M, err error) {

	flattenPipeline, err := parseJSONArray("validation/flatten_RawData.pipeline.json")
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, flattenPipeline...)

	matchers := []bson.M{}
	for dataType, schema := range jsonSchema {
		matchers = append(matchers, bson.M{
			"dataType": dataType,
			"$nor": []bson.M{
				{
					"$jsonSchema": bson.M{
						"bsonType": "object",
						"properties": bson.M{
							"dataObject": schema,
						},
					},
				},
			},
		})
	}

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"$or": matchers,
		},
	})
	return pipeline, nil
}

package engine

import (
	"encoding/json"
	"io/ioutil"

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

// GetRawDataValidationPipeline produit un pipeline pour retourner la listes des documents invalides depuis RawData.
func GetRawDataValidationPipeline() (pipeline []bson.M, err error) {
	dataType := "delai"

	flattenPipeline, err := parseJSONArray("validation/flatten_RawData.pipeline.json")
	if err != nil {
		return nil, err
	}

	jsonSchema, err := parseJSONObject("validation/" + dataType + ".schema.json")
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, flattenPipeline...)

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"dataType": dataType,
			"$nor": []bson.M{
				{
					"$jsonSchema": bson.M{
						"bsonType": "object",
						"properties": bson.M{
							"dataObject": jsonSchema,
						},
					},
				},
			},
		},
	})
	return pipeline, nil
}

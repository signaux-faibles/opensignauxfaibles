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

func LoadJSONSchemaFiles() (jsonSchema map[string]bson.M, err error) {
	jsonSchema = make(map[string]bson.M)

	dataType := "delai"
	jsonSchema[dataType], err = parseJSONObject("validation/" + dataType + ".schema.json")
	if err != nil {
		return nil, err
	}

	jsonSchema["bdf"], err = parseJSONObject("validation/bdSf.schema.json")
	if err != nil {
		return nil, err
	}

	return jsonSchema, nil
}

// GetRawDataValidationPipeline produit un pipeline pour retourner la listes des documents invalides depuis RawData.
func GetRawDataValidationPipeline(jsonSchema map[string]bson.M) (pipeline []bson.M, err error) {
	dataType := "delai"

	flattenPipeline, err := parseJSONArray("validation/flatten_RawData.pipeline.json")
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, flattenPipeline...)

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"$or": []bson.M{
				{
					"dataType": dataType,
					"$nor": []bson.M{
						{
							"$jsonSchema": bson.M{
								"bsonType": "object",
								"properties": bson.M{
									"dataObject": jsonSchema[dataType],
								},
							},
						},
					},
				},
				{
					"dataType": "bdf",
					"$nor": []bson.M{
						{
							"$jsonSchema": bson.M{
								"bsonType": "object",
								"properties": bson.M{
									"dataObject": jsonSchema["bdf"],
								},
							},
						},
					},
				},
			},
		},
	})
	return pipeline, nil
}

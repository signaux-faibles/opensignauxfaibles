package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"
)

func TestDataValidation(t *testing.T) {

	schemaProps := loadPropsFromSchema("../../validation/delai.schema.json")
	structProps := reflectPropsFromStruct(urssaf.Delai{})

	errors := diffMaps(schemaProps, structProps)
	if len(errors) > 0 {
		log.Println("Types are not deeply equal:")
		for _, err := range errors {
			log.Println("- " + err.Error())
		}
		t.FailNow()
	}
}

func reflectPropsFromStruct(structInstance interface{}) map[string]propertySchema {
	props := make(map[string]propertySchema)
	fields := reflect.TypeOf(structInstance)
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName != "" {
			fieldType := field.Type.Name()
			// support pointer types
			if fieldType == "" {
				fieldType = field.Type.Elem().Name()
			}
			// convert go types to javascript equivalents
			if fieldType == "int" || fieldType == "float64" {
				fieldType = "number"
			} else if fieldType == "Time" {
				fieldType = "date"
			}
			props[fieldName] = propertySchema{fieldType}
		}
	}
	return props
}

func loadPropsFromSchema(filePath string) map[string]propertySchema {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	var schema jsonSchema
	json.Unmarshal(byteValue, &schema)
	return schema.Properties
}

func diffMaps(schemaProps map[string]propertySchema, structProps map[string]propertySchema) []error {
	errors := []error{}
	commonKeys := []string{}
	for k := range schemaProps {
		if _, ok := structProps[k]; !ok {
			errors = append(errors, fmt.Errorf("property not found in Go struct: %v", k))
		} else {
			commonKeys = append(commonKeys, k)
		}
	}
	for k := range structProps {
		if _, ok := schemaProps[k]; !ok {
			errors = append(errors, fmt.Errorf("property not found in JSON Schema: %v", k))
		}
	}
	for _, k := range commonKeys {
		if structProps[k] != schemaProps[k] {
			errors = append(errors, fmt.Errorf("property types of \"%v\" don't match: %v <> %v", k, schemaProps[k], structProps[k]))
		}
	}
	return errors
}

type jsonSchema struct {
	Properties map[string]propertySchema `json:"properties"`
}

type propertySchema struct {
	BsonType string `json:"bsonType"`
}

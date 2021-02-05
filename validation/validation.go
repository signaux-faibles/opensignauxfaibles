package validation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

type jsonSchema struct {
	Properties map[string]propertySchema `json:"properties"`
}

type propertySchema struct {
	BsonType        string                    `json:"bsonType,omitempty"`
	Properties      map[string]propertySchema `json:"properties,omitempty"`
	RequiredProps   []string                  `json:"required,omitempty"`
	AdditionalProps bool                      `json:"additionalProperties,omitempty"`
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
		if !reflect.DeepEqual(structProps[k], schemaProps[k]) {
			errors = append(errors, fmt.Errorf("property types of \"%v\" don't match: %v <> %v", k, schemaProps[k], structProps[k]))
		}
	}
	return errors
}

func reflectPropsFromStruct(structInstance interface{}) map[string]propertySchema {
	return reflectPropsFromType(reflect.TypeOf(structInstance))
}

func reflectPropsFromType(structType reflect.Type) map[string]propertySchema {
	props := make(map[string]propertySchema)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName != "" && fieldName != "-" {
			fieldName = strings.ReplaceAll(fieldName, ",omitempty", "")
			fieldType := field.Type.Name()
			if field.Type.Kind() == reflect.Struct && fieldType != "Time" {
				props[fieldName] = propertySchema{
					BsonType:        "object",
					Properties:      reflectPropsFromType(field.Type),
					RequiredProps:   []string{"start", "end"}, // TODO: get list of props extracted above
					AdditionalProps: false,
				}
			} else {
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
				props[fieldName] = propertySchema{BsonType: fieldType}
			}
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

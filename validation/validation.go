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

type propertySchema struct {
	BsonType        string                    `json:"bsonType,omitempty"`
	Properties      map[string]propertySchema `json:"properties,omitempty"`
	RequiredProps   []string                  `json:"required,omitempty"`
	AdditionalProps bool                      `json:"additionalProperties,omitempty"`
}

func diffSchema(jsonSchema propertySchema, structSchema propertySchema) []error {
	errors := diffMaps(jsonSchema.Properties, structSchema.Properties)
	jsonReqProps := map[string]interface{}{}
	for _, k := range jsonSchema.RequiredProps {
		jsonReqProps[k] = true
	}
	structReqProps := map[string]interface{}{}
	for _, k := range structSchema.RequiredProps {
		structReqProps[k] = true
	}
	for _, k := range jsonSchema.RequiredProps {
		if _, ok := structReqProps[k]; !ok {
			errors = append(errors, fmt.Errorf("required property is marked as 'omitempty' in Go struct: %v", k))
		}
	}
	for _, k := range structSchema.RequiredProps {
		if _, ok := jsonReqProps[k]; !ok {
			errors = append(errors, fmt.Errorf("property not marked as 'required' in JSON Schema: %v", k))
		}
	}
	return errors
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
	return reflectStructType(reflect.TypeOf(structInstance)).Properties
}

func reflectStructType(structType reflect.Type) propertySchema {
	requiredProps := []string{}
	props := make(map[string]propertySchema)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Tag.Get("bson")
		if fieldName != "" && fieldName != "-" {
			if strings.Contains(fieldName, ",omitempty") == true {
				fieldName = strings.ReplaceAll(fieldName, ",omitempty", "")
			} else {
				requiredProps = append(requiredProps, fieldName)
			}
			fieldType := field.Type.Name()
			if field.Type.Kind() == reflect.Struct && fieldType != "Time" {
				props[fieldName] = reflectStructType(field.Type)
			} else {
				// support pointer types
				if fieldType == "" {
					fieldType = field.Type.Elem().Name()
				}
				// convert go types to BSON equivalents (cf https://docs.mongodb.com/manual/reference/operator/query/type/#document-type-available-types)
				if fieldType == "int" || fieldType == "int64" {
					fieldType = "int64"
				} else if fieldType == "float64" {
					fieldType = "double"
				} else if fieldType == "Time" {
					fieldType = "date"
				} else if fieldType == "bool" {
					fieldType = "bool"
				} else if fieldType == "string" {
					fieldType = "string"
				} else {
					log.Fatal("Unsupported type: " + fieldType)
				}
				props[fieldName] = propertySchema{BsonType: fieldType}
			}
		}
	}
	return propertySchema{
		BsonType:        "object",
		Properties:      props,
		RequiredProps:   requiredProps,
		AdditionalProps: false,
	}
}

func loadJSONSchema(filePath string) propertySchema {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	var schema propertySchema
	json.Unmarshal(byteValue, &schema)
	return schema
}

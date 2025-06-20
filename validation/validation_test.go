package validation

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"opensignauxfaibles/lib/apconso"
	"opensignauxfaibles/lib/apdemande"
	"opensignauxfaibles/lib/bdf"
	"opensignauxfaibles/lib/diane"
	"opensignauxfaibles/lib/ellisphere"
	"opensignauxfaibles/lib/paydex"
	"opensignauxfaibles/lib/sirene"
	sireneul "opensignauxfaibles/lib/sirene_ul"
	"opensignauxfaibles/lib/urssaf"
)

var _ = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestDiffSchema(t *testing.T) {
	t.Run("doit vérifier qu'un champ est optionnel s'il est taggé avec \"omitempty\"", func(t *testing.T) {
		type MyType struct {
			MandatoryField string `bson:"mandatoryField"`
			OptionalField  string `bson:"optionalField,omitempty"`
		}
		schema := propertySchema{
			BsonType: "object",
			Properties: map[string]propertySchema{
				"mandatoryField": {BsonType: "string"},
				"optionalField":  {BsonType: "string"},
			},
			RequiredProps:   []string{"mandatoryField"},
			AdditionalProps: false,
		}
		assert.Equal(t, []error{}, diffSchema(schema, reflectStructType(reflect.TypeOf(MyType{}))))
	})
}

func TestDiffMaps(t *testing.T) {
	t.Run("doit détecter une entrée manquante dans la map A", func(t *testing.T) {
		schemaProps := map[string]propertySchema{"a": {}, "b": {}}
		structProps := map[string]propertySchema{"a": {}, "b": {}, "c": {}}
		assert.ElementsMatch(t, diffMaps(schemaProps, structProps), []error{
			errors.New("property not found in JSON Schema: c"),
		})
	})
	t.Run("doit détecter une entrée manquante dans la map B", func(t *testing.T) {
		schemaProps := map[string]propertySchema{"a": {}, "b": {}, "c": {}}
		structProps := map[string]propertySchema{"a": {}, "b": {}}
		assert.ElementsMatch(t, diffMaps(schemaProps, structProps), []error{
			errors.New("property not found in Go struct: c"),
		})
	})
	t.Run("doit détecter une entrée dont le type ne correspond pas", func(t *testing.T) {
		schemaProps := map[string]propertySchema{"a": {BsonType: "number"}}
		structProps := map[string]propertySchema{"a": {BsonType: "string"}}
		assert.ElementsMatch(t, diffMaps(schemaProps, structProps), []error{
			errors.New("property types of \"a\" don't match: {number map[] [] false} <> {string map[] [] false}"),
		})
	})
}

func TestReflectStructType(t *testing.T) {
	t.Run("doit inclure tous les champs dans \"required\", sauf ceux taggés avec \"omitempty\"", func(t *testing.T) {
		type MyType struct {
			MandatoryField string `bson:"mandatoryField"`
			OptionalField  string `bson:"optionalField,omitempty"`
		}
		expectedSchema := propertySchema{
			BsonType: "object",
			Properties: map[string]propertySchema{
				"mandatoryField": {BsonType: "string"},
				"optionalField":  {BsonType: "string"},
			},
			RequiredProps:   []string{"mandatoryField"},
			AdditionalProps: false,
		}
		assert.Equal(t, expectedSchema, reflectStructType(reflect.TypeOf(MyType{})))
	})
}

func TestReflectPropsFromStruct(t *testing.T) {
	t.Run("doit extraire le nom JSON et type d'un champ de struct Go", func(t *testing.T) {
		type MyType struct {
			MyField string `bson:"myField"`
		}
		assert.Equal(t, map[string]propertySchema{"myField": {BsonType: "string"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit interpréter les types float64 et int comme number", func(t *testing.T) {
		type MyType struct {
			MyField1 int     `bson:"f1"`
			MyField2 float64 `bson:"f2"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {BsonType: "number"}, "f2": {BsonType: "number"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit reconnaitre le type des pointeurs", func(t *testing.T) {
		type MyType struct {
			MyField1 *int `bson:"f1"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {BsonType: "number"}}, reflectPropsFromStruct(MyType{}))
	})
}

func TestTypeAlignment(t *testing.T) {

	typesToCompare := map[string]interface{}{
		"apconso.schema.json":      apconso.APConso{},
		"apdemande.schema.json":    apdemande.APDemande{},
		"bdf.schema.json":          bdf.BDF{},
		"ccsf.schema.json":         urssaf.CCSF{},
		"compte.schema.json":       urssaf.Compte{},
		"cotisation.schema.json":   urssaf.Cotisation{},
		"debit.schema.json":        urssaf.Debit{},
		"delai.schema.json":        urssaf.Delai{},
		"diane.schema.json":        diane.Diane{},
		"effectif.schema.json":     urssaf.Effectif{},
		"effectif_ent.schema.json": urssaf.EffectifEnt{},
		"ellisphere.schema.json":   ellisphere.Ellisphere{},
		"paydex.schema.json":       paydex.Paydex{},
		"procol.schema.json":       urssaf.Procol{},
		"sirene.schema.json":       sirene.Sirene{},
		"sirene_ul.schema.json":    sireneul.SireneUL{},
		// NOTE: Au fur et à mesure qu'on ajoute des fichiers JSON Schema, penser à les couvrir ici.
	}

	t.Run("chaque fichier JSON Schema est couvert par un test d'alignement avec son type Go correspondant", func(t *testing.T) {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".schema.json") {
				jsonSchemaFile := file.Name()
				if _, ok := typesToCompare[jsonSchemaFile]; !ok {
					assert.Fail(t, "please add \""+jsonSchemaFile+"\" entry to typesToCompare, in validation_test.go")
				}
			}
		}
	})

	t.Run("chaque fichier JSON Schema est aligné avec le type Go retourné par le parseur correspondant", func(t *testing.T) {
		for jsonTypeName, structInstance := range typesToCompare {
			// TODO: remettre en marche le test apdemande
			// nécessite de revoir la validation de schéma pour prendre en compte les types nullables (*int et *float64 peuvent être null)
			if jsonTypeName == "apdemande.schema.json" {
				continue
			}
			t.Run(jsonTypeName, func(t *testing.T) {
				errors := diffTypeSchema(jsonTypeName, structInstance)
				if ok := assert.ElementsMatch(t, []error{}, errors); !ok {
					// affichage des champs non alignés, pour aider à la complétion
					structTypeName := reflect.TypeOf(structInstance).Name()
					t.Log(jsonTypeName + " is not aligned with struct type \"" + structTypeName + "\":")
					for _, err := range errors {
						t.Log("- " + err.Error())
					}
				}
			})
		}
	})

	t.Run("tout champ manquant dans le schema doit être rapporté", func(t *testing.T) {
		type Dummy struct {
			A string    `bson:"a"`
			B time.Time `bson:"b,omitempty"`
			C float64   `bson:"c"`
			D bool      `bson:"d,omitempty"`
			E bool      `bson:"e,omitempty"`
			Z bool
		}
		structSchema := reflectStructType(reflect.TypeOf(Dummy{}))
		var jsonSchema propertySchema
		json.Unmarshal([]byte(`{
			"title": "EntréeDummy",
			"bsonType": "object",
			"required": ["a"],
			"properties": {
				"a": {
					"bsonType": "string",
					"pattern": "^[0-9]{9}$"
				},
				"b": {
					"bsonType": "date"
				},
				"c": {
					"bsonType": "number"
				},
				"e": {
					"bsonType": "object"
				}
			},
			"additionalProperties": false
		}`), &jsonSchema)
		expectedErrors := []error{
			errors.New("property not found in JSON Schema: d"),
			errors.New("property types of \"e\" don't match: {object map[] [] false} <> {bool map[] [] false}"),
			errors.New("property not marked as 'required' in JSON Schema: c"),
		}
		errors := diffSchema(jsonSchema, structSchema)
		if ok := assert.ElementsMatch(t, expectedErrors, errors); !ok {
			// affichage des champs non alignés, pour aider à la complétion
			t.Log("JSON Schema is not aligned with struct type:")
			for _, err := range errors {
				t.Log("- " + err.Error())
			}
		}
	})
}

func diffTypeSchema(jsonSchemaFile string, structInstance interface{}) []error {
	jsonSchema := loadJSONSchema(jsonSchemaFile)
	structSchema := reflectStructType(reflect.TypeOf(structInstance))
	return diffSchema(jsonSchema, structSchema)
}

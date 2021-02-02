package validation

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

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
		schemaProps := map[string]propertySchema{"a": {"number"}}
		structProps := map[string]propertySchema{"a": {"string"}}
		assert.ElementsMatch(t, diffMaps(schemaProps, structProps), []error{
			errors.New("property types of \"a\" don't match: {number} <> {string}"),
		})
	})
}

func TestReflectPropsFromStruct(t *testing.T) {
	t.Run("doit extraire le nom JSON et type d'un champ de struct Go", func(t *testing.T) {
		type MyType struct {
			MyField string `json:"myField"`
		}
		assert.Equal(t, map[string]propertySchema{"myField": {"string"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit interpréter les types float64 et int comme number", func(t *testing.T) {
		type MyType struct {
			MyField1 int     `json:"f1"`
			MyField2 float64 `json:"f2"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {"number"}, "f2": {"number"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit reconnaitre le type des pointeurs", func(t *testing.T) {
		type MyType struct {
			MyField1 *int `json:"f1"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {"number"}}, reflectPropsFromStruct(MyType{}))
	})
}

func TestTypeAlignment(t *testing.T) {

	typesToCompare := map[string]interface{}{
		"delai.schema.json": urssaf.Delai{},
	}

	t.Run("tous les fichiers JSON Schema sont couverts par des tests", func(t *testing.T) {
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

	t.Run("chaque type décrit en JSON Schema doit correspondre à la structure retournée par le parseur correspondant", func(t *testing.T) {
		for jsonTypeName, structInstance := range typesToCompare {
			t.Run(jsonTypeName, func(t *testing.T) {
				errors := diffProps(jsonTypeName, structInstance)
				if len(errors) > 0 {
					log.Println("Types are not deeply equal:")
					for _, err := range errors {
						log.Println("- " + err.Error())
					}
					t.FailNow()
				}
			})
		}
	})

	t.Run("le type BDF n'est pas encore complètement couvert en JSON Schema", func(t *testing.T) {
		actualErrors := diffProps("bdf.schema.json", bdf.BDF{})
		assert.ElementsMatch(t, actualErrors, []error{
			errors.New("property not found in JSON Schema: delai_fournisseur"),
			errors.New("property not found in JSON Schema: dette_fiscale"),
			errors.New("property not found in JSON Schema: frais_financier"),
			errors.New("property not found in JSON Schema: arrete_bilan_bdf"),
			errors.New("property not found in JSON Schema: secteur"),
			errors.New("property not found in JSON Schema: taux_marge"),
			errors.New("property not found in JSON Schema: annee_bdf"),
			errors.New("property not found in JSON Schema: raison_sociale"),
			errors.New("property not found in JSON Schema: poids_frng"),
			errors.New("property not found in JSON Schema: financier_court_terme"),
		})
	})
}

func diffProps(jsonSchemaFile string, structInstance interface{}) []error {
	schemaProps := loadPropsFromSchema(jsonSchemaFile)
	structProps := reflectPropsFromStruct(structInstance)
	return diffMaps(schemaProps, structProps)
}

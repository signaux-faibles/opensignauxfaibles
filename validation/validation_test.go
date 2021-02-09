package validation

import (
	"errors"
	"flag"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/lib/diane"
	"github.com/signaux-faibles/opensignauxfaibles/lib/ellisphere"
	"github.com/signaux-faibles/opensignauxfaibles/lib/paydex"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sirene"
	sireneul "github.com/signaux-faibles/opensignauxfaibles/lib/sirene_ul"
	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file") // please keep this line until https://github.com/kubernetes-sigs/service-catalog/issues/2319#issuecomment-425200065 is fixed

func TestDiffSchema(t *testing.T) {
	t.Run("doit vérifier qu'un champ est optionnel s'il est taggé avec \"omitempty\"", func(t *testing.T) {
		type MyType struct {
			MandatoryField string `json:"mandatoryField"`
			OptionalField  string `json:"optionalField,omitempty"`
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
			MandatoryField string `json:"mandatoryField"`
			OptionalField  string `json:"optionalField,omitempty"`
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
			MyField string `json:"myField"`
		}
		assert.Equal(t, map[string]propertySchema{"myField": {BsonType: "string"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit interpréter les types float64 et int comme number", func(t *testing.T) {
		type MyType struct {
			MyField1 int     `json:"f1"`
			MyField2 float64 `json:"f2"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {BsonType: "number"}, "f2": {BsonType: "number"}}, reflectPropsFromStruct(MyType{}))
	})
	t.Run("doit reconnaitre le type des pointeurs", func(t *testing.T) {
		type MyType struct {
			MyField1 *int `json:"f1"`
		}
		assert.Equal(t, map[string]propertySchema{"f1": {BsonType: "number"}}, reflectPropsFromStruct(MyType{}))
	})
}

func TestTypeAlignment(t *testing.T) {

	type TypeToCompare struct {
		ParserStructInstance interface{} // instance type Go retourné par le parseur correspondant à un JSON Schema donné
		ExpectedErrors       []error     // erreurs attendues lors de la vérification d'alignement entre JSON Schema et type Go
	}

	typesToCompare := map[string]TypeToCompare{
		"apconso.schema.json":      {apconso.APConso{}, []error{}},
		"apdemande.schema.json":    {apdemande.APDemande{}, []error{}},
		"ccsf.schema.json":         {urssaf.CCSF{}, []error{}},
		"compte.schema.json":       {urssaf.Compte{}, []error{}},
		"cotisation.schema.json":   {urssaf.Cotisation{}, []error{}},
		"debit.schema.json":        {urssaf.Debit{}, []error{}},
		"delai.schema.json":        {urssaf.Delai{}, []error{}},
		"diane.schema.json":        {diane.Diane{}, []error{}},
		"effectif.schema.json":     {urssaf.Effectif{}, []error{}},
		"effectif_ent.schema.json": {urssaf.EffectifEnt{}, []error{}},
		"ellisphere.schema.json":   {ellisphere.Ellisphere{}, []error{}},
		"paydex.schema.json":       {paydex.Paydex{}, []error{}},
		"procol.schema.json":       {urssaf.Procol{}, []error{}},
		"sirene.schema.json":       {sirene.Sirene{}, []error{}},
		"sirene_ul.schema.json":    {sireneul.SireneUL{}, []error{}},
		"bdf.schema.json": {bdf.BDF{}, []error{ // bdf.schema.json n'est pas encore complet => la vérification va retourner les erreurs suivantes:
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
			errors.New("property not marked as 'required' in JSON Schema: annee_bdf"),
			errors.New("property not marked as 'required' in JSON Schema: arrete_bilan_bdf"),
			errors.New("property not marked as 'required' in JSON Schema: raison_sociale"),
			errors.New("property not marked as 'required' in JSON Schema: secteur"),
			errors.New("property not marked as 'required' in JSON Schema: poids_frng"),
			errors.New("property not marked as 'required' in JSON Schema: taux_marge"),
			errors.New("property not marked as 'required' in JSON Schema: delai_fournisseur"),
			errors.New("property not marked as 'required' in JSON Schema: dette_fiscale"),
			errors.New("property not marked as 'required' in JSON Schema: financier_court_terme"),
			errors.New("property not marked as 'required' in JSON Schema: frais_financier"),
		}},
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
		for jsonTypeName, typeDef := range typesToCompare {
			t.Run(jsonTypeName, func(t *testing.T) {
				errors := diffTypeSchema(jsonTypeName, typeDef.ParserStructInstance)
				if ok := assert.ElementsMatch(t, typeDef.ExpectedErrors, errors); !ok {
					// affichage des champs non alignés, pour aider à la complétion
					structTypeName := reflect.TypeOf(typeDef.ParserStructInstance).Name()
					t.Log(jsonTypeName + " is not aligned with struct type \"" + structTypeName + "\":")
					for _, err := range errors {
						t.Log("- " + err.Error())
					}
				}
			})
		}
	})
}

func diffTypeSchema(jsonSchemaFile string, structInstance interface{}) []error {
	jsonSchema := loadJSONSchema(jsonSchemaFile)
	structSchema := reflectStructType(reflect.TypeOf(structInstance))
	return diffSchema(jsonSchema, structSchema)
}

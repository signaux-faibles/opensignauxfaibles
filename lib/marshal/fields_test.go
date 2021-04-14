package marshal

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexColumnsFromCsvHeader(t *testing.T) {

	reader := makeCsvReader("year", "2020")
	idx, err := IndexColumnsFromCsvHeader(reader, struct{}{})

	if assert.NoError(t, err) {
		t.Run("permet de vérifier la présence de colonnes depuis un fichier CSV valide", func(t *testing.T) {
			hasYear, err := idx.HasFields([]string{"year"})
			assert.Equal(t, true, hasYear)
			assert.Equal(t, nil, err)
			hasBothFields, err := idx.HasFields([]string{"year", "dummy"})
			assert.Equal(t, false, hasBothFields)
			assert.EqualError(t, err, "Colonne dummy non trouvée. Abandon.")
		})

		t.Run("permet d'indexer chaque ligne depuis un fichier CSV valide", func(t *testing.T) {
			indexedRow := idx.IndexRow([]string{"2020"})
			assert.Equal(t, "2020", indexedRow.GetVal("year"))
		})
	}

}

func TestIndexedRow(t *testing.T) {

	idx := CreateColMapping(map[string]int{"data": 0}) // simule un fichier csv avec une colonne "data" à l'indice 0

	t.Run("GetOptionalVal retourne la valeur ou une chaine vide, selon la présence de la colonne", func(t *testing.T) {
		row := idx.IndexRow([]string{"2020"})
		dataVal, dataFound := row.GetOptionalVal("data")
		assert.Equal(t, "2020", dataVal)
		assert.Equal(t, true, dataFound)
		undefVal, undefFound := row.GetOptionalVal("colonne_inexistante")
		assert.Equal(t, "", undefVal)
		assert.Equal(t, false, undefFound)
	})

	t.Run("GetFloat64 retourne la valeur décimale ou nil (avec une erreur), selon la présence de la colonne", func(t *testing.T) {
		row := idx.IndexRow([]string{"20.20"})
		dataVal, err := row.GetFloat64("data")
		assert.Equal(t, 20.20, *dataVal)
		assert.Equal(t, nil, err)
		undefVal, undefErr := row.GetFloat64("colonne_inexistante")
		assert.Nil(t, undefVal)
		assert.EqualError(t, undefErr, "GetFloat64 failed to find column: colonne_inexistante")
	})

	t.Run("GetCommaFloat64 retourne la valeur décimale à virgule ou nil (avec une erreur), selon la présence de la colonne", func(t *testing.T) {
		row := idx.IndexRow([]string{"20,20"})
		dataVal, err := row.GetCommaFloat64("data")
		assert.Equal(t, 20.20, *dataVal)
		assert.Equal(t, nil, err)
		undefVal, undefErr := row.GetCommaFloat64("colonne_inexistante")
		assert.Nil(t, undefVal)
		assert.EqualError(t, undefErr, "GetCommaFloat64 failed to find column: colonne_inexistante")
	})

	t.Run("GetInt retourne la valeur entière ou nil (avec une erreur), selon la présence de la colonne", func(t *testing.T) {
		row := idx.IndexRow([]string{"20"})
		dataVal, err := row.GetInt("data")
		assert.Equal(t, 20, *dataVal)
		assert.Equal(t, nil, err)
		undefVal, undefErr := row.GetInt("colonne_inexistante")
		assert.Nil(t, undefVal)
		assert.EqualError(t, undefErr, "GetInt failed to find column: colonne_inexistante")
	})

	t.Run("GetBool retourne la valeur booléenne ou false (avec une erreur), selon la présence de la colonne", func(t *testing.T) {
		row := idx.IndexRow([]string{"true"})
		dataVal, err := row.GetBool("data")
		assert.Equal(t, true, dataVal)
		assert.Equal(t, nil, err)
		undefVal, undefErr := row.GetBool("colonne_inexistante")
		assert.Equal(t, false, undefVal)
		assert.EqualError(t, undefErr, "GetBool failed to find column: colonne_inexistante")
	})
}

func makeCsvReader(headerRow string, dataRow string) *csv.Reader {
	return csv.NewReader(strings.NewReader(headerRow + "\n" + dataRow))
}

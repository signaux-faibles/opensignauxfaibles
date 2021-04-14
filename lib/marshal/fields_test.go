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

	idx := CreateColMapping(map[string]int{"year": 0}) // simule un fichier csv avec une colonne "year" à l'indice 0

	t.Run("GetOptionalVal retourne une valeur si la colonne est trouvée", func(t *testing.T) {
		dataRow := "2020"
		row := idx.IndexRow([]string{dataRow})
		yearVal, yearFound := row.GetOptionalVal("year")
		assert.Equal(t, "2020", yearVal)
		assert.Equal(t, true, yearFound)
		undefVal, undefFound := row.GetOptionalVal("colonne_inexistante")
		assert.Equal(t, "", undefVal)
		assert.Equal(t, false, undefFound)
	})
}

func makeCsvReader(headerRow string, dataRow string) *csv.Reader {
	return csv.NewReader(strings.NewReader(headerRow + "\n" + dataRow))
}

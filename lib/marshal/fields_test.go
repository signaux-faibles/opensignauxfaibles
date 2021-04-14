package marshal

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexColumnsFromCsvHeader(t *testing.T) {

	t.Run("retourne un index fonctionnel des champs détectés, depuis un fichier CSV valide", func(t *testing.T) {
		type TupleStruct struct{}
		reader := csv.NewReader(strings.NewReader("year\n2020"))
		idx, err := IndexColumnsFromCsvHeader(reader, TupleStruct{})
		if assert.NoError(t, err) {
			hasYear, err := idx.HasFields([]string{"year"})
			assert.Equal(t, true, hasYear)
			assert.Equal(t, nil, err)
			hasBothFields, err := idx.HasFields([]string{"year", "dummy"})
			assert.Equal(t, false, hasBothFields)
			assert.EqualError(t, err, "Colonne dummy non trouvée. Abandon.")
		}
	})
}

func TestIndexedRow(t *testing.T) {

	t.Run("GetOptionalVal retourne une valeur si la colonne est trouvée", func(t *testing.T) {
		type TupleStruct struct{}
		headerRow := "year"
		dataRow := "2020"
		reader := csv.NewReader(strings.NewReader(headerRow + "\n" + dataRow))
		idx, err := IndexColumnsFromCsvHeader(reader, TupleStruct{})
		if assert.NoError(t, err) {
			row := idx.IndexRow([]string{dataRow})
			yearVal, yearFound := row.GetOptionalVal("year")
			assert.Equal(t, "2020", yearVal)
			assert.Equal(t, true, yearFound)
			undefVal, undefFound := row.GetOptionalVal("colonne_inexistante")
			assert.Equal(t, "", undefVal)
			assert.Equal(t, false, undefFound)
		}
	})
}

package minoos

import (
	"encoding/csv"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/tools/altares/test"
)

func Test_something(t *testing.T) {
	mc := NewWithClient(test.NewS3ForTest(t))
	expectedLines := 8000
	stockCSV := test.GenerateStockCSV(expectedLines)
	stats, err := stockCSV.Stat()
	require.NoError(t, err)
	slog.Debug("TU - stats du fichier csv généré", slog.Any("size", stats.Size()), slog.String("name", stats.Name()))
	mc.PutAltaresFile("a_file.csv", stockCSV)

	files := mc.ListAltaresFiles()
	assert.Len(t, files, 1)
	assert.Equal(t, "altares/a_file.csv", files[0])
	file := mc.GetAltaresFile("a_file.csv")
	assert.NotNil(t, file)

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, expectedLines+1) // +1 pour la ligne des headers
	assert.Len(t, records[0], 11)
}

func Test_newLocalS3(t *testing.T) {
	// Initialize minio client object.
	minioClient := test.NewS3ForTest(t)
	mc := NewWithClient(minioClient)
	assert.NotNil(t, mc)
}

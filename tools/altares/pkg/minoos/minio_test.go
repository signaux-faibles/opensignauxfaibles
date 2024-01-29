package minoos

import (
	"encoding/csv"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/tools/altares/test"
)

func Test_put_and_read_one_file(t *testing.T) {
	mc := New(test.NewS3ForTest(t), test.FakeBucketName())
	//t.Cleanup(mc.CleanupVersionedBucket)
	expectedLines := 8000
	stockCSV := test.GenerateStockCSV(expectedLines)
	stats, err := stockCSV.Stat()
	require.NoError(t, err)
	slog.Debug("TU - stats du fichier csv généré", slog.Any("size", stats.Size()), slog.String("name", stats.Name()))
	mc.PutAltaresFile("stock.csv", stockCSV)

	files := mc.ListAltaresFiles()
	assert.Len(t, files, 1)
	assert.Equal(t, "altares/stock.csv", files[0])
	file := mc.GetAltaresFile("altares/stock.csv")
	assert.NotNil(t, file)

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	require.NoError(t, err)
	assert.Len(t, records, expectedLines+1) // +1 pour la ligne des headers
	assert.Len(t, records[0], 11)
}

func Test_put_and_list_many_files(t *testing.T) {
	mc := New(test.NewS3ForTest(t), test.FakeBucketName())
	//t.Cleanup(mc.CleanupVersionedBucket)
	stock := test.CreateRandomFile()
	mc.PutAltaresFile("stock.csv", stock)
	increment1 := test.CreateRandomFile()
	mc.PutAltaresFile("2312/increment1", increment1)
	increment2 := test.CreateRandomFile()
	mc.PutAltaresFile("2401/increment2", increment2)

	files := mc.ListAltaresFiles()
	assert.Len(t, files, 3)
	assert.Contains(t, files, "altares/stock.csv")
	assert.Contains(t, files, "altares/2312/increment1")
	assert.Contains(t, files, "altares/2401/increment2")
}

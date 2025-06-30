package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/marshal"
)

var fakeCsv faker.Faker

func init() {
	fakeCsv = faker.New()
}

func Test_writeLinesToCSV(t *testing.T) {
	exportPath := filepath.Join(os.TempDir(), fakeCsv.Lorem().Word())
	viper.Set("export.path", exportPath)
	batchKey := BatchKey("2310")

	tuple := Test{}
	tuples := map[string]marshal.Tuple{
		fakeCsv.Lorem().Word(): tuple,
	}
	writeLinesToCSV(batchKey, tuples)
	assert.FileExists(t, filepath.Join(exportPath, batchKey, tuple.Type()+".csv"))
}

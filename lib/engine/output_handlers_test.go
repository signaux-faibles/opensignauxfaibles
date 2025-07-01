package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var fakeCsv faker.Faker

type TestTuple struct{}

func (TestTuple) Key() string {
	return "000000000"
}

func (TestTuple) Scope() string {
	return "entreprise"
}

func (TestTuple) Type() string {
	return "test"
}

func (TestTuple) Headers() []string {
	return []string{"test_header"}
}
func (TestTuple) Values() []string {
	return []string{"test_value"}
}

func init() {
	fakeCsv = faker.New()
}

func Test_writeLinesToCSV(t *testing.T) {
	exportPath := filepath.Join(os.TempDir(), fakeCsv.Lorem().Word())
	viper.Set("export.path", exportPath)
	batchKey := "2310"

	tuple := TestTuple{}
	writeLineToCSV(batchKey, tuple)
	assert.FileExists(t, filepath.Join(exportPath, batchKey, tuple.Type()+".csv"))
}

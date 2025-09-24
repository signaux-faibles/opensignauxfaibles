package prepareimport

import (
	"bytes"
	"compress/gzip"
	"opensignauxfaibles/lib/base"
	"os"
	"path/filepath"
	"testing"

	"github.com/jaswdr/faker"
)

var fake faker.Faker

func init() {
	fake = faker.New()
}

var dummyBatchKey = base.NewSafeBatchKey("1802")

// CreateTempFiles creates a temporary directory with a batch of files, and clean up after the execution of tests
func CreateTempFiles(t *testing.T, batchkey base.BatchKey, filenames []string) string {
	contentPerFile := map[string][]byte{}
	for _, filename := range filenames {
		contentPerFile[filename] = []byte{}
	}
	return CreateTempFilesWithContent(t, batchkey, contentPerFile)
}

// CreateTempFilesWithContent creates a temporary directory with a batch of files, and clean up after the execution of tests
func CreateTempFilesWithContent(t *testing.T, batchkey base.BatchKey, contentPerFile map[string][]byte) string {
	t.Helper()
	parentDir, err := os.MkdirTemp(os.TempDir(), "example")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Cleanup(func() { _ = os.RemoveAll(parentDir) })

	batchDir := filepath.Join(parentDir, batchkey.String())
	_ = os.Mkdir(batchDir, 0777)
	for filename := range contentPerFile {
		tmpFilename := filepath.Join(batchDir, filename)
		if err := os.WriteFile(tmpFilename, contentPerFile[filename], 0666); err != nil {
			t.Fatal(err.Error())
		}
	}

	return parentDir
}

func GzipString(src string) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write([]byte(src))
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func SomeText(length int) string {
	return fake.Lorem().Text(length)
}

func SomeTextAsBytes(length int) []byte {
	return fake.Lorem().Bytes(length)
}

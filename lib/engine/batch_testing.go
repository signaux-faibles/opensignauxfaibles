package engine

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func NewSafeBatchKey(key string) BatchKey {
	batchKey, err := NewBatchKey(key)
	if err != nil {
		log.Fatal(err)
	}
	return batchKey
}

func MockBatch(filetype ParserType, batchFiles []BatchFile) AdminBatch {
	batch := AdminBatch{
		Key:   "1902",
		Files: BatchFiles{filetype: batchFiles},
		Params: AdminBatchParams{
			DateDebut: time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // January 1st, 2019
			DateFin:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // February 1st, 2019
		},
	}
	return batch
}

// -----------------------------------------------------

type MockBatchFile struct {
	content string
}

func (MockBatchFile) Filename() string   { return "mockfile" }
func (MockBatchFile) Path() string       { return "./mockfile" }
func (MockBatchFile) IsCompressed() bool { return false }

func (m MockBatchFile) Open() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.content)), nil
}

func NewMockBatchFile(content string) BatchFile {
	return MockBatchFile{content}
}

// -----------------------------------------------------

type OpenFailsBatchFile struct {
	MockBatchFile
}

func (m OpenFailsBatchFile) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("error from Open()")
}
package engine

import (
	"log"
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


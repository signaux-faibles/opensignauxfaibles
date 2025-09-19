package base

import "time"

func MockBatch(filetype ParserType, filepaths []string) AdminBatch {
	batchFiles := []BatchFile{}
	for _, file := range filepaths {
		batchFiles = append(batchFiles, NewDummyBatchFile(file))
	}
	batch := AdminBatch{
		Files: BatchFiles{filetype: batchFiles},
		Params: AdminBatchParams{
			DateDebut: time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // January 1st, 2019
			DateFin:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // February 1st, 2019
		},
	}
	return batch
}

// NewDummyBatchFile creates a BatchFile for testing contexts where the actual
// basepath does not matter
func NewDummyBatchFile(relativepath string) BatchFile {
	return NewBatchFile(".", relativepath)
}

func NewDummyBatchFileFromBatch(relativepath string, batch BatchKey) BatchFile {
	return NewBatchFileFromBatch(".", batch, relativepath)
}

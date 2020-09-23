package base

import (
	"errors"
	"time"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID    `json:"id" bson:"_id"`
	Files         BatchFiles `json:"files" bson:"files"`
	Name          string     `json:"name" bson:"name"`
	Readonly      bool       `json:"readonly" bson:"readonly"`
	CompleteTypes []string   `json:"complete_types" bson:"complete_types"`
	Params        struct {
		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
	} `json:"params" bson:"param"`
}

// New crée un nouveau batch
func (batch *AdminBatch) New(batchKey string) error {
	if !IsBatchID(batchKey) {
		return errors.New("Valeur de batch non autorisée")
	}
	batch.ID.Key = batchKey
	batch.ID.Type = "batch"
	batch.Files = BatchFiles{}
	return nil
}

func IsBatchID(batchID string) bool {
	if len(batchID) < 4 {
		return false
	}
	_, err := time.Parse("0601", batchID[0:4])
	if len(batchID) > 4 && batchID[4] != '_' {
		return false
	}
	return err == nil
}

// AdminID Collection key
type AdminID struct {
	Key  string `json:"key" bson:"key"`
	Type string `json:"type" bson:"type"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

// MockBatch with a map[type][]filepaths
func MockBatch(filetype string, filepaths []string) AdminBatch {
	fileMap := map[string][]string{filetype: filepaths}
	batch := AdminBatch{
		Files: BatchFiles(fileMap),
	}
	return batch
}

package base

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID          `json:"id" bson:"_id"`
	Files         BatchFiles       `json:"files" bson:"files"`
	Name          string           `json:"name" bson:"name"`
	Readonly      bool             `json:"readonly" bson:"readonly"`
	CompleteTypes []string         `json:"complete_types" bson:"complete_types"`
	Params        adminBatchParams `json:"params" bson:"param"`
}

type adminBatchParams struct {
	DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
	DateFin         time.Time `json:"date_fin" bson:"date_fin"`
	DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
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

// IsBatchID retourne `true` si `batchID` est un identifiant de Batch.
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
type BatchFiles map[string][]BatchFile

type BatchFile string

func (file BatchFile) FilePath() string {
	return rePrefix.ReplaceAllString(string(file), "") // c.a.d. suppression du préfixe éventuellement trouvé
}

func (file BatchFile) IsCompressed() bool {
	return file.Prefix() == "gzip:" || strings.HasSuffix(string(file), ".gz")
}

func (file BatchFile) Prefix() string {
	return rePrefix.FindString(string(file))
}

var rePrefix = regexp.MustCompile("^[a-z]*:")

// MockBatch with a map[type][]filepaths
func MockBatch(filetype string, filepaths []BatchFile) AdminBatch {
	fileMap := map[string][]BatchFile{filetype: filepaths}
	batch := AdminBatch{
		Files: BatchFiles(fileMap),
		Params: adminBatchParams{
			DateDebut: time.Date(2019, 0, 1, 0, 0, 0, 0, time.UTC), // January 1st, 2019
			DateFin:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), // February 1st, 2019
		},
	}
	return batch
}

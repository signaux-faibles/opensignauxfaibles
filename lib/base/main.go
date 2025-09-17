package base

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID     AdminID          `json:"id"`
	Files  BatchFiles       `json:"files"`
	Params AdminBatchParams `json:"params"`
}

type AdminBatchParams struct {
	DateDebut       time.Time `json:"date_debut"`
	DateFin         time.Time `json:"date_fin"`
	DateFinEffectif time.Time `json:"date_fin_effectif"`
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

// AdminID represents the "id" property of an AdminBatch object.
type AdminID struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[ParserType][]BatchFile

// BatchFile encapsule un fichier mentionné dans un Batch
type BatchFile string

var reScheme = regexp.MustCompile("^[a-z]*:")

func (file BatchFile) scheme() string {
	scheme := reScheme.FindString(string(file))
	return scheme
}

func (file BatchFile) RelativePath() string {
	pathWithoutScheme := reScheme.ReplaceAllString(string(file), "")
	return pathWithoutScheme
}

// FilePath retourne le chemin vers le fichier, sans le schéma
// (base path)
func (file BatchFile) FilePath() string {
	return filepath.Join(viper.GetString("APP_DATA"), file.RelativePath())
}

// IsCompressed est vrai si le fichier est compressé
func (file BatchFile) IsCompressed() bool {
	return file.scheme() == "gzip:" ||
		strings.HasSuffix(file.RelativePath(), ".gz")
}

// MockBatch with a map[type][]filepaths
func MockBatch(filetype ParserType, filepaths []string) AdminBatch {
	batchFiles := []BatchFile{}
	for _, file := range filepaths {
		batchFiles = append(batchFiles, BatchFile(file))
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

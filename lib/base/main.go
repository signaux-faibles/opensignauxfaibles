package base

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID          `json:"id"`
	Files         BatchFiles       `json:"files"`
	Name          string           `json:"name"`
	Readonly      bool             `json:"readonly"`
	CompleteTypes []string         `json:"complete_types"`
	Params        AdminBatchParams `json:"params"`
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
type BatchFiles map[string][]BatchFile

// BatchFile encapsule un fichier mentionné dans un Batch
type BatchFile struct {
	// Scheme prefix, e.g. `gzip:`
	Scheme string

	// Directory path where data is to be looked for
	basePath string

	// Relative path inside the base path
	relativePath string
}

var reScheme = regexp.MustCompile("^[a-z]*:")

// NewBatchFile créé un chemin d'un fichier encapsulé dans un batch avec la
// variable d'environnement "APP_DATA" comme chemin de base
func NewBatchFile(path string) BatchFile {
	return NewBatchFileWithBasePath(path, viper.GetString("APP_DATA"))
}

// MarshalJSON implements json.Marshaller interface
func (file BatchFile) MarshalJSON() ([]byte, error) {
	value := file.Scheme + file.relativePath
	return json.Marshal(value)
}

// UnmarshalJSON implements json.Unmarshaller interface
func (file *BatchFile) UnmarshalJSON(raw []byte) error {
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}

	batchFile := NewBatchFile(value)
	file.Scheme = batchFile.Scheme
	file.relativePath = batchFile.relativePath
	file.basePath = batchFile.basePath

	return nil
}

// NewBatchFileWithBasePath crée un chemin d'un fichier encapsulé dans un
// batch
func NewBatchFileWithBasePath(path, basePath string) BatchFile {
	scheme := reScheme.FindString(path)
	pathWithoutScheme := reScheme.ReplaceAllString(path, "")
	return BatchFile{Scheme: scheme, relativePath: pathWithoutScheme, basePath: basePath}
}

// FilePath retourne le chemin vers le fichier, sans le schéma
// (base path)
func (file BatchFile) FilePath() string {
	return filepath.Join(file.basePath, file.relativePath)
}

func (file BatchFile) RelativePath() string {
	return file.relativePath
}

// IsCompressed est vrai si le fichier est compressé
func (file BatchFile) IsCompressed() bool {
	return file.Scheme == "gzip:" ||
		strings.HasSuffix(file.relativePath, ".gz")
}

// MockBatch with a map[type][]filepaths
func MockBatch(filetype string, filepaths []string) AdminBatch {
	batchFiles := []BatchFile{}
	for _, file := range filepaths {
		batchFiles = append(batchFiles, NewBatchFile(file))
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

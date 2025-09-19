package base

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	Key    BatchKey         `json:"key"`
	Files  BatchFiles       `json:"files"`
	Params AdminBatchParams `json:"params"`
}

type AdminBatchParams struct {
	DateDebut time.Time `json:"date_debut"`
	DateFin   time.Time `json:"date_fin"`
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

// BatchFiles fichiers mappés par type
type BatchFiles map[ParserType][]BatchFile

// GetFilterFile returns the filter file.
func (files BatchFiles) GetFilterFile() (BatchFile, error) {
	if files["filter"] == nil || len(files["filter"]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 filter file, found: %s", files["filter"])
	}
	return files["filter"][0], nil
}

func (files BatchFiles) GetSireneULFile() (BatchFile, error) {
	if files[SireneUl] == nil || len(files[SireneUl]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 sireneUL filter file, found %s", files[SireneUl])
	}
	return files[SireneUl][0], nil
}

// GetEffectifFile returns the effectif file.
func (files BatchFiles) GetEffectifFile() (BatchFile, error) {
	if files["effectif"] == nil || len(files["effectif"]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 effectif file, found: %s", files["effectif"])
	}
	return files["effectif"][0], nil
}

// BatchFile is the relevant metadata of a file that needs to be imported
type BatchFile interface {
	Filename() string

	// Relative path given a base path
	// Usually begins with the batch key
	RelativePath() string

	AbsolutePath() string
	IsCompressed() bool
}

type batchFile struct {
	basepath     string
	relativePath string
}

func (file *batchFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(file.RelativePath())
}

func (file *batchFile) UnmarshalJSON(pathBytes []byte) error {
	file.basepath = viper.GetString("APP_DATA")
	file.relativePath = string(pathBytes)
	return nil
}

func (file *batchFile) Filename() string {
	return filepath.Base(file.relativePath)
}

func (file *batchFile) RelativePath() string {
	return file.relativePath
}

// AbsolutePath retourne le chemin vers le fichier, sans le schéma
// (base path)
func (file *batchFile) AbsolutePath() string {
	return filepath.Join(file.basepath, file.RelativePath())
}

// IsCompressed est vrai si le fichier est compressé
func (file *batchFile) IsCompressed() bool {
	return strings.HasSuffix(file.RelativePath(), ".gz")
}

// NewBatchFileFromBatch is a helper function that allows to give a batch and filename instead of an
// actual relative path. The batch key is used as subdirectory.
func NewBatchFileFromBatch(basepath string, batch BatchKey, filename string) BatchFile {
	return &batchFile{basepath, path.Join(batch.String(), filename)}
}

func NewBatchFile(basepath string, relativepath string) BatchFile {
	return &batchFile{basepath, relativepath}
}

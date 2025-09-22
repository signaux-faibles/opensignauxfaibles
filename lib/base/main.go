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

func (files *BatchFiles) UnmarshalJSON(data []byte) error {
	var temp map[ParserType][]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*files = make(BatchFiles)

	for parserType, paths := range temp {
		var parserFiles []BatchFile
		for _, path := range paths {
			parserFiles = append(parserFiles, NewBatchFile(path))
		}

		(*files)[parserType] = parserFiles
	}
	return nil
}

// BatchFile is the relevant metadata of a file that needs to be imported
type BatchFile interface {
	Filename() string

	// Relative path given a base path
	// Usually begins with the batch key
	Path() string

	IsCompressed() bool
}

type batchFile struct {
	path string

	// Explicitely specify if the file is compressed
	// If set to nil, then the ".gz" suffix of the filename will determine if
	// the file is compressed.
	isCompressed *bool
}

func (file *batchFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(file.Path())
}

func (file *batchFile) UnmarshalJSON(pathBytes []byte) error {
	file.path = path.Join(viper.GetString("APP_DATA"), string(pathBytes))
	return nil
}

func (file *batchFile) Filename() string {
	return filepath.Base(file.Path())
}

func (file *batchFile) Path() string {
	return file.path
}

// IsCompressed est vrai si le fichier est compressé
func (file *batchFile) IsCompressed() bool {
	if file.isCompressed != nil {
		return *file.isCompressed
	}

	return strings.HasSuffix(file.Path(), ".gz")
}

//--------- BatchFile creation helpers --------

// NewBatchFile create a BatchFile from a given absolute or relative file.
// If several path segments are given, they are concatenated with path.Join.
func NewBatchFile(pathSegments ...string) BatchFile {
	return &batchFile{path.Join(pathSegments...), nil}
}

func NewCompressedBatchFile(path string) BatchFile {
	True := true
	return &batchFile{path, &True}
}

// NewBatchFileFromBatch is a helper function that allows to give a batch and filename instead of an
// actual relative path. The batch key is used as subdirectory.
func NewBatchFileFromBatch(basepath string, batch BatchKey, filename string) BatchFile {
	return NewBatchFile(basepath, batch.String(), filename)
}

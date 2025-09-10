package prepareimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// PopulateFilesProperty populates the "files" property of an Admin object, given a path.
func PopulateFilesProperty(pathname string, batchKey BatchKey) (FilesProperty, []string) {
	batchPath := path.Join(pathname, batchKey.String())
	filenames, _ := ReadFilenames(batchPath)
	var augmentedFiles []DataFile
	for _, file := range filenames {
		augmentedFiles = append(augmentedFiles, NewDataFile(file, batchPath))
	}
	return PopulateFilesPropertyFromDataFiles(augmentedFiles, batchKey)
}

// PopulateFilesPropertyFromDataFiles populates the "files" property of an Admin object, given a list of Data files.
func PopulateFilesPropertyFromDataFiles(filenames []DataFile, batchKey BatchKey) (FilesProperty, []string) {
	filesProperty := FilesProperty{}
	unsupportedFiles := []string{}
	for _, filename := range filenames {
		filetype := filename.DetectFileType()
		if filetype == "" {
			unsupportedFiles = append(unsupportedFiles, batchKey.Path()+filename.GetFilename())
			continue
		}
		if _, exists := filesProperty[filetype]; !exists {
			filesProperty[filetype] = []BatchFile{}
		}
		batchFileToAdd := newBatchFile(batchKey, filename.GetFilename())
		if strings.HasSuffix(filename.GetOriginalFilename(), ".gz") {
			size := filename.GetSize()
			if size == nil {
				panic(errors.New("file size could not be found for: " + batchFileToAdd.Name()))
			}
			batchFileToAdd.AddGzippedSize(*size)
		}
		filesProperty[filetype] = append(filesProperty[filetype], batchFileToAdd)
	}
	return filesProperty, unsupportedFiles
}

// ReadFilenames returns the name of files found at the provided path.
func ReadFilenames(path string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

// FilesProperty represents the "files" property of an Admin object.
type FilesProperty map[ValidFileType][]BatchFile

// GetFilterFile returns the filter file.
func (fp FilesProperty) GetFilterFile() (BatchFile, error) {
	if fp["filter"] == nil || len(fp["filter"]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 filter file, found: %s", fp["filter"])
	}
	return fp["filter"][0], nil
}

func (fp FilesProperty) GetSireneULFile() (BatchFile, error) {
	if fp[sireneUl] == nil || len(fp[sireneUl]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 sireneUL filter file, found %s", fp[sireneUl])
	}
	return fp[sireneUl][0], nil
}

// GetEffectifFile returns the effectif file.
func (fp FilesProperty) GetEffectifFile() (BatchFile, error) {
	if fp["effectif"] == nil || len(fp["effectif"]) != 1 {
		return nil, fmt.Errorf("batch requires just 1 effectif file, found: %s", fp["effectif"])
	}
	return fp["effectif"][0], nil
}

// BatchFile represents a file that is listed in a FilesProperty entry.
type BatchFile interface {
	BatchKey() BatchKey
	Name() string
	Path() string
	AbsolutePath(parentDir string) string
	GetGzippedSize() uint64     // in bytes
	AddGzippedSize(size uint64) // in bytes
}

func newBatchFile(batchKey BatchKey, filename string) BatchFile {
	return &batchFile{
		batchKey: batchKey,
		filename: filename,
	}
}

type batchFile struct {
	batchKey    BatchKey
	filename    string
	gzippedSize uint64 // in bytes
}

func (file *batchFile) BatchKey() BatchKey {
	return file.batchKey
}

func (file *batchFile) Name() string {
	return file.filename
}

// Path retourne le chemin relatif du fichier, avec un préfixe "gzip:" si celui-ci est compressé.
func (file *batchFile) Path() string {
	return file.AbsolutePath("")
}

// AbsolutePath retourne le chemin absolu du fichier, avec un préfixe "gzip:" si celui-ci est compressé.
func (file *batchFile) AbsolutePath(parentDir string) string {
	var prefix string
	if file.gzippedSize > 0 {
		prefix = "gzip:"
	}
	return prefix + path.Join(parentDir, file.batchKey.Path(), file.filename)
}

func (file *batchFile) AddGzippedSize(size uint64) {
	file.gzippedSize = size
}

func (file *batchFile) GetGzippedSize() uint64 {
	return file.gzippedSize
}

// MarshalJSON will be called when serializing the AdminObject.
func (file *batchFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(file.Path())
}

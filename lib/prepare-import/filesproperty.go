package prepareimport

import (
	"opensignauxfaibles/lib/engine"
	"os"
	"path"
)

// PopulateFilesProperty populates the "files" property of an Admin object, given a path.
func PopulateFilesProperty(pathname string, batchKey engine.BatchKey) (engine.BatchFiles, []string) {
	batchPath := path.Join(pathname, batchKey.String())
	filenames, _ := ReadFilenames(batchPath)

	var batchFiles []engine.BatchFile
	for _, file := range filenames {
		batchFiles = append(batchFiles, engine.NewBatchFileFromBatch(pathname, batchKey, file))
	}
	return PopulateFilesPropertyFromDataFiles(batchFiles)
}

// PopulateFilesPropertyFromDataFiles populates the "files" property of an Admin object, given a list of Data files.
func PopulateFilesPropertyFromDataFiles(files []engine.BatchFile) (engine.BatchFiles, []string) {
	filesProperty := engine.BatchFiles{}

	unsupportedFiles := []string{}

	for _, file := range files {
		parserType := ExtractParserTypeFromFilename(file.Filename())

		if parserType == "" {
			unsupportedFiles = append(unsupportedFiles, file.Path())
			continue
		}
		if _, exists := filesProperty[parserType]; !exists {
			filesProperty[parserType] = []engine.BatchFile{}
		}

		filesProperty[parserType] = append(filesProperty[parserType], file)
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

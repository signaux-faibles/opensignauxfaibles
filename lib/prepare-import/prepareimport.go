// Package prepareimport deals with all operations that need to be
// performed before the import runs, e.g. defining exactly which files will be
// imported and their type.
package prepareimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/filter"
)

// PrepareImport generates an Admin object from files found at given pathname,
// in the "batchKey" directory, on the file system.
// func PrepareImport(basepath string, batchKey engine.BatchKey, filterWriter engine.FilterWriter) (engine.AdminBatch, error) {
func PrepareImport(basepath string, batchKey engine.BatchKey, w engine.FilterWriter) (engine.AdminBatch, error) {

	slog.Debug(string("Listing data files in " + batchKey + "/ ..."))

	batchPath := path.Join(basepath, batchKey.String())

	if _, err := os.ReadDir(batchPath); err != nil {
		return engine.AdminBatch{}, fmt.Errorf("could not find directory %s in provided path", batchKey.String())
	}

	var err error
	batchFiles, unsupportedFiles := PopulateFilesProperty(basepath, batchKey)

	// To complete the FilesProperty, we need:
	// - a filter file (created from an effectif file, at the batch/parent level)

	effectifFile, _ := batchFiles.GetEffectifFile()
	filterFile, _ := batchFiles.GetFilterFile()
	sireneULFile, _ := batchFiles.GetSireneULFile()

	if effectifFile != nil {
		slog.Debug("Found effectif file: " + effectifFile.Path())
	}

	if filterFile != nil {
		slog.Debug("Found filter file: " + filterFile.Path())
	}

	if sireneULFile != nil {
		slog.Debug("Found sireneUL file: " + sireneULFile.Path())
	}

	if filterFile == nil && effectifFile == nil {
		return engine.AdminBatch{}, errors.New("filter is missing: batch should include a filter or one effectif file")
	}

	// if needed, create a filter file from the effectif file
	if filterFile == nil {
		slog.Debug("Writing filter file")
		if err = createFilterFromEffectifAndSirene(
			w,
			effectifFile.Path(),
			sireneULFile.Path(),
		); err != nil {
			return engine.AdminBatch{}, err
		}
	}

	// add the filter to filesProperty
	if batchFiles["filter"] == nil && filterFile != nil {
		slog.Debug("Adding filter file to batch ...")
		batchFiles[engine.Filter] = append(batchFiles[engine.Filter], filterFile)
	}

	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}

	return engine.AdminBatch{
		Key:    batchKey,
		Files:  batchFiles,
		Params: populateParamProperty(batchKey),
	}, err
}

func createFilterFromEffectifAndSirene(filterWriter engine.FilterWriter, effectifFilePath string, sireneULFilePath string) error {
	categoriesJuridiqueFilter := filter.CategorieJuridiqueFilter(sireneULFilePath)

	// Create the filter
	sirenFilter, err := filter.Create(
		effectifFilePath, // input: the effectif file
		filter.DefaultNbMois,
		filter.DefaultMinEffectif,
		filter.DefaultNbIgnoredCols,
		categoriesJuridiqueFilter,
	)
	if err != nil {
		return err
	}

	// Write the filter
	return filterWriter.Write(sirenFilter)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// InferBatchProvider infers the batch imports given the filenames
// Implements BatchProvider interface
type InferBatchProvider struct {
	Path     string
	BatchKey engine.BatchKey
}

func (p InferBatchProvider) Get() (engine.AdminBatch, error) {
	var batch engine.AdminBatch

	// TODO only temp
	var w = &filter.MemoryFilterWriter{}
	//

	batch, err := PrepareImport(p.Path, p.BatchKey, w)
	if _, ok := err.(UnsupportedFilesError); ok {
		slog.Warn(fmt.Sprintf("Des fichiers non-identifiés sont présents : %v", err))
	} else if err != nil {
		return batch, fmt.Errorf("une erreur est survenue en préparant l'import : %w", err)
	}

	slog.Info("Batch inféré avec succès")

	batchJSON, _ := json.MarshalIndent(batch, "", "  ")
	if batchJSON != nil {
		slog.Info(string(batchJSON))
	}
	return batch, nil
}

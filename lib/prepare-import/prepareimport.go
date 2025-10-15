package prepareimport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname,
// in the "batchKey" directory, on the file system.
func PrepareImport(basepath string, batchKey engine.BatchKey) (engine.AdminBatch, error) {

	fmt.Println("Listing data files in " + batchKey + "/ ...")

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
		fmt.Println("Found effectif file: " + effectifFile.Path())
	}

	if filterFile != nil {
		fmt.Println("Found filter file: " + filterFile.Path())
	}

	if sireneULFile != nil {
		fmt.Println("Found sireneUL file: " + sireneULFile.Path())
	}

	if filterFile == nil && effectifFile == nil {
		return engine.AdminBatch{}, errors.New("filter is missing: batch should include a filter or one effectif file")
	}

	// if needed, create a filter file from the effectif file
	if filterFile == nil {
		filterFile = engine.NewBatchFileFromBatch(basepath, batchKey, "filter_siren.csv")

		fmt.Println("Generating filter file: " + filterFile.Path() + " ...")
		if err = createFilterFromEffectifAndSirene(
			filterFile.Path(),
			effectifFile.Path(),
			sireneULFile.Path(),
		); err != nil {
			return engine.AdminBatch{}, err
		}
	}

	// add the filter to filesProperty
	if batchFiles["filter"] == nil && filterFile != nil {
		fmt.Println("Adding filter file to batch ...")
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

func createFilterFromEffectifAndSirene(filterFilePath string, effectifFilePath string, sireneULFilePath string) error {
	if fileExists(filterFilePath) {
		return errors.New("about to overwrite existing filter file: " + filterFilePath)
	}
	filterWriter, err := os.Create(filterFilePath)
	if err != nil {
		return err
	}
	categoriesJuridiqueFilter := createfilter.CategorieJuridiqueFilter(sireneULFilePath)

	return createfilter.CreateFilter(
		filterWriter,     // output: the filter file
		effectifFilePath, // input: the effectif file
		createfilter.DefaultNbMois,
		createfilter.DefaultMinEffectif,
		createfilter.DefaultNbIgnoredCols,
		categoriesJuridiqueFilter,
	)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes. Source: https://stackoverflow.com/a/21061062/592254
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// InferBatchProvider infers the batch imports given the filenames
// Implements BatchProvider interface
type InferBatchProvider struct {
	Path     string
	BatchKey engine.BatchKey
}

func (p InferBatchProvider) Get() (engine.AdminBatch, error) {
	var batch engine.AdminBatch
	batch, err := PrepareImport(p.Path, p.BatchKey)

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

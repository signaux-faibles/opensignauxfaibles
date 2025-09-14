package prepareimport

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"opensignauxfaibles/lib/prepare-import/createfilter"
)

// PrepareImport generates an Admin object from files found at given pathname of the file system.
func PrepareImport(pathname string, batchKey BatchKey, providedDateFinEffectif string) (AdminBatch, error) {

	batchPath := getBatchPath(pathname, batchKey)
	println("Listing data files in " + batchPath + "/ ...")
	if _, err := os.ReadDir(path.Join(pathname, batchPath)); err != nil {
		return AdminBatch{}, fmt.Errorf("could not find directory %s in provided path", batchPath)
	}

	var err error
	filesProperty, unsupportedFiles := PopulateFilesProperty(pathname, batchKey)

	// To complete the FilesProperty, we need:
	// - a filter file (created from an effectif file, at the batch/parent level)
	// - a dateFinEffectif value (provided as parameter, or detected from effectif file)

	var dateFinEffectif time.Time
	effectifFile, _ := filesProperty.GetEffectifFile()
	filterFile, _ := filesProperty.GetFilterFile()
	sireneULFile, _ := filesProperty.GetSireneULFile()
	if (effectifFile == nil || filterFile == nil) && batchKey.IsSubBatch() {
		println("Looking for effectif and/or filter file in " + batchKey.GetParentBatch() + " ...")
		parentFilesProperty, _ := PopulateFilesProperty(pathname, newSafeBatchKey(batchKey.GetParentBatch()))
		if effectifFile == nil {
			effectifFile, _ = parentFilesProperty.GetEffectifFile()
		}
		if filterFile == nil {
			filterFile, _ = parentFilesProperty.GetFilterFile()
		}
	}

	if effectifFile != nil {
		println("Found effectif file: " + effectifFile.Name())
	}

	if filterFile != nil {
		println("Found filter file: " + filterFile.Name())
	}

	if sireneULFile != nil {
		println("Found sireneUL file: " + sireneULFile.Name())
	}

	// if needed, create a filter file from the effectif file
	if filterFile == nil {
		if effectifFile == nil {
			return AdminBatch{}, errors.New("filter is missing: batch should include a filter or one effectif file")
		}
		effectifFilePath := effectifFile.AbsolutePath(pathname)
		sireneULFilePath := sireneULFile.AbsolutePath(pathname)
		effectifBatch := effectifFile.BatchKey()
		filterFile = newBatchFile(effectifBatch, "filter_siren_"+effectifBatch.String()+".csv")
		println("Generating filter file: " + filterFile.Path() + " ...")
		if err = createFilterFromEffectifAndSirene(path.Join(pathname, filterFile.Path()), effectifFilePath, sireneULFilePath); err != nil {
			return AdminBatch{}, err
		}
	}

	// add the filter to filesProperty
	if filesProperty["filter"] == nil && filterFile != nil {
		if batchKey.IsSubBatch() {
			// copy the filter into the sub-batch's directory
			println("Copying filter file to " + filterFile.Path() + " ...")
			src := path.Join(pathname, filterFile.Path())
			dest := path.Join(pathname, batchKey.GetParentBatch(), batchKey.Path(), filterFile.Name())
			err = copy(src, dest)
			if err != nil {
				return AdminBatch{}, err
			}
			filterFile = newBatchFile(batchKey, filterFile.Name())
		}
		println("Adding filter file to batch ...")
		filesProperty["filter"] = append(filesProperty["filter"], filterFile)
	}

	if effectifFile != nil {
		println("Detecting dateFinEffectif from effectif file ...")
		effectifFilePath := effectifFile.AbsolutePath(pathname)
		dateFinEffectif, err = createfilter.DetectDateFinEffectif(effectifFilePath, createfilter.DefaultNbIgnoredCols) // TODO: éviter de lire le fichier Effectif deux fois
		if err != nil {
			return AdminBatch{}, err
		}
	}

	// make sure we have date_fin_effectif
	if dateFinEffectif.IsZero() {
		println("Still missing date_fin_effectif => parsing CLI parameter ...")
		dateFinEffectif, err = time.Parse("2006-01-02", providedDateFinEffectif)
		if err != nil {
			return AdminBatch{}, errors.New("date_fin_effectif is missing or invalid: " + providedDateFinEffectif)
		}
	}

	if len(unsupportedFiles) > 0 {
		err = UnsupportedFilesError{unsupportedFiles}
	}

	return AdminBatch{
		ID:            IDProperty{batchKey, "batch"},
		Files:         populateFilesPaths(filesProperty),
		CompleteTypes: populateCompleteTypesProperty(filesProperty),
		Param:         populateParamProperty(batchKey, NewDateFinEffectif(dateFinEffectif)),
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

func getBatchPath(pathname string, batchKey BatchKey) string {
	if batchKey.IsSubBatch() {
		return path.Join(batchKey.GetParentBatch(), batchKey.String())
	}
	return batchKey.String()
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

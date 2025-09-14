package prepareimport

import (
	"fmt"
	"opensignauxfaibles/lib/base"
	"strconv"
	"strings"
	"time"
)

// AdminBatch represents a document going to be stored in the Admin db collection.
type AdminBatch struct {
	ID            base.AdminID                    `json:"id,omitempty"`
	CompleteTypes []base.ValidFileType            `json:"complete_types,omitempty"`
	Files         map[base.ValidFileType][]string `json:"files,omitempty"`
	Param         ParamProperty                   `json:"param,omitempty"`
}

// ParamProperty represents the "param" property of an Admin object.
type ParamProperty struct {
	DateDebut       time.Time `json:"date_debut"`
	DateFin         time.Time `json:"date_fin"`
	DateFinEffectif time.Time `json:"date_fin_effectif"`
}

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "type de fichier non supporté : " + strings.Join(err.UnsupportedFiles, ", ")
}

func populateParamProperty(batchKey BatchKey, dateFinEffectif DateFinEffectif) ParamProperty {
	year, _ := strconv.Atoi("20" + batchKey.String()[0:2])
	month, _ := strconv.Atoi(batchKey.String()[2:4])
	return ParamProperty{
		DateDebut:       time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		DateFin:         time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC),
		DateFinEffectif: dateFinEffectif.Date(),
	}
}

func populateCompleteTypesProperty(filesProperty FilesProperty) []base.ValidFileType {
	completeTypes := []base.ValidFileType{}
	for _, typeName := range defaultCompleteTypes {
		if _, ok := filesProperty[typeName]; ok {
			completeTypes = append(completeTypes, typeName)
		}
	}
	for typeName, thresholdInBytes := range thresholdPerGzippedFileType {
		if files, ok := filesProperty[typeName]; ok {
			if len(files) != 1 {
				panic(fmt.Errorf("'complete' file detection can only work if there is only 1 file per type, found %v for type %v", len(files), typeName))
			}
			file := files[0]
			if file.GetGzippedSize() >= thresholdInBytes {
				println(fmt.Sprintf("Info: file \"%v\" was marked as \"complete\" because it's a gzipped file which size reached the threshold of %v bytes", file.Name(), thresholdInBytes))
				completeTypes = append(completeTypes, typeName)
			}
		}
	}
	return completeTypes
}

func populateFilesPaths(filesProperty FilesProperty) map[base.ValidFileType][]string {
	r := make(map[base.ValidFileType][]string)
	for k, v := range filesProperty {

		var paths []string
		for _, batchFile := range v {
			paths = append(paths, batchFile.Path())
		}

		r[k] = paths
	}
	return r
}

// types of files that are always provided as "complete"
var defaultCompleteTypes = []base.ValidFileType{
	base.Apconso,
	base.Apdemande,
	base.Effectif,
	base.EffectifEnt,
	base.Sirene,
	base.SireneUl,
}

// types of files that will be considered as "complete" if their gzipped size reach a certain threshold (in bytes)
var thresholdPerGzippedFileType = map[base.ValidFileType]uint64{
	base.Cotisation: 143813078,
	base.Delai:      1666199,
	base.Procol:     1646193,
	base.Debit:      254781489,
}

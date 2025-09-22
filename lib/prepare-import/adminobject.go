package prepareimport

import (
	"opensignauxfaibles/lib/base"
	"strconv"
	"strings"
	"time"
)

// UnsupportedFilesError is an Error object that lists files that were not supported.
type UnsupportedFilesError struct {
	UnsupportedFiles []string
}

func (err UnsupportedFilesError) Error() string {
	return "type de fichier non supporté : " + strings.Join(err.UnsupportedFiles, ", ")
}

func populateParamProperty(batchKey base.BatchKey) base.AdminBatchParams {
	year, _ := strconv.Atoi("20" + batchKey.String()[0:2])
	month, _ := strconv.Atoi(batchKey.String()[2:4])
	return base.AdminBatchParams{
		DateDebut: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		DateFin:   time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC),
	}
}

package prepareimport

import (
	"opensignauxfaibles/lib/engine"
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

func populateParamProperty(batchKey engine.BatchKey) engine.AdminBatchParams {
	year, _ := strconv.Atoi("20" + batchKey.String()[0:2])
	month, _ := strconv.Atoi(batchKey.String()[2:4])
	dateFin := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return engine.AdminBatchParams{
		DateDebut: dateFin.AddDate(-10, 0, 0), // 10 ans avant DateFin
		DateFin:   dateFin,
	}
}

package altares

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

var FIELDS = []int{
	1,
	18,
	21,
	21,
	22,
	23,
	24,
	26,
	27,
	30,
}

var END_OF_FILE_REGEXP = regexp.MustCompile("Fin du fichier : total (?P<nblines>\\d+) ligne\\(s\\)")

func ConvertIncrement(incrementFilename string, output io.Writer) {
	slog.Info("démarrage de la conversion du fichier incrémental", slog.String("filename", incrementFilename))
	inputFile, err := os.Open(incrementFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			panic(errors.Wrap(closeErr, "erreur à la fermeture du fichier"))
		}
	}()

	fromISO8859_15toUTF8 := charmap.ISO8859_15.NewDecoder()
	convertReader := fromISO8859_15toUTF8.Reader(inputFile)
	reader := csv.NewReader(convertReader)
	reader.Comma = ';'

	w := csv.NewWriter(output)
	defer w.Flush()

	// discard headers
	headers, err := reader.Read()
	utils.ManageError(err, "erreur lors de la lecture des headers")
	slog.Debug("description des headers", slog.Any("headers", headers))

	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok && err.Err != csv.ErrFieldCount {
			slog.Warn("prbleme", slog.Any("error", err))
			continue
		}
		//if err == io.EOF {
		//	return
		//}
		if isIncrementalEndOfFile(record) {
			return
		}
		out := selectFields(record)

		if out != nil {
			err = w.Write(out)
			if err != nil {
				panic(err)
			}
		}
	}
}

func selectFields(record []string) []string {
	var data []string
	for _, field := range FIELDS {
		//if field > len(record)-1 {
		//	if
		//	slog.Error(
		//		"erreur de longueur de ligne",
		//		slog.Any("record", record),
		//	)
		//	utils.ManageError(fmt.Errorf("erreur de longueur de ligne, attendu : %d, obtenu %d", field, len(record)-1), "mais que se passe t il ?"))
		//	return nil
		//}
		data = append(data, record[field])
	}
	return data
}

func isIncrementalEndOfFile(record []string) bool {
	if len(record) > 2 {
		return false
	}
	//names := END_OF_FILE_REGEXP.SubexpNames()
	//slog.Info("les captures", slog.Any("names", names))
	line := END_OF_FILE_REGEXP.FindStringSubmatch(record[0])
	if len(line) != 2 {
		utils.ManageError(fmt.Errorf("erreur de fin de fichier : %+v", record), "problème avec la fin de fichier")
	}
	slog.Info("fin du fichier incrémental", slog.Any("expectedLines", line[1]))
	return true
}

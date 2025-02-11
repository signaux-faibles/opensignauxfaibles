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

var mappingMensuel = mapping{
	siren:                           simpleConversion(1),
	etat_organisation:               simpleConversion(18),
	code_paydex:                     simpleConversion(21),
	nbr_jrs_retard:                  simpleConversion(22),
	nbr_fournisseurs:                simpleConversion(23),
	encours_etudies:                 simpleConversion(24),
	note_100_alerteur_plus_30:       simpleConversion(26),
	note_100_alerteur_plus_90_jours: simpleConversion(27),
	date_valeur:                     simpleConversion(30),
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

	readAllRowsUntil(reader, w, mappingMensuel, true, isIncrementalEndOfFile)
}

func isIncrementalEndOfFile(record []string) bool {
	if len(record) > 2 {
		return false
	}
	line := END_OF_FILE_REGEXP.FindStringSubmatch(record[0])
	if len(line) != 2 {
		utils.ManageError(fmt.Errorf("erreur de fin de fichier : %+v", record), "problème avec la fin de fichier")
	}
	slog.Info("fin du fichier incrémental", slog.Any("expectedLines", line[1]))
	return true
}

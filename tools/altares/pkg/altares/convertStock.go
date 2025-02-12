package altares

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

var mappingStock = mapping{
	siren:             simpleConversion(0),
	etat_organisation: simpleConversion(1),
	code_paydex:       simpleConversion(2),
	nbr_jrs_retard:    simpleConversion(3),
	nbr_fournisseurs:  simpleConversion(4),
	//encours_etudies:   simpleConversion(6),
	encours_etudies:                 advancedConversion(5, useDotAsDecimalDelimiter),
	note_100_alerteur_plus_30:       simpleConversion(6),
	note_100_alerteur_plus_90_jours: simpleConversion(7),
	date_valeur:                     advancedConversion(8, datifyStock),
}

func ConvertStock(stockFilename string, output io.Writer) {
	slog.Info("conversion du fichier stock", slog.String("filename", stockFilename))
	inputFile, err := os.Open(stockFilename)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := inputFile.Close()
		if closeErr != nil {
			panic(errors.Wrap(closeErr, "erreur à la fermeture du fichier"))
		}
	}()
	//inputFileWithoutBOM, encoding := utfbom.Skip(inputFile)
	//slog.Info("encodage du fichier stock détecté", slog.String("encoding", encoding.String()))
	reader := csv.NewReader(inputFile)
	reader.TrimLeadingSpace = true
	reader.Comma = ','

	w := csv.NewWriter(output)
	defer w.Flush()

	readAllRows(reader, w, mappingStock, true)
}

func datifyStock(s string) string {
	return s[6:10] + "-" + s[3:5] + "-" + s[0:2]
}
func useDotAsDecimalDelimiter(s string) string {
	return s[6:10] + "-" + s[3:5] + "-" + s[0:2]
}

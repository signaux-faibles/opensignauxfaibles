package altares

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

const STOCK_FIRST_LINE = "SIREN;ETAT_ORGANISATION;CODE_PAYDEX;PAYDEX;NBR_JRS_RETARD;NBR_FOURNISSEURS;ENCOURS_ETUDIES;NBR_EXPERIENCES_PAIEMENT;NOTE100_ALERTEUR_PLUS_30;NOTE100_ALERTEUR_PLUS_90_JOURS;DATE_VALEUR"

var mappingStock = mapping{
	siren:                           simpleConversion(0),
	etat_organisation:               simpleConversion(1),
	code_paydex:                     simpleConversion(2),
	nbr_jrs_retard:                  simpleConversion(4),
	nbr_fournisseurs:                simpleConversion(5),
	encours_etudies:                 simpleConversion(6),
	note_100_alerteur_plus_30:       simpleConversion(8),
	note_100_alerteur_plus_90_jours: simpleConversion(9),
	date_valeur:                     advancedConversion(10, datifyStock),
}

func ConvertStock(stockFilename string, output io.Writer) {
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
	convertStockFile(inputFile, output)
}

func convertStockFile(inputFile *os.File, output io.Writer) {
	slog.Info("conversion du fichier stock", slog.String("filename", inputFile.Name()))

	reader := csv.NewReader(inputFile)
	reader.TrimLeadingSpace = true
	reader.Comma = ';'

	w := csv.NewWriter(output)
	defer w.Flush()

	readAllRows(reader, w, mappingStock, true)
}

func datifyStock(s string) string {
	return s[6:10] + "-" + s[3:5] + "-" + s[0:2]
}

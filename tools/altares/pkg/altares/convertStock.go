package altares

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

const HEADER_TO_REMOVE = "NBR_EXPERIENCES_PAIEMENT"

var FORMATTERS = map[int]func(string) string{
	9: datifyStock,
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
	reader := csv.NewReader(inputFile)
	reader.Comma = ';'

	w := csv.NewWriter(output)
	defer w.Flush()

	slog.Debug("manipule les headers")
	columnToRemove := manageHeaders(reader, w)
	slog.Debug("suppression d'une colonne", slog.String("status", "start"), slog.Int("index", columnToRemove))

	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok {
			switch err.Err {
			case csv.ErrFieldCount:
				slog.Warn(
					"erreur lors de la lecture du fichier stock, enregistrement rejeté",
					slog.Any("error", err.Err),
					slog.Any("line", reader.InputOffset()),
					slog.Any("record", record),
				)
				continue
			default:
				slog.Error("erreur lors de la lecture", slog.Any("error", err))
				utils.ManageError(err, "erreur pendant la suppression de colonne")
			}
		}
		if err == io.EOF {
			slog.Info(
				"suppression d'une colonne",
				slog.String("status", "end"),
				slog.Int("index", columnToRemove),
				slog.Any("row", reader.InputOffset()),
			)
			return
		}
		newRecord := removeColumn(record, columnToRemove)
		formatValues(newRecord)
		err = w.Write(newRecord)
		utils.ManageError(err, "erreur pendant l'écriture du fichier de sortie")
	}
}

func formatValues(record []string) {
	for idx, value := range record {
		formatter := FORMATTERS[idx]
		if formatter != nil {
			record[idx] = formatter(value)
		}
	}
}

func datifyStock(s string) string {
	return s[6:10] + "-" + s[3:5] + "-" + s[0:2]
}

func manageHeaders(reader *csv.Reader, w *csv.Writer) int {
	columnToRemove := -1
	headers, err := reader.Read()
	utils.ManageError(err, "erreur à la lecture des headers du fichier stock")
	slog.Debug("description des headers", slog.Any("headers", headers))
	for idx, header := range headers {
		if header == HEADER_TO_REMOVE {
			columnToRemove = idx
		}
		headers[idx] = formatHeader(header)
	}
	err = w.Write(removeColumn(headers, columnToRemove))
	if err != nil {
		panic(err)
	}
	slog.Info("colonne à supprimer", slog.Int("idx", columnToRemove))
	return columnToRemove
}

func removeColumn(record []string, remove int) []string {
	var r []string
	for idx, value := range record {
		if idx != remove {
			r = append(r, value)
		}
	}
	return r
}

func formatHeader(input string) string {
	r := strcase.ToSnakeWithIgnore(input, ".,")
	r = strings.TrimSpace(r)
	r = strings.Replace(r, "etat", "état", 1)
	r = strings.Replace(r, "etudies", "étudiés", 1)
	return r
}

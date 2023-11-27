package altares

import (
	"encoding/csv"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

var columnsToRemove = []string{"NBR_EXPERIENCES_PAIEMENT"}

var formatters = map[int]func(string) string{
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
	inputFileWithoutBOM, encoding := utfbom.Skip(inputFile)
	slog.Info("encodage du fichier stock détecté", slog.String("encoding", encoding.String()))
	reader := csv.NewReader(inputFileWithoutBOM)
	reader.TrimLeadingSpace = true
	reader.Comma = ';'

	w := csv.NewWriter(output)
	defer w.Flush()

	slog.Debug("manipule les headers")
	columnsToRemove := manageHeaders(reader, w)
	slog.Debug("conversion du fichier stock", slog.String("status", "start"))

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
			slog.Debug(
				"conversion du fichier stock",
				slog.String("status", "end"),
				slog.Any("row", reader.InputOffset()),
			)
			return
		}
		newRecord := removeColumns(record, columnsToRemove...)
		formatValues(newRecord)
		err = w.Write(newRecord)
		utils.ManageError(err, "erreur pendant l'écriture du fichier de sortie")
	}
}

func formatValues(record []string) {
	for idx, value := range record {
		formatter := formatters[idx]
		if formatter != nil {
			record[idx] = formatter(value)
		}
	}
}

func datifyStock(s string) string {
	return s[6:10] + "-" + s[3:5] + "-" + s[0:2]
}

func manageHeaders(reader *csv.Reader, w *csv.Writer) []int {
	var idxToRemove []int
	headers, err := reader.Read()
	utils.ManageError(err, "erreur à la lecture des headers du fichier stock")
	slog.Debug("description des headers", slog.Any("headers", headers))
	for idx, header := range headers {
		if slices.Contains(columnsToRemove, header) {
			idxToRemove = append(idxToRemove, idx)
		}
		headers[idx] = formatHeader(header)
	}
	err = w.Write(removeColumns(headers, idxToRemove...))
	if err != nil {
		panic(err)
	}
	slog.Info("colonnes à supprimer", slog.Any("idx", idxToRemove))
	return idxToRemove
}

func removeColumns(record []string, remove ...int) []string {
	if remove == nil {
		return record
	}
	var r []string
	for idx, value := range record {
		if !slices.Contains(remove, idx) {
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

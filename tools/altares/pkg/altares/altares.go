package altares

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/pkg/errors"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

type column int

type mapping map[column]conversion

type conversion struct {
	source  int
	convert func(string) string
}

const RESULT_FIRST_LINE = "siren,état_organisation,code_paydex,nbr_jrs_retard,nbr_fournisseurs,encours_étudiés,note_100_alerteur_plus_30,note_100_alerteur_plus_90_jours,date_valeur"

const (
	siren column = iota
	etat_organisation
	code_paydex
	nbr_jrs_retard
	nbr_fournisseurs
	encours_etudies
	note_100_alerteur_plus_30
	note_100_alerteur_plus_90_jours
	date_valeur
)

var allColumns = []column{siren, etat_organisation, code_paydex, nbr_jrs_retard, nbr_fournisseurs, encours_etudies, note_100_alerteur_plus_30, note_100_alerteur_plus_90_jours, date_valeur}

func (c column) String() string {
	switch c {
	case siren:
		return "siren"
	case etat_organisation:
		return "état_organisation"
	case code_paydex:
		return "code_paydex"
	case nbr_jrs_retard:
		return "nbr_jrs_retard"
	case nbr_fournisseurs:
		return "nbr_fournisseurs"
	case encours_etudies:
		return "encours_étudiés"
	case note_100_alerteur_plus_30:
		return "note_100_alerteur_plus_30"
	case note_100_alerteur_plus_90_jours:
		return "note_100_alerteur_plus_90_jours"
	case date_valeur:
		return "date_valeur"
	}
	utils.ManageError(fmt.Errorf("type de colonne inconnu : %d", c), "erreur très bizarre")
	return ""
}

func simpleConversion(idx int) conversion {
	return conversion{
		source:  idx,
		convert: func(s string) string { return s },
	}
}

func advancedConversion(idx int, f func(s string) string) conversion {
	return conversion{
		source:  idx,
		convert: f,
	}
}

func readAllRows(r *csv.Reader, w *csv.Writer, m mapping, skipHeaders bool) {
	readAllRowsUntil(r, w, m, skipHeaders, nil)
}

func readAllRowsUntil(r *csv.Reader, w *csv.Writer, m mapping, skipHeaders bool, eofDetector func([]string) bool) {
	if skipHeaders {
		// discard headers
		headers, err := r.Read()
		utils.ManageError(err, "erreur lors de la lecture des headers")
		slog.Debug("description des headers", slog.Any("headers", headers))
	}
	for {
		record, err := r.Read()
		if eofDetector != nil && eofDetector(record) {
			slog.Info("fichier terminé")
			return
		}
		if err, ok := err.(*csv.ParseError); ok {
			switch err.Err {
			case csv.ErrFieldCount:
				slog.Warn(
					"erreur lors de la lecture du fichier stock, enregistrement rejeté",
					slog.Any("error", err.Err),
					slog.Any("record", record),
				)
				continue
			default:
				slog.Error("erreur lors de la lecture", slog.Any("error", err))
				utils.ManageError(err, "erreur pendant la suppression de colonne")
			}
		}
		if err == io.EOF {
			slog.Info("fichier terminé")
			return
		}
		out, err := convertRow(record, m)
		utils.ManageError(err, "erreur à la conversion d'une ligne du fichier")
		if out != nil {
			err = w.Write(out)
			utils.ManageError(err, "erreur à l'écriture du fichier converti")
		}
	}
}

func convertRow(record []string, m mapping) ([]string, error) {
	var r []string
	if record == nil {
		slog.Warn("enregistrement nil")
		return nil, nil
	}
	if m == nil {
		slog.Warn("mapping non défini")
		return record, nil
	}
	if len(record) < len(m) {
		slog.Warn("moins de colonnes dans l'enregistrement que dans le conversion")
	}
	for _, colonne := range allColumns {
		actualConversion, found := m[colonne]
		if !found {
			return nil, fmt.Errorf("erreur de définition de mapping pour la colonne %d", colonne)
		}
		if len(record) <= actualConversion.source {
			return nil, fmt.Errorf("problème de conversion de la ligne, l'index du mapping recherché %v est inférieur à la longueur de l'enregistrement %v",
				actualConversion.source, len(record))
		}
		value := record[actualConversion.source]
		r = append(r, actualConversion.convert(value))
	}
	return r, nil
}

func WriteHeaders(output io.Writer) {
	writer := csv.NewWriter(output)
	defer writer.Flush()
	var headers []string
	for _, currentCol := range allColumns {
		headers = append(headers, currentCol.String())
	}
	err := writer.Write(headers)
	utils.ManageError(err, "erreur pendant  l'écriture des headers")
}

func Convert(filename string, output io.Writer) {
	// Ouvrir le fichier en lecture
	slog.Info("démarrage de la conversion du fichier ", slog.String("filename", filename))
	fichier, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := fichier.Close()
		if closeErr != nil {
			panic(errors.Wrap(closeErr, "erreur à la fermeture du fichier"))
		}
	}()

	// Créer un scanner pour lire le fichier ligne par ligne
	scanner := bufio.NewScanner(fichier)

	// Lire la première ligne
	if scanner.Scan() {
		premiereLigne := scanner.Text()
		if IsStockFile(premiereLigne) {
			convertStockFile(fichier, output)
		} else if IsIncrementFile(premiereLigne) {
			convertIncrementFile(fichier, output)
		} else if IsResultFile(premiereLigne) {
			slog.Info("le fichier est le résultat d'un process précédent", slog.String("filename", filename))
		} else {
			utils.ManageError(errors.New("fichier de type inconnu"), "erreur à la conversion du fichier", slog.Any("filename", filename))
		}
	} else if err := scanner.Err(); err != nil {
		utils.ManageError(err, "erreur à la lecture de la première ligne du fichier", slog.Any("filename", filename))
	}
}

func IsStockFile(firstLine string) bool {
	return strings.EqualFold(STOCK_FIRST_LINE, firstLine)
}

func IsIncrementFile(firstLine string) bool {
	return strings.EqualFold(INCREMENT_FIRST_LINE, firstLine)
}

func IsResultFile(firstLine string) bool {
	return strings.EqualFold(RESULT_FIRST_LINE, firstLine)
}

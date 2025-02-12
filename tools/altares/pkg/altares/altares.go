package altares

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

type column int

type mapping map[column]conversion

type conversion struct {
	source  int
	convert func(string) string
}

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
		slog.Warn("moins de colonnes dans l'enregistrement que dans la conversion")
	}
	for _, colonne := range allColumns {
		actualConversion, found := m[colonne]
		if !found {
			return nil, fmt.Errorf("erreur de définition de mapping pour la colonne %d", colonne)
		}
		if len(record) <= actualConversion.source {
			slog.Warn(fmt.Sprintf("erreur sur la ligne %v", record))
			slog.Warn(fmt.Sprintf("erreur sur la ligne %v", m))
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

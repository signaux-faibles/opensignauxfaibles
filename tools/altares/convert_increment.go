package main

import (
	"encoding/csv"
	"io"
	"os"
)

const EMPTY int = -1

var FIELDS = []int{
	1,
	18,
	21,
	21,
	22,
	23,
	24,
	EMPTY,
	26,
	27,
	30,
}

var DATES = []int{
	10,
}

var HEADERS = []string{
	"SIREN",
	"ETAT_ORGANISATION",
	"CODE_PAYDEX",
	"PAYDEX",
	"NBR_JRS_RETARD",
	"NBR_FOURNISSEURS",
	"ENCOURS_ETUDIES",
	"NBR_EXPERIENCES_PAIEMENT",
	"NOTE100_ALERTEUR_PLUS_30",
	"NOTE100_ALERTEUR_PLUS_90_JOURS",
	"DATE_VALEUR",
}

func main() {
	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(inputFile)
	//reader.LazyQuotes = true
	reader.Comma = ';'

	w := csv.NewWriter(os.Stdout)
	w.Comma = reader.Comma

	// discard headers
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	err = w.Write(HEADERS)
	if err != nil {
		panic(err)
	}

	for {
		record, err := reader.Read()
		if err, ok := err.(*csv.ParseError); ok && err.Err != csv.ErrFieldCount {
			continue
		}
		if err == io.EOF {
			w.Flush()
			return
		}
		out := output(record)

		if out != nil {
			err = w.Write(datifyOutput(out))
			if err != nil {
				panic(err)
			}
		}
	}
}

func output(record []string) []string {
	var data []string
	for _, field := range FIELDS {
		if field > len(record)-1 {
			return nil
		}
		if field == EMPTY {
			data = append(data, "")
		} else {
			data = append(data, record[field])
		}
	}
	return data
}

func datifyOutput(output []string) []string {
	for _, rank := range DATES {
		output[rank] = datify(output[rank])
	}
	return output
}

func datify(s string) string {
	return s[8:10] + "/" + s[5:7] + "/" + s[0:4]
}

// Fichier incrément
// 0 Référence Client
// 1 Siren
// 2 Siret
// 3 Raison sociale 1
// 4 Raison sociale 2
// 5 Enseigne
// 6 Sigle
// 7 Complément d'adresse
// 8 Adresse
// 9 Distribution spéciale
// 10 Code postal et bureau distributeur
// 11 Pays
// 12 Code postal
// 13 Ville
// 14 Qualité Etablissement
// 15 Code type d'établissement
// 16 Libellé type d'établissement
// 17 Etat d'activité établissement
// 18 Etat d'activité entreprise
// 19 Etat de procédure collective
// 20 Diffusible
// 21 Paydex
// 22 Retard moyen de paiements (j)
// 23 Nombre de fournisseurs analysés
// 24 Montant total des encours étudiés (¤)
// 25 Montant total des encours échus non réglés (¤)
// 26 FPI 30+
// 27 FPI 90+
// 28 Code du mouvement
// 29 Libellé du mouvement
// 30 Date d'effet du mouvement

// Fichier stock
// SIREN
// ETAT_ORGANISATION
// CODE_PAYDEX
// PAYDEX
// NBR_JRS_RETARD
// NBR_FOURNISSEURS
// ENCOURS_ETUDIES
// NBR_EXPERIENCES_PAIEMENT
// NOTE100_ALERTEUR_PLUS_30
// NOTE100_ALERTEUR_PLUS_90_JOURS
// DATE_VALEUR

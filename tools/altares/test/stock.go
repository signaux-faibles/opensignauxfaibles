package test

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"opensignauxfaibles/tools/altares/pkg/altares"
	"opensignauxfaibles/tools/altares/pkg/utils"
)

func GenerateStockCSV(lines int) *os.File {
	temp, err := os.CreateTemp(os.TempDir(), "stock_*.csv")
	utils.ManageError(err, "erreur à la création du fichier stock")
	headers := []string{
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
	writer := csv.NewWriter(temp)
	writer.Comma = ';'
	err = writer.Write(headers)
	utils.ManageError(err, "erreur à l'écriture des headers")
	for i := 0; i < lines; i++ {
		err = writer.Write(newStockLine())
		utils.ManageError(err, "erreur d'écriture de la ligne")
	}
	writer.Flush()
	err = temp.Close()
	utils.ManageError(err, "erreur à la fermeture du fichier")
	open, err := os.Open(temp.Name())
	utils.ManageError(err, "erreur à la réouverture du fichier")
	return open
}

func newStockLine() []string {
	codePaydex, paydexLabel, nbJours := aPaydex()
	return []string{
		aSiret(),
		anEtatOrganisation(),
		codePaydex,
		paydexLabel,
		nbJours,
		aNbFournisseurs(),
		aEncoursEtudies(),
		Fake.Lorem().Word(), // cette ligne est coupée normalement
		aNote100(),
		aNote100(),
		aStockDateValeur(),
	}
}

func aStockDateValeur() string {
	between := Fake.Time().TimeBetween(
		time.Now().AddDate(-10, 0, 0),
		time.Now().AddDate(0, -1, 0),
	)
	return strings.ReplaceAll(between.Format(time.DateOnly), "/", "-")
}

func aNote100() string {
	return strconv.Itoa(Fake.IntBetween(0, 100))
}

func aEncoursEtudies() string {
	float := Fake.RandomFloat(3, 1, 999999)
	return strconv.FormatFloat(float, 'f', -1, 64)
}

func aNbFournisseurs() string {
	nb := Fake.IntBetween(1, 150)
	return strconv.Itoa(nb)
}

func aPaydex() (string, string, string) {
	code := Fake.RandomStringMapKey(altares.Paydex)
	label := altares.Paydex[code]
	nbJours := ""
	before, _, found := strings.Cut(label, " ")
	if found {
		_, err := strconv.Atoi(before)
		if err == nil {
			nbJours = before
		}
	}
	return code, label, nbJours
}

func anEtatOrganisation() string {
	return Fake.RandomStringElement([]string{"Actif", "Fermé", "Liquidé"})
}

func aSiret() string {
	return fmt.Sprintf("%09d", Fake.IntBetween(1000, 9999))
}

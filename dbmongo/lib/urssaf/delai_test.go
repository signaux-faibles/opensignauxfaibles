package urssaf

import (
	"fmt"
	"testing"
	"time"

	"github.com/signaux-faibles/gournal"
)

func TestDelai(t *testing.T) {
	field := map[string]int{
		"NumeroCompte":      2,
		"NumeroContentieux": 3,
		"DateCreation":      4,
		"DateEcheance":      5,
		"DureeDelai":        6,
		"Denomination":      7,
		"Indic6m":           8,
		"AnneeCreation":     9,
		"MontantEcheancier": 10,
		"Stade":             11,
		"Action":            12,
	}
	test_row_1 := []string{"092", "44", "111111111111111111", "2222222222222", "25/09/2013", "04/10/2016", "444", "Test", "SUP", "2014", "200000.00", "PR", "SUR PO"}
	var tracker gournal.Tracker
	test_1, tracker := readLine(test_row_1, field, "33333333333333", tracker)

	loc, _ := time.LoadLocation("Europe/Paris")
	date1 := time.Date(2013, 9, 25, 0, 0, 0, 0, loc)
	date2 := time.Date(2016, 10, 4, 0, 0, 0, 0, loc)
	fmt.Println(test_1)
	fmt.Println(date1)
	fmt.Println(date2)

	expected_1 := Delai{
		key:               "33333333333333",
		NumeroCompte:      "111111111111111111",
		NumeroContentieux: "2222222222222",
		DateCreation:      date1,
		DateEcheance:      date2,
		DureeDelai:        444,
		Denomination:      "Test",
		Indic6m:           "SUP",
		AnneeCreation:     2014,
		MontantEcheancier: 200000,
		Stade:             "PR",
		Action:            "SUR PO",
	}

	if !compareDelais(test_1, expected_1) {
		t.Error("Delai parser is not working well")
	}
}

func compareDelais(delai1 Delai, delai2 Delai) bool {
	return delai1.NumeroCompte == delai2.NumeroCompte &&
		delai1.NumeroContentieux == delai2.NumeroContentieux &&
		delai1.DateCreation.Equal(delai2.DateCreation) &&
		delai1.DateEcheance.Equal(delai2.DateEcheance) &&
		delai1.DureeDelai == delai2.DureeDelai &&
		delai1.Denomination == delai2.Denomination &&
		delai1.Indic6m == delai2.Indic6m &&
		delai1.AnneeCreation == delai2.AnneeCreation &&
		delai1.MontantEcheancier == delai2.MontantEcheancier &&
		delai1.Stade == delai2.Stade &&
		delai1.Action == delai2.Action
}

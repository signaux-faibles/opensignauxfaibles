package sirene

import (
	"testing"
	"time"

	"github.com/signaux-faibles/gournal"
)

func TestSirene(t *testing.T) {
	test_row_1 := []string{"005520135", "00038", "00552013500038", "O", "2007-04-20", "", "", "", "2008-01-04T17:54:12", "true", "2", "", "70", "", "RUE", "DE LAUSANNE", "01220", "DIVONNE LES BAINS", "", "", "01143", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "2007-11-19", "F", "", "", "", "", "17.1P", "NAFRev1", "N", "6.143361", "46.357915", "0.91", "housenumber", "70 Rue de Lausanne 01220 Divonne-les-Bains", "ADRNIVX_0000002100454561", "G", "70 RUE DE LAUSANNE", ""}

	var tracker gournal.Tracker
	test_1, tracker := readLineEtablissement(test_row_1, tracker)

	loc, _ := time.LoadLocation("Europe/Paris")
	creation := time.Date(2007, 4, 20, 0, 0, 0, 0, loc)
	longitude := 6.143361
	lattitude := 46.357915

	expected_1 := Sirene{
		Siren:       "005520135",
		Nic:         "00038",
		Siege:       true,
		NumVoie:     "70",
		IndRep:      "",
		TypeVoie:    "RUE",
		CodePostal:  "01220",
		Cedex:       "",
		Departement: "01",
		Commune:     "DIVONNE LES BAINS",
		APE:         "171P",
		Creation:    &creation,
		Longitude:   &longitude,
		Lattitude:   &lattitude,
		Adresse: [6]string{
			"",
			"",
			"DE LAUSANNE",
			"01220",
			"DIVONNE LES BAINS",
			"70 Rue de Lausanne 01220 Divonne-les-Bains",
		},
	}
	if !compareSirene(test_1, expected_1) {
		t.Error("Structure of Sirene is not read as expected")
	}
}

func compareSirene(siren1 Sirene, siren2 Sirene) bool {
	return (siren1.Siren == siren2.Siren) &&
		(siren1.Nic == siren2.Nic) &&
		(siren1.Siege == siren2.Siege) &&
		(siren1.NumVoie == siren2.NumVoie) &&
		(siren1.IndRep == siren2.IndRep) &&
		(siren1.TypeVoie == siren2.TypeVoie) &&
		(siren1.CodePostal == siren2.CodePostal) &&
		(siren1.Cedex == siren2.Cedex) &&
		(siren1.Departement == siren2.Departement) &&
		(siren1.Commune == siren2.Commune) &&
		(siren1.APE == siren2.APE) &&
		(*siren1.Creation).Equal(*siren2.Creation) &&
		(*siren1.Longitude == *siren2.Longitude) &&
		(*siren1.Lattitude == *siren2.Lattitude) &&
		(siren1.Adresse == siren2.Adresse)
}

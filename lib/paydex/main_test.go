package paydex

import (
	"bufio"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Le Paydex – index de paiement – mesure statistiquement la régularité de
// paiement d’une entreprise vis-à-vis de ses fournisseurs.
// Il est exprimé en nombre de jours de retard de paiement moyen,
// basé sur trois expériences de paiement minimum
// (provenant de trois fournisseurs distincts).

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestPaydex(t *testing.T) {
	t.Run("can read an actual file", func(t *testing.T) {
		filePath := "../../../paydex/E_202011095813_Retro-Paydex_20201207.csv" // => SIREN;NB_JOURS;NB_JOURS_LIB;DATE_VALEUR
		// filePath := "../../../paydex/E_202011095813_Identite_20201207.csv" // => R�f�rence Client;Siren;Siret;Raison sociale 1;Raison sociale 2;Enseigne;Sigle;Compl�ment d'adresse;Adresse;Distribution sp�ciale;Code postal et bureau distributeur;Pays;Code postal;Ville;Qualit� Etablissement;Code type d'�tablissement;Libell� type d'�tablissement;Etat d'activit� �tablissement;Etat d'activit� entreprise;Etat de proc�dure collective;Diffusible;
		file, err := os.Open(filePath)
		if err != nil {
			t.Error(err)
		}
		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'
		row, err := reader.Read()
		log.Println(row)
		row, err = reader.Read()
		log.Println(row)
	})

	t.Run("can parse a line", func(t *testing.T) {
		row := []string{"000000001", "2", "2 jours", "15/12/2018"}
		expected := Paydex{
			Siren:   "000000001",
			Periode: time.Date(2018, 12, 01, 00, 00, 00, 0, time.UTC),
			Jours:   2,
		}
		actual := parsePaydexLine(row)
		assert.Equal(t, expected, actual)
	})
}

package paydex

import (
	"bufio"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"testing"
)

// Le Paydex – index de paiement – mesure statistiquement la régularité de
// paiement d’une entreprise vis-à-vis de ses fournisseurs.
// Il est exprimé en nombre de jours de retard de paiement moyen,
// basé sur trois expériences de paiement minimum
// (provenant de trois fournisseurs distincts).

var update = flag.Bool("update", false, "Update the expected test values in golden file")

func TestPaydex(t *testing.T) {
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
}

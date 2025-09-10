package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFileTypeFromFilename(t *testing.T) {

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category ValidFileType
	}{
		// guessed from urssaf files found on stockage/goub server
		{"sigfaible_debits.csv", debit},
		{"sigfaible_cotisdues.csv", cotisation},
		{"sigfaible_pcoll.csv", procol},
		{"sigfaible_etablissement_utf8.csv", adminUrssaf},
		{"sigfaible_effectif_siret.csv", effectif},
		{"sigfaible_effectif_siren.csv", effectifEnt},
		{"sigfaible_delais.csv", delai},
		{"sigfaible_ccsf.csv", ccsf},

		// compressed version of urssaf files
		{"sigfaible_ccsf.csv.gz", ccsf},
		{"sigfaible_etablissement_utf8.csv.gz", adminUrssaf}, // sfdata parser name: "comptes"
		{"sigfaible_cotisdues.csv.gz", cotisation},
		{"sigfaible_debits.csv.gz", debit},
		{"sigfaible_delais.csv.gz", delai},
		{"sigfaible_effectif_siren.csv.gz", effectifEnt},
		{"sigfaible_effectif_siret.csv.gz", effectif},
		{"sigfaible_pcoll.csv.gz", procol},

		// guessed from dgefp files
		{"consommation_ap.csv", apconso},
		{"demande_ap.csv", apdemande},

		// others
		{"Diane_Export_4.txt", diane},
		{"sigfaibles_debits.csv", debit},
		{"sigfaibles_debits2.csv", debit},
		{"diane_req_2002.csv", diane},
		{"diane_req_dom_2002.csv", diane},
		{"effectif_dom.csv", effectif},
		{"filter_siren_2002.csv", filter},
		{"sireneUL.csv", sireneUl},
		{"StockEtablissement_utf8_geo.csv", sirene},
		{"StockEtablissement_utf8_geo.csv", sirene},
		{"E_202011095813_Retro-Paydex_20201207.csv", paydex},
		{"E_202011095813_Identite_20201207.csv", ""}, // not paydex
		{"Ellisphère-Tête de groupe-FinalV2-2015.xlsx", ellisphere},
	}
	for _, testCase := range cases {
		t.Run("should return "+string(testCase.category)+" for file "+testCase.name, func(t *testing.T) {
			got := ExtractFileTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}

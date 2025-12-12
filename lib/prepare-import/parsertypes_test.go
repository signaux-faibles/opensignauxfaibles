package prepareimport

import (
	"opensignauxfaibles/lib/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractParserTypeFromFilename(t *testing.T) {

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category engine.ParserType
	}{
		// guessed from urssaf files found on stockage/goub server
		{"sigfaible_debits.csv", engine.Debit},
		{"sigfaible_cotisdues.csv", engine.Cotisation},
		{"sigfaible_pcoll.csv", engine.Procol},
		{"sigfaible_effectif_siret.csv", engine.Effectif},
		{"sigfaible_effectif_siren.csv", engine.EffectifEnt},
		{"sigfaible_delais.csv", engine.Delai},
		{"sigfaible_ccsf.csv", engine.Ccsf},

		// compressed version of urssaf files
		{"sigfaible_ccsf.csv.gz", engine.Ccsf},
		{"sigfaible_cotisdues.csv.gz", engine.Cotisation},
		{"sigfaible_debits.csv.gz", engine.Debit},
		{"sigfaible_delais.csv.gz", engine.Delai},
		{"sigfaible_effectif_siren.csv.gz", engine.EffectifEnt},
		{"sigfaible_effectif_siret.csv.gz", engine.Effectif},
		{"sigfaible_pcoll.csv.gz", engine.Procol},

		// guessed from dgefp files
		{"consommation_ap.csv", engine.Apconso},
		{"demande_ap.csv", engine.Apdemande},

		// others
		{"Diane_Export_4.txt", engine.Diane},
		{"sigfaibles_debits.csv", engine.Debit},
		{"sigfaibles_debits2.csv", engine.Debit},
		{"diane_req_2002.csv", engine.Diane},
		{"diane_req_dom_2002.csv", engine.Diane},
		{"effectif_dom.csv", engine.Effectif},
		{"filter.csv", engine.Filter},
		{"filter_siren_2002.csv", engine.Filter},
		{"sireneUL.csv", engine.SireneUl},
		{"StockEtablissement_utf8_geo.csv", engine.Sirene},
		{"StockEtablissement_utf8_geo.csv", engine.Sirene},
		{"E_202011095813_Retro-Paydex_20201207.csv", engine.Paydex},
		{"E_202011095813_Identite_20201207.csv", ""}, // not paydex
		{"Ellisphère-Tête de groupe-FinalV2-2015.xlsx", engine.Ellisphere},
	}
	for _, testCase := range cases {
		t.Run("should return "+string(testCase.category)+" for file "+testCase.name, func(t *testing.T) {
			got := ExtractParserTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}

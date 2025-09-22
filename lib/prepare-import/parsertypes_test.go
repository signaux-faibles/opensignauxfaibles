package prepareimport

import (
	"opensignauxfaibles/lib/base"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractParserTypeFromFilename(t *testing.T) {

	// inspired by https://github.com/golang/go/wiki/TableDrivenTests
	cases := []struct {
		name     string
		category base.ParserType
	}{
		// guessed from urssaf files found on stockage/goub server
		{"sigfaible_debits.csv", base.Debit},
		{"sigfaible_cotisdues.csv", base.Cotisation},
		{"sigfaible_pcoll.csv", base.Procol},
		{"sigfaible_etablissement_utf8.csv", base.AdminUrssaf},
		{"sigfaible_effectif_siret.csv", base.Effectif},
		{"sigfaible_effectif_siren.csv", base.EffectifEnt},
		{"sigfaible_delais.csv", base.Delai},
		{"sigfaible_ccsf.csv", base.Ccsf},

		// compressed version of urssaf files
		{"sigfaible_ccsf.csv.gz", base.Ccsf},
		{"sigfaible_etablissement_utf8.csv.gz", base.AdminUrssaf}, // sfdata parser name: "comptes"
		{"sigfaible_cotisdues.csv.gz", base.Cotisation},
		{"sigfaible_debits.csv.gz", base.Debit},
		{"sigfaible_delais.csv.gz", base.Delai},
		{"sigfaible_effectif_siren.csv.gz", base.EffectifEnt},
		{"sigfaible_effectif_siret.csv.gz", base.Effectif},
		{"sigfaible_pcoll.csv.gz", base.Procol},

		// guessed from dgefp files
		{"consommation_ap.csv", base.Apconso},
		{"demande_ap.csv", base.Apdemande},

		// others
		{"Diane_Export_4.txt", base.Diane},
		{"sigfaibles_debits.csv", base.Debit},
		{"sigfaibles_debits2.csv", base.Debit},
		{"diane_req_2002.csv", base.Diane},
		{"diane_req_dom_2002.csv", base.Diane},
		{"effectif_dom.csv", base.Effectif},
		{"filter_siren_2002.csv", base.Filter},
		{"sireneUL.csv", base.SireneUl},
		{"StockEtablissement_utf8_geo.csv", base.Sirene},
		{"StockEtablissement_utf8_geo.csv", base.Sirene},
		{"E_202011095813_Retro-Paydex_20201207.csv", base.Paydex},
		{"E_202011095813_Identite_20201207.csv", ""}, // not paydex
		{"Ellisphère-Tête de groupe-FinalV2-2015.xlsx", base.Ellisphere},
	}
	for _, testCase := range cases {
		t.Run("should return "+string(testCase.category)+" for file "+testCase.name, func(t *testing.T) {
			got := ExtractParserTypeFromFilename(testCase.name)
			assert.Equal(t, testCase.category, got)
		})
	}
}

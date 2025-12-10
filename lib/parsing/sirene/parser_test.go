package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSirene.json")
var testData = engine.NewBatchFile("testData", "sireneTestData.csv")

func TestSirene(t *testing.T) {
	engine.TestParserOutput(t, NewSireneParser(), engine.NewEmptyCache(), testData, golden, *update)

	t.Run("should fail if a required column is missing", func(t *testing.T) {
		output := engine.RunParserInline(t, NewSireneParser(), []string{"siren"}) // many columns are missing
		assert.Equal(t, []engine.Tuple(nil), output.Tuples, "should return no tuples")
		assert.Regexp(t, regexp.MustCompile("column [^ ]+ not found"), engine.GetFatalError(output))
	})
}

func TestExtractDepartement(t *testing.T) {
	tests := []struct {
		name        string
		codeCommune string
		want        string
	}{
		{
			name:        "département métropolitain classique",
			codeCommune: "75001",
			want:        "75",
		},
		{
			name:        "département 01 (Ain)",
			codeCommune: "01000",
			want:        "01",
		},
		{
			name:        "Corse-du-Sud (200)",
			codeCommune: "2A000",
			want:        "2A",
		},
		{
			name:        "Corse-du-Sud (201)",
			codeCommune: "2A100",
			want:        "2A",
		},
		{
			name:        "Haute-Corse (202)",
			codeCommune: "2B200",
			want:        "2B",
		},
		{
			name:        "Haute-Corse (206)",
			codeCommune: "2B600",
			want:        "2B",
		},
		{
			name:        "Guadeloupe (971)",
			codeCommune: "97110",
			want:        "971",
		},
		{
			name:        "Martinique (972)",
			codeCommune: "97200",
			want:        "972",
		},
		{
			name:        "Guyane (973)",
			codeCommune: "97300",
			want:        "973",
		},
		{
			name:        "Réunion (974)",
			codeCommune: "97400",
			want:        "974",
		},
		{
			name:        "Saint-Pierre-et-Miquelon (975)",
			codeCommune: "97500",
			want:        "975",
		},
		{
			name:        "Mayotte (976)",
			codeCommune: "97600",
			want:        "976",
		},
		{
			name:        "Saint-Barthélemy (977)",
			codeCommune: "97700",
			want:        "977",
		},
		{
			name:        "Saint-Martin (978)",
			codeCommune: "97800",
			want:        "978",
		},
		{
			name:        "Wallis-et-Futuna (986)",
			codeCommune: "98600",
			want:        "986",
		},
		{
			name:        "Polynésie française (987)",
			codeCommune: "98700",
			want:        "987",
		},
		{
			name:        "Nouvelle-Calédonie (988)",
			codeCommune: "98800",
			want:        "988",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDepartement(tt.codeCommune)
			assert.Equal(t, tt.want, got)
		})
	}
}

var sampleRow = []string{
	"956520573", "44947", "67323298386574", "t", "2007-04-20", "", "", "", "",
	"false", "2", "", "", "", "CHE", "dguylittc", "64276", "rjkxla", "", "",
	"73187", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	"", "2007-11-19", "A", "", "", "", "", "52.21Z", "NAFRev2", "w", "6.125600",
	"14.106745", "4.74", "rmyvud", "pycqborak26441jvefqp", "52756_a195_mpo13m",
	"p", "hkcilbeed", "",
}

var headers = []string{
	"siren", "nic", "siret", "statutDiffusionEtablissement", "dateCreationEtablissement", "trancheEffectifsEtablissement", "anneeEffectifsEtablissement", "activitePrincipaleRegistreMetiersEtablissement", "dateDernierTraitementEtablissement", "etablissementSiege", "nombrePeriodesEtablissement", "complementAdresseEtablissement", "numeroVoieEtablissement", "indiceRepetitionEtablissement", "typeVoieEtablissement", "libelleVoieEtablissement", "codePostalEtablissement", "libelleCommuneEtablissement", "libelleCommuneEtrangerEtablissement", "distributionSpecialeEtablissement", "codeCommuneEtablissement", "codeCedexEtablissement", "libelleCedexEtablissement", "codePaysEtrangerEtablissement", "libellePaysEtrangerEtablissement", "complementAdresse2Etablissement", "numeroVoie2Etablissement", "indiceRepetition2Etablissement", "typeVoie2Etablissement", "libelleVoie2Etablissement", "codePostal2Etablissement", "libelleCommune2Etablissement", "libelleCommuneEtranger2Etablissement", "distributionSpeciale2Etablissement", "codeCommune2Etablissement", "codeCedex2Etablissement", "libelleCedex2Etablissement", "codePaysEtranger2Etablissement", "libellePaysEtranger2Etablissement", "dateDebut", "etatAdministratifEtablissement", "enseigne1Etablissement", "enseigne2Etablissement", "enseigne3Etablissement", "denominationUsuelleEtablissement", "activitePrincipaleEtablissement", "nomenclatureActivitePrincipaleEtablissement", "caractereEmployeurEtablissement", "longitude", "latitude", "geo_score", "geo_type", "geo_adresse", "geo_id", "geo_ligne", "geo_l4", "geo_l5",
}

var colIndex, _ = parsing.HeaderIndexer{Dest: Sirene{}}.Index(headers, true)

func TestNAFFiltering(t *testing.T) {
	testcases := []struct {
		name                 string
		nomenclatureActivite string
		shouldFilter         bool
	}{
		{
			"NAFRev2 should be included",
			"NAFRev2",
			false,
		},
		{
			"NAFRev1 should be excluded",
			"NAFRev1",
			true,
		},
		{
			"NAF1993 should be excluded",
			"NAF1993",
			true,
		},
		{
			"NAP should be excluded",
			"NAP",
			true,
		},
		{
			"missing nomenclature should be excluded",
			"",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			row := make([]string, len(sampleRow))
			copy(row, sampleRow)

			index, _ := colIndex.Get("nomenclatureActivitePrincipaleEtablissement")
			row[index] = tc.nomenclatureActivite

			parser := &sireneRowParser{}
			res := &engine.ParsedLineResult{}

			parser.ParseRow(row, res, colIndex)

			if tc.shouldFilter {
				assert.NotNil(t, res.FilterError)
			} else {
				assert.Nil(t, res.FilterError)
			}

		})
	}

}

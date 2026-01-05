package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing"
)

var update = flag.Bool("update", false, "Update the expected test values in golden file")

var golden = filepath.Join("testData", "expectedSirene.json")
var testData = engine.NewBatchFile("testData", "sireneTestData.csv")

func TestSirene(t *testing.T) {
	engine.TestParserOutput(t, NewSireneParser(), testData, golden, *update)

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

func TestSireneParser(t *testing.T) {
	testCases := []struct {
		csvRow   []string
		expected Sirene
	}{
		{
			csvRow: []string{
				"123456789", "00012", "12345678900012", "O", "2010-05-15", "", "", "", "",
				"true", "1", "Bat A", "10", "B", "RUE", "de la Paix", "75001", "Paris", "",
				"", "75101", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
				"", "2010-05-15", "A", "", "", "", "", "47.11A", "NAFRev2", "O", "2.352341",
				"48.864716", "0.9", "housenumber", "10 BIS RUE de la Paix 75001 Paris", "75101_1234_00001",
				"w", "10 BIS RUE de la Paix 75001 Paris", "",
			},
			expected: Sirene{
				Siren:             "123456789",
				Nic:               "00012",
				Siret:             "12345678900012",
				Siege:             true,
				ComplementAdresse: "Bat A",
				NumVoie:           "10",
				IndRep:            "BIS",
				TypeVoie:          "RUE",
				Voie:              "de la Paix",
				CodePostal:        "75001",
				Commune:           "Paris",
				CodeCommune:       "75101",
				Departement:       "75",
				APE:               "47.11A",
				Creation:          parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2010-05-15")),
				Longitude:         2.352341,
				Latitude:          48.864716,
				EstActif:          true,
			},
		},
		{
			csvRow: []string{
				"987654321", "00025", "98765432100025", "O", "2015-03-20", "", "", "", "",
				"false", "1", "", "25", "", "AV", "des Champs-Elysées", "13001", "Marseille", "",
				"", "13201", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
				"", "2015-03-20", "F", "", "", "", "", "56.10A", "NAFRev2", "O", "5.369780",
				"43.296482", "0.85", "street", "25 AVENUE des Champs-Elysées 13001 Marseille", "13201_5678_00001",
				"w", "25 AVENUE des Champs-Elysées", "",
			},
			expected: Sirene{
				Siren:        "987654321",
				Nic:          "00025",
				Siret:        "98765432100025",
				Siege:        false,
				NumVoie:      "25",
				IndRep:       "",
				TypeVoie:     "AVENUE",
				Voie:         "des Champs-Elysées",
				CodePostal:   "13001",
				Commune:      "Marseille",
				CodeCommune:  "13201",
				Departement:  "13",
				APE:          "56.10A",
				Creation:     parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2015-03-20")),
				Longitude:    5.369780,
				Latitude:     43.296482,
				EstActif:     false,
			},
		},
		{
			csvRow: []string{
				"111222333", "00030", "11122233300030", "O", "2018-11-10", "", "", "", "",
				"true", "1", "", "3", "T", "CHE", "du Moulin", "97110", "Pointe-à-Pitre", "",
				"", "97120", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
				"", "2018-11-10", "A", "", "", "", "", "62.01Z", "NAFRev2", "N", "",
				"", "", "", "", "",
				"", "", "",
			},
			expected: Sirene{
				Siren:        "111222333",
				Nic:          "00030",
				Siret:        "11122233300030",
				Siege:        true,
				NumVoie:      "3",
				IndRep:       "TER",
				TypeVoie:     "CHEMIN",
				Voie:         "du Moulin",
				CodePostal:   "97110",
				Commune:      "Pointe-à-Pitre",
				CodeCommune:  "97120",
				Departement:  "971",
				APE:          "62.01Z",
				Creation:     parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2018-11-10")),
				Longitude:    0,
				Latitude:     0,
				EstActif:     true,
			},
		},
	}

	parser := NewSireneParser()
	header := "siren,nic,siret,statutDiffusionEtablissement,dateCreationEtablissement,trancheEffectifsEtablissement,anneeEffectifsEtablissement,activitePrincipaleRegistreMetiersEtablissement,dateDernierTraitementEtablissement,etablissementSiege,nombrePeriodesEtablissement,complementAdresseEtablissement,numeroVoieEtablissement,indiceRepetitionEtablissement,typeVoieEtablissement,libelleVoieEtablissement,codePostalEtablissement,libelleCommuneEtablissement,libelleCommuneEtrangerEtablissement,distributionSpecialeEtablissement,codeCommuneEtablissement,codeCedexEtablissement,libelleCedexEtablissement,codePaysEtrangerEtablissement,libellePaysEtrangerEtablissement,complementAdresse2Etablissement,numeroVoie2Etablissement,indiceRepetition2Etablissement,typeVoie2Etablissement,libelleVoie2Etablissement,codePostal2Etablissement,libelleCommune2Etablissement,libelleCommuneEtranger2Etablissement,distributionSpeciale2Etablissement,codeCommune2Etablissement,codeCedex2Etablissement,libelleCedex2Etablissement,codePaysEtranger2Etablissement,libellePaysEtranger2Etablissement,dateDebut,etatAdministratifEtablissement,enseigne1Etablissement,enseigne2Etablissement,enseigne3Etablissement,denominationUsuelleEtablissement,activitePrincipaleEtablissement,nomenclatureActivitePrincipaleEtablissement,caractereEmployeurEtablissement,longitude,latitude,geo_score,geo_type,geo_adresse,geo_id,geo_ligne,geo_l4,geo_l5"
	for _, tc := range testCases {
		instance := parser.New(parsing.CreateReader(header, ",", tc.csvRow))
		instance.Init(engine.NoFilter, nil)
		res := &engine.ParsedLineResult{}
		instance.ReadNext(res)

		assert.Empty(t, res.Errors)
		require.GreaterOrEqual(t, len(res.Tuples), 1)

		sirene, ok := res.Tuples[0].(Sirene)
		if !ok {
			t.Fatal("tuple type is not as expected")
		}

		assert.Equal(t, tc.expected, sirene)
	}
}

func TestSireneMissingColumns(t *testing.T) {
	t.Run("should fail if one column misses", func(t *testing.T) {
		parser := NewSireneParser()
		instance := parser.New(parsing.CreateReader("siren,nic", ",", []string{}))
		err := instance.Init(engine.NoFilter, nil)

		assert.Error(t, err, "should report a fatal error")
		assert.Regexp(t, regexp.MustCompile("column [^ ]+ not found"), err.Error())
	})

	t.Run("should fail if etablissementSiege column is missing", func(t *testing.T) {
		parser := NewSireneParser()
		instance := parser.New(parsing.CreateReader("siren,nic,siret,statutDiffusionEtablissement,dateCreationEtablissement,trancheEffectifsEtablissement", ",", []string{}))
		err := instance.Init(engine.NoFilter, nil)

		assert.Error(t, err, "should report a fatal error")
		assert.Contains(t, err.Error(), "column etablissementSiege not found")
	})
}

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

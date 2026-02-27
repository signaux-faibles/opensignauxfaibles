package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"strings"
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

func TestSireneParser(t *testing.T) {
	testCases := []struct {
		fields   map[string]string
		expected Sirene
	}{
		{
			fields: map[string]string{
				"siren":                                    "123456789",
				"nic":                                      "00012",
				"etablissementSiege":                      "true",
				"complementAdresseEtablissement":          "Bat A",
				"numeroVoieEtablissement":                 "10",
				"indiceRepetitionEtablissement":           "B",
				"typeVoieEtablissement":                   "RUE",
				"libelleVoieEtablissement":                "de la Paix",
				"codePostalEtablissement":                 "75001",
				"libelleCommuneEtablissement":             "Paris",
				"codeCommuneEtablissement":                "75101",
				"activitePrincipaleEtablissement":         "47.11A",
				"nomenclatureActivitePrincipaleEtablissement": "NAFRev2",
				"dateCreationEtablissement":               "2010-05-15",
				"etatAdministratifEtablissement":          "A",
				"longitude":                               "2.352341",
				"latitude":                                "48.864716",
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
			fields: map[string]string{
				"siren":                                    "987654321",
				"nic":                                      "00025",
				"etablissementSiege":                      "false",
				"numeroVoieEtablissement":                 "25",
				"typeVoieEtablissement":                   "AV",
				"libelleVoieEtablissement":                "des Champs-Elysées",
				"codePostalEtablissement":                 "13001",
				"libelleCommuneEtablissement":             "Marseille",
				"codeCommuneEtablissement":                "13201",
				"activitePrincipaleEtablissement":         "56.10A",
				"nomenclatureActivitePrincipaleEtablissement": "NAFRev2",
				"dateCreationEtablissement":               "2015-03-20",
				"etatAdministratifEtablissement":          "F",
				"longitude":                               "5.369780",
				"latitude":                                "43.296482",
			},
			expected: Sirene{
				Siren:       "987654321",
				Nic:         "00025",
				Siret:       "98765432100025",
				Siege:       false,
				NumVoie:     "25",
				TypeVoie:    "AVENUE",
				Voie:        "des Champs-Elysées",
				CodePostal:  "13001",
				Commune:     "Marseille",
				CodeCommune: "13201",
				Departement: "13",
				APE:         "56.10A",
				Creation:    parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2015-03-20")),
				Longitude:   5.369780,
				Latitude:    43.296482,
				EstActif:    false,
			},
		},
		{
			fields: map[string]string{
				"siren":                                    "111222333",
				"nic":                                      "00030",
				"etablissementSiege":                      "true",
				"numeroVoieEtablissement":                 "3",
				"indiceRepetitionEtablissement":           "T",
				"typeVoieEtablissement":                   "CHE",
				"libelleVoieEtablissement":                "du Moulin",
				"codePostalEtablissement":                 "97110",
				"libelleCommuneEtablissement":             "Pointe-à-Pitre",
				"codeCommuneEtablissement":                "97120",
				"activitePrincipaleEtablissement":         "62.01Z",
				"nomenclatureActivitePrincipaleEtablissement": "NAFRev2",
				"dateCreationEtablissement":               "2018-11-10",
				"etatAdministratifEtablissement":          "A",
			},
			expected: Sirene{
				Siren:       "111222333",
				Nic:         "00030",
				Siret:       "11122233300030",
				Siege:       true,
				NumVoie:     "3",
				IndRep:      "TER",
				TypeVoie:    "CHEMIN",
				Voie:        "du Moulin",
				CodePostal:  "97110",
				Commune:     "Pointe-à-Pitre",
				CodeCommune: "97120",
				Departement: "971",
				APE:         "62.01Z",
				Creation:    parsing.TimePtr(parsing.MustParseTime("2006-01-02", "2018-11-10")),
				EstActif:    true,
			},
		},
	}

	parser := NewSireneParser()
	header := strings.Join(sireneColumns, ",")
	for _, tc := range testCases {
		row, err := makeRow(tc.fields)
		require.NoError(t, err)
		instance := parser.New(strings.NewReader(header + "\n" + row))
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

	header := strings.Join(sireneColumns, ",")
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			row, err := makeRow(map[string]string{
				"siren":                                    "123456789",
				"nic":                                      "00012",
				"etablissementSiege":                       "true",
				"etatAdministratifEtablissement":           "A",
				"activitePrincipaleEtablissement":          "47.11A",
				"nomenclatureActivitePrincipaleEtablissement": tc.nomenclatureActivite,
			})
			require.NoError(t, err)
			parser := NewSireneParser()
			instance := parser.New(strings.NewReader(header + "\n" + row))
			instance.Init(engine.NoFilter, nil)
			res := &engine.ParsedLineResult{}
			instance.ReadNext(res)

			if tc.shouldFilter {
				assert.NotNil(t, res.FilterError)
			} else {
				assert.Nil(t, res.FilterError)
			}
		})
	}
}

package sirene

import (
	"flag"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/engine"
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
		name       string
		codePostal string
		want       string
	}{
		{
			name:       "département métropolitain classique",
			codePostal: "75001",
			want:       "75",
		},
		{
			name:       "département 01 (Ain)",
			codePostal: "01000",
			want:       "01",
		},
		{
			name:       "Corse-du-Sud (200)",
			codePostal: "20000",
			want:       "2A",
		},
		{
			name:       "Corse-du-Sud (201)",
			codePostal: "20100",
			want:       "2A",
		},
		{
			name:       "Haute-Corse (202)",
			codePostal: "20200",
			want:       "2B",
		},
		{
			name:       "Haute-Corse (206)",
			codePostal: "20600",
			want:       "2B",
		},
		{
			name:       "Guadeloupe (971)",
			codePostal: "97110",
			want:       "971",
		},
		{
			name:       "Martinique (972)",
			codePostal: "97200",
			want:       "972",
		},
		{
			name:       "Guyane (973)",
			codePostal: "97300",
			want:       "973",
		},
		{
			name:       "Réunion (974)",
			codePostal: "97400",
			want:       "974",
		},
		{
			name:       "Saint-Pierre-et-Miquelon (975)",
			codePostal: "97500",
			want:       "975",
		},
		{
			name:       "Mayotte (976)",
			codePostal: "97600",
			want:       "976",
		},
		{
			name:       "Saint-Barthélemy (977)",
			codePostal: "97700",
			want:       "977",
		},
		{
			name:       "Saint-Martin (978)",
			codePostal: "97800",
			want:       "978",
		},
		{
			name:       "Wallis-et-Futuna (986)",
			codePostal: "98600",
			want:       "986",
		},
		{
			name:       "Polynésie française (987)",
			codePostal: "98700",
			want:       "987",
		},
		{
			name:       "Nouvelle-Calédonie (988)",
			codePostal: "98800",
			want:       "988",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDepartement(tt.codePostal)
			assert.Equal(t, tt.want, got)
		})
	}
}

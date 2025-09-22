package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/prepare-import/createfilter"
	"opensignauxfaibles/lib/prepare-import/prepareimport"
)

var goldenAdminObject = createfilter.ReadGoldenFile("end_to_end_golden.txt")
var emptyAsString, _ = json.MarshalIndent(base.AdminBatch{}, "", "  ")

func Test_prepare(t *testing.T) {
	effectifData, err := os.ReadFile("./createfilter/test_data.csv")
	if err != nil {
		t.Fatal(err)
	}

	type want struct {
		adminObject string
		error       string
	}

	tests := []struct {
		name  string
		batch string
		want  want
	}{
		{
			"test avec tous les bons paramètres",
			"1802",
			want{adminObject: goldenAdminObject, error: prepareimport.UnsupportedFilesError{}.Error()},
		},
		{
			"test avec un mauvais paramètre batch",
			"180",
			want{adminObject: string(emptyAsString), error: "la clé du batch doit respecter le format requis AAMM"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gzipString, _ := prepareimport.GzipString(prepareimport.SomeText(254781489))

			buildedBatchKey, _ := base.NewBatchKey(tt.batch)

			parentDir := prepareimport.CreateTempFilesWithContent(t, buildedBatchKey, map[string][]byte{
				"sigfaibles_effectif_siret.csv":            effectifData,
				"sigfaibles_debits.csv":                    prepareimport.SomeTextAsBytes(254784321489),
				"unsupported.csv":                          prepareimport.SomeTextAsBytes(254788761489),
				"E_202011095813_Retro-Paydex_20201207.csv": prepareimport.SomeTextAsBytes(25477681489),
				"sigfaible_pcoll.csv.gz":                   gzipString,
				"sireneUL.csv":                             ReadFileData(t, "createfilter/test_uniteLegale.csv"),
			})

			actual, err2 := prepare(parentDir, tt.batch)
			assert.ErrorContains(t, err2, tt.want.error)
			objectBytes, err := json.MarshalIndent(actual, "", "  ")
			assert.NoError(t, err)
			assert.Equal(t, tt.want.adminObject, stripAllBasepaths(string(objectBytes), parentDir+"/"))
		})
	}
}

func stripAllBasepaths(jsonStr, basepath string) string {
	return strings.ReplaceAll(jsonStr, basepath, "")
}

func ReadFileData(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

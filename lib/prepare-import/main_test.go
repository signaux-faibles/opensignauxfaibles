package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/prepare-import/createfilter"
	"opensignauxfaibles/lib/prepare-import/prepareimport"
)

var goldenAdminObject = createfilter.ReadGoldenFile("end_to_end_golden.txt")
var emptyAsString, _ = json.MarshalIndent(prepareimport.AdminBatch{}, "", "  ")

func Test_prepare(t *testing.T) {
	effectifData, err := os.ReadFile("./createfilter/test_data.csv")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		batch       string
		finEffectif string
	}
	type want struct {
		adminObject string
		error       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"test avec tous les bons paramètres",
			args{"1802", "2018-01-01"},
			want{adminObject: goldenAdminObject, error: prepareimport.UnsupportedFilesError{}.Error()},
		},
		{
			"test avec un mauvais paramètre batch",
			args{"180", "2018-01-01"},
			want{adminObject: string(emptyAsString), error: "la clé du batch doit respecter le format requis AAMM"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gzipString, _ := prepareimport.GzipString(prepareimport.SomeText(254781489))

			buildedBatchKey, _ := prepareimport.NewBatchKey(tt.args.batch)
			parentDir := prepareimport.CreateTempFilesWithContent(t, buildedBatchKey, map[string][]byte{
				"sigfaibles_effectif_siret.csv":            effectifData,
				"sigfaibles_debits.csv":                    prepareimport.SomeTextAsBytes(254784321489),
				"unsupported.csv":                          prepareimport.SomeTextAsBytes(254788761489),
				"E_202011095813_Retro-Paydex_20201207.csv": prepareimport.SomeTextAsBytes(25477681489),
				"sigfaible_pcoll.csv.gz":                   gzipString,
				"sireneUL.csv":                             ReadFileData(t, "createfilter/test_uniteLegale.csv"),
			})
			actual, err2 := prepare(parentDir, tt.args.batch, tt.args.finEffectif)
			assert.ErrorContains(t, err2, tt.want.error)
			objectBytes, err := json.MarshalIndent(actual, "", "  ")
			assert.NoError(t, err)
			assert.Equal(t, tt.want.adminObject, string(objectBytes))
		})
	}
}

func ReadFileData(t *testing.T, filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

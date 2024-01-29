//go:build local

package main

import (
	"encoding/csv"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var EXPECTED_HEADERS = []string{"siren", "état_organisation", "code_paydex", "nbr_jrs_retard", "nbr_fournisseurs", "encours_étudiés", "note_100_alerteur_plus_30", "note_100_alerteur_plus_90_jours", "date_valeur"}

func Test_convertAndConcat(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "output_*.csv")
	log.Print("fichier généré ", output.Name())
	require.NoError(t, err)
	convertAndConcat(
		[]string{
			"resources/SF_DATA_20230706.txt",
			"resources/2312/S_202011095834-3_202310020319.csv",
			"resources/2312/S_202011095834-3_202311010315.csv",
			"resources/2401/S_202011095834-3_202312011707.csv",
			"resources/2401/S_202011095834-3_202401010331.csv",
		},
		output,
	)
	err = output.Close()
	require.NoError(t, err)

	generated, err := os.Open(output.Name())
	require.NoError(t, err)
	defer output.Close()

	csvR := csv.NewReader(generated)
	// read headers
	headers, err := csvR.Read()
	require.NoError(t, err)
	assert.Equal(t, EXPECTED_HEADERS, headers)
	// première ligne -> provient du fichier stock
	expectedFirstLine := []string{"005480546", "Fermé", "", "", "4", "4390.532", "", "", "2022-05-15"}
	firstLine, err := csvR.Read()
	require.NoError(t, err)
	assert.Equal(t, expectedFirstLine, firstLine)

	//// toutes les lignes
	//var lastLine []string
	//for {
	//	currentLine, err := csvR.Read()
	//	if err == io.EOF {
	//		break
	//	}
	//	lastLine = currentLine // on n'est pas à la fin du fichier, alors peut-être est-ce la dernière ligne ?
	//	if err != nil {
	//		t.Errorf("erreur de lecture du fichier cible : %v", err)
	//	}
	//	assert.Len(t, currentLine, len(EXPECTED_HEADERS))
	//}

	//expectedLastLine := []string{"999990286", "Actif", "070", "15", "12", "70873", "22", "13", "2023-10-30"}
	//assert.Equal(t, expectedLastLine, lastLine)
}

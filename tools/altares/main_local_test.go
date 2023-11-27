//go:build local

package main

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var EXPECTED_HEADERS = []string{"\ufeffsiren", "état_organisation", "code_paydex", "paydex", "nbr_jrs_retard", "nbr_fournisseurs", "encours_étudiés", "note_100_alerteur_plus_30", "note_100_alerteur_plus_90_jours", "date_valeur"}

func Test_convertAndConcat(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "output_*.csv")
	require.NoError(t, err)
	convertAndConcat(
		[]string{"resources/SF_DATA_20230706.txt", "resources/S_202011095834-3_202310020319.csv", "resources/S_202011095834-3_202311010315.csv"},
		output,
	)
	err = output.Close()
	require.NoError(t, err)

	output, err = os.Open(output.Name())
	require.NoError(t, err)

	csvR := csv.NewReader(output)
	// read headers
	headers, err := csvR.Read()
	require.NoError(t, err)
	assert.Equal(t, EXPECTED_HEADERS, headers)
}

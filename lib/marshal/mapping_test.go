package marshal

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/stretchr/testify/assert"
)

func TestGetSiret(t *testing.T) {

	stdTime1, _ := time.Parse("2006-02-01", "2015-01-01")
	stdTime2, _ := time.Parse("2006-02-01", "2016-01-01")
	var stdMapping = map[string][]SiretDate{
		"abc": {
			{"01234567891011", stdTime1},
			{"87654321091011", stdTime2},
		},
	}
	var batch = base.AdminBatch{}
	testCases := []struct {
		compte      string
		date        string
		mapping     Comptes
		expectError bool
		expected    string
	}{
		{"abc", "2015-01-01", map[string][]SiretDate{}, true, ""},
		{"abc", "2014-06-01", stdMapping, false, "01234567891011"},
		{"abc", "2014-12-31", stdMapping, false, "01234567891011"},
		{"abc", "2015-01-01", stdMapping, false, "87654321091011"},
		{"abc", "2016-01-01", stdMapping, true, ""},
	}

	for ind, tc := range testCases {

		var cache = Cache{"comptes": tc.mapping}

		time, _ := time.Parse("2006-02-01", tc.date)
		actual, err := GetSiret(tc.compte, &time, cache, &batch)
		if err != nil && !tc.expectError {
			t.Fatalf("Unexpected error during cache request in test %d: %v", ind, err)
		}
		if err == nil && tc.expectError {
			t.Fatalf("Expected error missing during cache request in test %d: %v", ind, err)
		}
		expected := tc.expected
		if actual != expected {
			t.Log(actual)
			t.Log(expected)
			t.Fatalf("GetSiret failed on test %d", ind)
		}
	}
}

func TestReadSiretMapping(t *testing.T) {

	t.Run("readSiretMapping doit être insensible à la casse des en-têtes de colonnes", func(t *testing.T) {
		batch := base.AdminBatch{}
		csvHeader := `UrsSaf_geStion;DEp;ComptE;Etat_Compte;SIren;Siret;Date_creA_siret;DatE_disp_sirEt;Cle_mD5`
		_, err := readSiretMapping(strings.NewReader(csvHeader), Cache{}, &batch)
		assert.NoError(t, err)
	})

	t.Run("readSiretMapping doit rapporter une erreur s'il manque une colonne", func(t *testing.T) {
		batch := base.AdminBatch{}
		csvHeader := `UrsSaf_geStion`
		_, err := readSiretMapping(strings.NewReader(csvHeader), Cache{}, &batch)
		assert.Error(t, err, "Colonne Compte non trouvée.")
	})

	t.Run("readSiretMapping doit produire les mêmes tuples que d'habitude", func(t *testing.T) {

		var batch = base.AdminBatch{}

		stdTime1, _ := time.Parse("2006-02-01", "2899-01-01")
		stdTime2, _ := time.Parse("2006-02-01", "2015-01-01")
		stdExpected1 := Comptes{
			"abc": []SiretDate{{"01234567891011", stdTime1}},
		}

		stdExpected2 := Comptes{
			"abc": []SiretDate{{"01234567891011", stdTime2}},
		}

		stdExpected3 := Comptes{
			"abc": []SiretDate{
				{"01234567891011", stdTime2},
				{"87654321091011", stdTime1},
			},
		}

		stdFilterCache := Cache{
			"filter": SirenFilter{"012345678": true},
		}

		expectedHeader := "Urssaf_gestion;Dep;Compte;Etat_compte;Siren;Siret;Date_crea_siret;Date_disp_siret"

		testCases := []struct {
			csv         string
			cache       Cache
			expectError bool
			expected    Comptes
		}{
			// No closing date
			{expectedHeader + `
		;;"abc";;;"01234567891011";;""`, Cache{}, false, stdExpected1},
			// With closing date
			{expectedHeader + `
		;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected2},
			// With filtered siret
			{expectedHeader + `
		;;"abc";;;"01234567891011";;"1150101"`, stdFilterCache, false, stdExpected2},
			// With two entries, including excluded siret
			{expectedHeader + `
		;;"abc";;;"01234567891011";;"1150101"
		;;"abc";;;"87654321091011";;""`, stdFilterCache, false, stdExpected2}, // i.e. no mapping stored for 87654321091011, because it's not included by Filter
			// With two entries 1
			{expectedHeader + `
		;;"abc";;;"01234567891011";;"1150101"
		;;"abc";;;"87654321091011";;""`, Cache{}, false, stdExpected3},
			// With two entries 2: different order
			{expectedHeader + `
	    ;;"abc";;;"87654321091011";;""
	    ;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected3},
			// With invalid siret
			{expectedHeader + `
		  ;;"abc";;;"8765432109101A";;""
	    ;;"abc";;;"01234567891011";;"1150101"`, Cache{}, false, stdExpected2},
		}

		for ind, tc := range testCases {
			actual, err := readSiretMapping(strings.NewReader(tc.csv), tc.cache, &batch)
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error during file reading in test %d: %v", ind, err)
			}
			if err == nil && tc.expectError {
				t.Fatalf("expected error missing during  file reading in test %d: %v", ind, err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("Actual %v", actual)
				t.Logf("Expected %v", tc.expected)
				t.Fatalf("ReadSiretMapping failed test %d", ind)
			}
		}
	})
}

func TestGetCompteSiretMapping(t *testing.T) {
	t.Run("GetCompteSiretMapping can read from compressed admin_urssaf file with `gzip:` prefix", func(t *testing.T) {
		expectedComptes := []string{"111982477292496174", "450359886246036238", "636043216536562844"}
		compressedFileData := compressFileData(t, "../urssaf/testData/comptesTestData.csv")
		compressedFile := CreateTempFileWithContent(t, compressedFileData.Bytes())
		cache := NewCache()
		batch := base.MockBatch("admin_urssaf", []string{"gzip:" + compressedFile.Name()})
		actual, err := GetCompteSiretMapping(cache, &batch, OpenAndReadSiretMapping)
		if assert.NoError(t, err) {
			assert.EqualValues(t, expectedComptes, actual.GetSortedKeys())
		}
	})

	t.Run("other test cases", func(t *testing.T) {
		stdTime1, _ := time.Parse("2006-02-01", "2899-01-01")
		stdTime2, _ := time.Parse("2006-02-01", "2016-01-01")
		stdExpected1 := Comptes{
			"abc": []SiretDate{{"01234567891011", stdTime1}},
		}
		stdExpected2 := Comptes{
			"abc": []SiretDate{{"01234567891011", stdTime2}},
		}

		// When file is read, returnd stdExpected1
		mockOpenFile := func(s1 string, s2 base.BatchFile, c Comptes, ca Cache, ba *base.AdminBatch) (Comptes, error) {
			for key := range stdExpected1 {
				c[key] = stdExpected1[key]
			}
			return c, nil
		}

		testCases := []struct {
			cache       Cache
			batch       base.AdminBatch
			expectError bool
			expected    Comptes
		}{
			// Basic reading from file
			{NewCache(), base.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected1},
			// Cache superseeds reading from file
			{Cache{"comptes": stdExpected2}, base.MockBatch("admin_urssaf", []string{"a"}), false, stdExpected2},
			// No cache, no file = error
			{NewCache(), base.MockBatch("otherStuff", []string{"a"}), true, nil},
		}

		for ind, tc := range testCases {

			actual, err := GetCompteSiretMapping(tc.cache, &tc.batch, mockOpenFile)
			if err != nil && !tc.expectError {
				t.Fatalf("Unexpected error during mapping request in test %d: %v", ind, err)
			}
			if err == nil && tc.expectError {
				t.Fatalf("Expected error missing during mapping request in test %d: %v", ind, err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("Actual %v", actual)
				t.Logf("Expected %v", tc.expected)
				t.Fatalf("getSiretMapping failed test %d", ind)
			}
		}
	})
}

func TestGetSiretFromComptesMapping(t *testing.T) {
	t.Run("doit retourner une erreur si un compte n'est pas associé à un numéro de SIRET, à la date d'aujourd'hui", func(t *testing.T) {
		comptes := MockComptesMapping(map[string]string{})
		compteManquant := "636043216536562844"
		// test
		date := time.Now()
		siret, err := GetSiretFromComptesMapping(compteManquant, &date, comptes)
		assert.Equal(t, "", siret)
		assert.Equal(t, "Pas de siret associé au compte "+compteManquant+" à la période "+date.String(), err.Error())
	})

	t.Run("doit retourner une erreur si un compte n'est pas associé à un numéro de SIRET, à une date passée", func(t *testing.T) {
		comptes := MockComptesMapping(map[string]string{})
		compteManquant := "636043216536562844"
		// test
		date := time.Time{}
		siret, err := GetSiretFromComptesMapping(compteManquant, &date, comptes)
		assert.Equal(t, "", siret)
		assert.Equal(t, "Pas de siret associé au compte "+compteManquant+" à la période "+date.String(), err.Error())
	})
}

func compressFileData(t *testing.T, filePath string) (compressedData bytes.Buffer) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	zw := gzip.NewWriter(&compressedData)
	if _, err = zw.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return compressedData
}

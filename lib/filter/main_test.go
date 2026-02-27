package filter

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing/effectif"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var outGoldenFile = "testData/test_golden.txt"
var errGoldenFile = "testData/test_golden_err.txt"

var updateGoldenFile = flag.Bool("update", false, "Update the expected test values in golden file")

func TestCreateFilter(t *testing.T) {
	t.Run("create_filter golden file", func(t *testing.T) {

		var cmdOutput bytes.Buffer
		var cmdError = *bytes.NewBufferString("") // default: no error

		filter, err := Create(engine.NewBatchFile("testData/test_data.csv"), DefaultNbMois, DefaultMinEffectif)
		if err != nil {
			cmdError = *bytes.NewBufferString(err.Error())
		} else {
			csvWriter := NewCsvWriter(&cmdOutput)
			err = csvWriter.Write(filter)
			if err != nil {
				cmdError = *bytes.NewBufferString(err.Error())
			}
		}
		sortedCmdOutput := bytes.NewBufferString(sortOutput(cmdOutput.String()))
		expectedOutput := DiffWithGoldenFile(outGoldenFile, *updateGoldenFile, *sortedCmdOutput)
		expectedError := DiffWithGoldenFile(errGoldenFile, *updateGoldenFile, cmdError)

		assert.Equal(t, string(expectedOutput), sortedCmdOutput.String())
		assert.Equal(t, string(expectedError), cmdError.String())
	})
}

func sortOutput(s string) string {
	lines := strings.Split(s[0:len(s)-1], "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n") + "\n"
}

// Règle: si et seulement si au moins un établissement a eu pendant au moins
// une période un effectif >= 10, on veut l'avoir en base de données, avec
// tous les autres établissements de cette entreprise.
// cf https://github.com/signaux-faibles/opensignauxfaibles/issues/199
func TestOutputPerimeter(t *testing.T) {
	// test de non regression
	t.Run("le département de l'entreprise n'est pas considéré comme une valeur d'effectif", func(t *testing.T) {
		// setup conditions and expectations
		minEffectif := 10
		expectedSirens := []string{"222222222", "333333333"}

		period1, _ := time.Parse("2006-01-02", "2010-01-01")
		period2, _ := time.Parse("2006-01-02", "2010-01-02")

		csvLines := effectif.MakeEffectifEntCSV(
			[]time.Time{period1, period2},
			map[string][]int{
				"000000000": {4, 4}, "111111111": {4, 4},
				"222222222": {14, 14}, // ✅ siren retenu car 14 est bien un effectif ≥ 10
				"333333333": {14, 14}, // ✅ siren retenu car 14 est bien un effectif ≥ 10
			})
		// test: run outputPerimeter() on csv lines
		actualSirens := getOutputPerimeter(csvLines, DefaultNbMois, minEffectif)
		sort.Strings(actualSirens)

		// assert
		assert.Equal(t, expectedSirens, actualSirens)
	})

	t.Run("outputPerimeter ne doit pas contenir deux fois le même siren", func(t *testing.T) {
		// setup conditions and expectations
		minEffectif := 1
		expectedSirens := []string{"111111111", "333333333"}
		csvContent := strings.Join([]string{
			"siren;eff201011;rais_soc",
			"111111111;1;ENTREPRISE", // première entreprise ayant 111111111 comme siren
			"111111111;1;ENTREPRISE", // deuxième entreprise ayant 111111111 comme siren
			"333333333;1;ENTREPRISE",
		}, "\n")
		// test: run outputPerimeter() on csv lines
		actualSirens := getOutputPerimeter(csvContent, DefaultNbMois, minEffectif)
		sort.Strings(actualSirens)
		// assert
		assert.Equal(t, expectedSirens, actualSirens)
	})
}

// wrapper to run outputPerimeter() on a csv string
func getOutputPerimeter(csvContent string, nbMois, minEffectif int) (actualSirens []string) {
	mockFile := engine.NewMockBatchFile(csvContent)
	extractor, err := newEffectifDataExtractor(mockFile)
	if err != nil {
		panic(err)
	}
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	perimeter, err := getImportPerimeter(mockFile, nbMois, minEffectif, extractor)
	if err != nil {
		panic(err)
	}
	for siren := range perimeter {
		fmt.Fprintln(writer, siren)
	}
	writer.Flush()
	return strings.Split(strings.TrimSpace(output.String()), "\n")
}

func TestIsInsidePerimeter(t *testing.T) {
	nbMois := 3 // => seules les valeurs d'effectif des 3 derniers mois vont être considérées
	minEffectif := 10
	testCases := []struct {
		input    []string
		expected bool
	}{
		{[]string{"10", "9", "4", "7", "5"}, false},  // ❌ l'effectif ≥10 date de plus de 3 mois
		{[]string{"10", "20", "4", "7", "5"}, false}, // ❌ l'effectif ≥10
		{[]string{"10", "9", "12", "7", "5"}, true},  // ✅ un effectif ≥10 a été trouvé dans la fenêtre des 3 mois
		{[]string{"10", "9", "12", "", ""}, true},    // ✅ l'absence des 2 dernières valeurs d'effectif n'influe pas
		{[]string{"10", "9", "5", "", ""}, false},    // ❌ l'absence des 2 dernières valeurs d'effectif n'influe pas
		{[]string{"10", "9", "", "", ""}, false},     // ❌ l'absence des 3 dernières valeurs d'effectif n'influe pas
	}

	for i, tc := range testCases {
		t.Run("Test case "+strconv.Itoa(i), func(t *testing.T) {
			shouldKeep := isInsidePerimeter(tc.input, nbMois, minEffectif)
			assert.Equal(t, tc.expected, shouldKeep)
		})
	}
}

func TestGuessLastNonMissing(t *testing.T) {
	testCases := []struct {
		inputCsv string
		expected int
	}{
		{"h;h\n1;", 1},
		{"h;h\n;1", 0},
		{"h;h\n1;1", 0},
		{"h;h\n;", 2},
		{"h;h\n;\n;1", 0},
		{"h;h\n1;\n;", 1},
		{"h;h\n1;\n1;", 1},
	}

	for i, tc := range testCases {
		t.Run("Test case without ignored "+strconv.Itoa(i), func(t *testing.T) {
			mockFile := engine.NewMockBatchFile(tc.inputCsv)
			lastNonMissing, err := guessLastNMissing(mockFile, 0)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, lastNonMissing)
		})
	}

	testCasesIgnore := []struct {
		inputCsv string
		expected int
	}{
		{"h;h;h\n1;;1", 1},
		{"h;h;h\n;1;1", 0},
		{"h;h;h\n1;1;1", 0},
		{"h;h;h\n;;1", 2},
		{"h;h;h\n;;1\n;1;1", 0},
		{"h;h;h\n1;;1\n;;1", 1},
		{"h;h;h\n1;;1\n1;;1", 1},
	}

	for i, tc := range testCasesIgnore {
		t.Run("Test case without ignored "+strconv.Itoa(i), func(t *testing.T) {
			mockFile := engine.NewMockBatchFile(tc.inputCsv)
			lastNonMissing, err := guessLastNMissing(mockFile, 1)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, lastNonMissing)
		})
	}
}

func TestCheck(t *testing.T) {
	// Helper to read test effectif data
	effectifData := readTestData(t, "testData/test_data.csv")

	// Mock readers
	validFilterReader := &MemoryFilterReader{Filter: engine.NoFilter}
	invalidFilterReader := &MemoryFilterReader{Filter: nil}

	var nilInterfaceReader Reader
	var nilPointerReader *StandardReader

	testCases := []struct {
		name         string
		batchFiles   engine.BatchFiles
		filterReader Reader
		expectError  bool
	}{
		{
			"Filtre valide explicitement fourni par l'utilisateur -> OK",
			engine.BatchFiles{
				"filter": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			validFilterReader,
			false,
		},
		{
			"Fichier effectif valide -> OK",
			engine.BatchFiles{
				"effectif_ent": []engine.BatchFile{engine.NewMockBatchFile(effectifData)},
			},
			invalidFilterReader,
			false,
		},
		{
			"Pas de fichier filtre ou effectif ou filtre en base -> NOK",
			engine.BatchFiles{
				"debits": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			invalidFilterReader,
			true,
		},
		{
			"Pas de fichier filtre ou effectif mais filtre en base -> OK",
			engine.BatchFiles{
				"debits": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			validFilterReader,
			false,
		},
		{
			// Si r = nil.(*Reader), r.Read() retourne NoFilter
			"Pointeur de Reader nil -> OK",
			engine.BatchFiles{
				"debits": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			nilPointerReader,
			false,
		},
		{
			// Si r = nil.(engine.FilterReader), r.Read() est illicite
			"Pointeur d'interface nil -> NOK",
			engine.BatchFiles{
				"debits": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			nilInterfaceReader,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Check(tc.filterReader, tc.batchFiles)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateState(t *testing.T) {
	// Helper to read test effectif data
	effectifData := readTestData(t, "testData/test_data.csv")

	testCases := []struct {
		name        string
		batchFiles  engine.BatchFiles
		expectWrite bool
	}{
		{
			"Effectif file present -> filter should be written",
			engine.BatchFiles{
				"effectif_ent": []engine.BatchFile{engine.NewMockBatchFile(effectifData)},
			},
			true,
		},
		{
			"Explicit filter file present -> filter should NOT be written",
			engine.BatchFiles{
				"filter": []engine.BatchFile{engine.NewMockBatchFile("siren\n012345678")},
			},
			false,
		},
		{
			"No effectif or filter file -> filter should NOT be written",
			engine.BatchFiles{
				"debits": []engine.BatchFile{engine.NewMockBatchFile("")},
			},
			false,
		},
		{
			"Both effectif and filter files present -> filter should NOT be written",
			engine.BatchFiles{
				"effectif_ent": []engine.BatchFile{engine.NewMockBatchFile(effectifData)},
				"filter":       []engine.BatchFile{engine.NewMockBatchFile("siren\n012345678")},
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFilterWriter := &MemoryFilterWriter{}
			err := UpdateState(mockFilterWriter, tc.batchFiles)
			assert.NoError(t, err)

			if tc.expectWrite {
				assert.NotNil(t, mockFilterWriter.Filter, "Expected filter to be written")
			} else {
				assert.Nil(t, mockFilterWriter.Filter, "Expected filter NOT to be written")
			}
		})
	}
}

// readTestData reads test data from a file
func readTestData(t *testing.T, filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

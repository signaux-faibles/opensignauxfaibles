package marshal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/stretchr/testify/assert"
)

// MockComptesMapping ...
func MockComptesMapping(mapping map[string]string) Comptes {

	mockComptes := make(Comptes)
	MakeSiretDateArray := func(siret string) []SiretDate {
		longAgo, _ := time.Parse("2006-01-02", "9999-01-02")
		return []SiretDate{
			{
				Siret: siret,
				Date:  longAgo,
			},
		}
	}
	for compte, siret := range mapping {
		mockComptes[compte] = MakeSiretDateArray(siret)
	}
	return mockComptes
}

type tuplesAndEvents = struct {
	Tuples []Tuple `json:"tuples"`
	Events []Event `json:"events"`
}

// GetFatalError retourne le message d'erreur fatale obtenu suite à une
// opération de parsing, ou une chaine vide.
func GetFatalError(output tuplesAndEvents) string {
	headFatal := GetFatalErrors(output.Events[0])
	if headFatal == nil || len(headFatal) < 1 {
		return ""
	}
	if len(headFatal) > 1 {
		log.Println(headFatal)
		log.Fatal("headFatal should never contain more than one item")
	}
	return headFatal[0].(string)
}

// GetFatalErrors retourne les messages d'erreurs fatales obtenus suite à une
// opération de parsing, ou nil.
func GetFatalErrors(event Event) []interface{} {
	reportData, _ := event.ParseReport()
	headFatal, ok := reportData["headFatal"].([]interface{})
	if ok != true {
		return nil
	}
	return headFatal
}

// RunParserInline returns Tuples and Events resulting from the execution of a
// Parser on a given list of rows, with an empty Cache.
func RunParserInline(t *testing.T, parser Parser, rows []string) (output tuplesAndEvents) {
	return RunParserInlineEx(t, NewCache(), parser, rows)
}

// RunParserInlineEx returns Tuples and Events resulting from the execution of a
// Parser on a given list of rows.
func RunParserInlineEx(t *testing.T, cache Cache, parser Parser, rows []string) (output tuplesAndEvents) {
	csvData := strings.Join(rows, "\n")
	csvFile := CreateTempFileWithContent(t, []byte(csvData)) // will clean up after the test
	return RunParser(parser, cache, csvFile.Name())
}

// RunParser returns Tuples and Events resulting from the execution of a
// Parser on a given input file.
func RunParser(
	parser Parser,
	cache Cache,
	inputFile string,
) (output tuplesAndEvents) {
	batch := base.MockBatch(parser.GetFileType(), []string{inputFile})
	tuples, events := ParseFilesFromBatch(cache, &batch, parser)

	// intercepter et afficher les évènements pendant l'importation
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			event.Date = time.Time{}
			output.Events = append(output.Events, event)
		}
	}()

	for tuple := range tuples {
		output.Tuples = append(output.Tuples, tuple)
	}

	wg.Wait()
	return output
}

// TestParserOutput compares output Tuples and output Events with JSON stored
// in a golden file. If update = true, the the golden file is updated.
func TestParserOutput(
	t *testing.T,
	parser Parser,
	cache Cache,
	inputFile string,
	goldenFile string,
	update bool,
) {
	var output = RunParser(parser, cache, inputFile)

	actual, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if update {
		ioutil.WriteFile(goldenFile, []byte(actual), 0644)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal("Could not open golden file" + err.Error())
	}

	assert.Equal(t, string(expected), string(actual))
}

// CreateTempFileWithContent créée un fichier temporaire et le supprime
// après le passage (ou échec) du test.
func CreateTempFileWithContent(t *testing.T, content []byte) *os.File {
	t.Helper()
	tmpfile, err := ioutil.TempFile("", "createTempFileWithContent")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile
}

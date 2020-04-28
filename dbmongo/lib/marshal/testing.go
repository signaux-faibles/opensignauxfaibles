package marshal

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/stretchr/testify/assert"
)

func compareFields(field1 Field, field2 Field) string {
	var res string
	if field1.GoName != field2.GoName {
		res = res + "Different struct fields GoName; " + field1.GoName + ", " + field2.GoName
	}
	if field1.CSVName != field2.CSVName {
		res = res + "Different struct fields CSVName; " + field1.CSVName + ", " + field2.CSVName
	}
	if field1.CSVCol != field2.CSVCol {
		res = res + "Different struct fields CSVCol; " + strconv.Itoa(field1.CSVCol) + ", " + strconv.Itoa(field2.CSVCol)
	}
	if field1.JSONName != field2.JSONName {
		res = res + "Different struct fields JSONName; " + field1.JSONName + ", " + field2.JSONName
	}
	if field1.Parser != field2.Parser {
		res = res + "Different struct fields Parser; " + field1.Parser + ", " + field2.Parser
	}
	if field1.IfEmpty != field2.IfEmpty {
		res = res + "Different struct fields IfEmpty; " + field1.IfEmpty + ", " + field2.IfEmpty
	}
	if field1.IfInvalid != field2.IfInvalid {
		res = res + "Different struct fields IfInvalid; " + field1.IfInvalid + ", " + field2.IfInvalid
	}
	if field1.ValidityRegex != field2.ValidityRegex {
		res = res + "Different struct fields ValidityRegex; " + field1.ValidityRegex + ", " + field2.ValidityRegex
	}
	return res
}

// MockComptesMapping ...
func MockComptesMapping(mapping map[string]string) Comptes {

	mockComptes := make(Comptes)
	MakeSiretDateArray := func(siret string) []SiretDate {
		longAgo, _ := time.Parse("2006-01-02", "9999-01-02")
		return []SiretDate{
			SiretDate{
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

// TestParserTupleOutput helps to test the output of a Parser. It compares
// output Tuples with JSON stored in a golden file. If update = true, the
// the golden file is updated.
func TestParserTupleOutput(
	t *testing.T,
	parser engine.Parser,
	cache engine.Cache,
	parserType string,
	inputFile string,
	goldenFile string,
	update bool,
) {
	batch := engine.MockBatch(parserType, []string{inputFile})
	var events chan engine.Event
	var tuples chan engine.Tuple
	tuples, events = parser(cache, &batch)
	var firstCriticalEvent *engine.Event = nil

	// intercepter et afficher les évènements pendant l'importation
	go func() {
		for event := range events {
			t.Logf("[%s] event: %v", event.Priority, event.Comment)
			if event.Priority == engine.Critical && firstCriticalEvent == nil {
				firstCriticalEvent = &event
			}
		}
	}()

	actualJsons := []string{}
	for tuple := range tuples {
		t.Log("tuple:", tuple)
		json, err := engine.GetJson(tuple)
		if err != nil {
			log.Fatal(err)
		}
		actualJsons = append(actualJsons, string(json))
	}

	actual := "[" + strings.Join(actualJsons, ",") + "]"
	if update {
		ioutil.WriteFile(goldenFile, []byte(actual), 0644)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal("Could not open golden file" + err.Error())
	}

	if firstCriticalEvent != nil {
		assert.FailNow(t, "Caught Critical event: ", firstCriticalEvent.Comment)
	}

	assert.Equal(t, string(expected), string(actual))

}

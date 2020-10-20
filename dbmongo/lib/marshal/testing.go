package marshal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
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

// RunParser returns Tuples and Events resulting from the execution of a
// Parser on a given input file.
func RunParser(
	parser BatchParser,
	cache Cache,
	parserType string,
	inputFile string,
) (output tuplesAndEvents) {
	batch := base.MockBatch(parserType, []string{inputFile})
	var events chan Event
	var tuples chan Tuple
	tuples, events = parser(cache, &batch)

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
	parser BatchParser,
	cache Cache,
	parserType string,
	inputFile string,
	goldenFile string,
	update bool,
) {
	var output = RunParser(parser, cache, parserType, inputFile)

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

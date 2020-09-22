package marshal

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/engine"
	"github.com/stretchr/testify/assert"
)

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
	cache base.Cache,
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
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait() // pour éviter "panic: Log in goroutine after TestEffectif has completed"
	go func() {
		defer wg.Done()
		for event := range events {
			t.Logf("[%s] event: %v", event.Priority, event.Comment)
			if event.Priority == engine.Critical && firstCriticalEvent == nil {
				firstCriticalEvent = &event
			}
		}
	}()

	actualJsons := []string{}
	for tuple := range tuples {
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

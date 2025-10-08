package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"opensignauxfaibles/lib/base"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// DiscardTuple ignore les données
func DiscardTuple(tuples chan Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

// FailDataSink is a sink that always fails
type FailDataSink struct{}

func (s *FailDataSink) ProcessOutput(ctx context.Context, ch chan Tuple) error {
	time.Sleep(500 * time.Millisecond)
	return fmt.Errorf("this sink always fails")
}

type FailSinkFactory struct{}

func (FailSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	return &FailDataSink{}, nil
}

type DiscardReportSink struct {
	counter int
}

func (s DiscardReportSink) Process(ch chan Report) error {
	for range ch {
		s.counter++
	}
	return nil
}

type tuplesAndReports = struct {
	Tuples  []Tuple  `json:"tuples"`
	Reports []Report `json:"reports"`
}

// GetFatalError retourne le message d'erreur fatale obtenu suite à une
// opération de parsing, ou une chaine vide.
func GetFatalError(output tuplesAndReports) string {
	headFatal := GetFatalErrors(output.Reports[0])
	if headFatal == nil || len(headFatal) < 1 {
		return ""
	}
	if len(headFatal) > 1 {
		log.Println(headFatal) // pour aider au débogage en cas d'échec du test
		log.Fatal("headFatal should never contain more than one item")
	}
	return headFatal[0]
}

// GetFatalErrors retourne les messages d'erreurs fatales obtenus suite à une
// opération de parsing, ou nil.
func GetFatalErrors(report Report) []string {
	return report.HeadFatal
}

// ConsumeFatalErrors récupère les erreurs fatales depuis un canal d'évènements
func ConsumeFatalErrors(ch chan Report) []string {
	var fatalErrors []string
	for event := range ch {
		headFatal := GetFatalErrors(event)
		fatalErrors = append(fatalErrors, headFatal...)
	}
	return fatalErrors
}

// RunParserInline returns Tuples and Reports resulting from the execution of a
// Parser on a given list of rows, with an empty Cache.
func RunParserInline(t *testing.T, parser Parser, rows []string) (output tuplesAndReports) {
	return RunParserInlineEx(t, NewEmptyCache(), parser, rows)
}

// RunParserInlineEx returns Tuples and Reports resulting from the execution of a
// Parser on a given list of rows.
func RunParserInlineEx(t *testing.T, cache Cache, parser Parser, rows []string) (output tuplesAndReports) {
	csvData := strings.Join(rows, "\n")
	csvFile := CreateTempFileWithContent(t, []byte(csvData)) // will clean up after the test
	return RunParser(parser, cache, base.NewBatchFile(csvFile.Name()))
}

// TestParserOutput compares output Tuples and output Reports with JSON stored
// in a golden file. If update = true, the the golden file is updated.
func TestParserOutput(
	t *testing.T,
	parser Parser,
	cache Cache,
	inputFile base.BatchFile,
	goldenFile string,
	update bool,
) {
	var output = RunParser(parser, cache, inputFile)

	actual, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if update {
		_ = os.WriteFile(goldenFile, actual, 0644)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal("Could not open golden file" + err.Error())
	}

	assert.Equal(t, string(expected), string(actual))
}

// RunParser returns Tuples and Reports resulting from the execution of a
// Parser on a given input file.
// TimeStamps are set to 0 for reproducibility
func RunParser(
	parser Parser,
	cache Cache,
	inputFile base.BatchFile,
) (output tuplesAndReports) {
	ctx := context.Background()
	batch := base.MockBatch(parser.Type(), []base.BatchFile{inputFile})
	tuples, events := ParseFilesFromBatch(ctx, cache, &batch, parser, NoFilter{})

	// intercepter et afficher les évènements pendant l'importation
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range events {
			event.StartDate = time.Time{}
			output.Reports = append(output.Reports, event)
		}
	}()

	for tuple := range tuples {
		output.Tuples = append(output.Tuples, tuple)
	}

	wg.Wait()
	return output
}

// CreateTempFileWithContent créée un fichier temporaire et le supprime
// après le passage (ou échec) du test.
func CreateTempFileWithContent(t *testing.T, content []byte) *os.File {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "createTempFileWithContent")
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

type TestTuple struct {
	Test1 string `csv:"test1" sql:"test1"`
	Test2 *int   `csv:"test2" sql:"test2"`
	Test3 string
	Test4 *time.Time `csv:"test4" sql:"test4"`
}

func (TestTuple) Key() string           { return "" }
func (TestTuple) Scope() string         { return "" }
func (TestTuple) Type() base.ParserType { return "" }

// -----------------------------------------------------
// Test DataSinkFactory and Data Sink implementations
// -----------------------------------------------------

type TestSinkFactory struct{}

func (TestSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	return &DiscardDataSink{}, nil
}

// DiscardSinkFactory discards all data, regardless of the parser
type DiscardSinkFactory struct{}

func (f *DiscardSinkFactory) CreateSink(parserType base.ParserType) (DataSink, error) {
	return &DiscardDataSink{}, nil
}

type DiscardDataSink struct {
	Counter int
}

func (s *DiscardDataSink) ProcessOutput(ctx context.Context, ch chan Tuple) error {
	for range ch {
		s.Counter++
	}
	return nil
}

// -----------------------------------------------------
// Test PacakgeRegistry mock implementation
// -----------------------------------------------------

type EmptyRegistry struct{}

func (r EmptyRegistry) Resolve(base.ParserType) Parser { return nil }
func (r EmptyRegistry) All() []Parser                  { return []Parser{} }

// -----------------------------------------------------
// Test Filter skips nothing
// -----------------------------------------------------

type NoFilter struct{}

func (f NoFilter) ShouldSkip(string) bool { return false }

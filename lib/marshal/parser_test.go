package marshal

import (
	"errors"
	"log"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/stretchr/testify/assert"
)

func TestParseFilesFromBatch(t *testing.T) {

	t.Run("intérrompt le parsing en cas d'erreur d'initialisation", func(t *testing.T) {
		batch := base.AdminBatch{Files: base.BatchFiles{"dummy": {"dummy.csv"}}}
		_, eventChan := ParseFilesFromBatch(NewCache(), &batch, &dummyParser{initError: errors.New("error from Init()")})
		fatalErrors := ConsumeFatalErrors(eventChan)
		assert.Equal(t, []string{"Fatal: error from Init()"}, fatalErrors)
	})

	t.Run("ne rapporte pas d'erreur de fermeture en cas d'erreur d'ouverture", func(t *testing.T) {
		batch := base.AdminBatch{Files: base.BatchFiles{"dummy": {"dummy.csv"}}}
		_, eventChan := ParseFilesFromBatch(NewCache(), &batch, &dummyParser{})
		fatalErrors := ConsumeFatalErrors(eventChan)
		assert.Equal(t, []string{"Fatal: error from Open()"}, fatalErrors)
	})
}

func TestParseTuplesFromLine(t *testing.T) {

	t.Run("ne comptabilise pas plus d'une erreur par ligne", func(t *testing.T) {
		var parsedLine ParsedLineResult
		parsedLine.AddRegularError(errors.New("error 1"))
		parsedLine.AddRegularError(errors.New("error 2"))
		parsedLine.AddTuple(dummyTuple{})

		tracker := NewParsingTracker()
		parseTuplesFromLine(parsedLine, &SirenFilter{}, &tracker, make(chan Tuple))
		tracker.Next()

		report := tracker.Report("", "")
		assert.Equal(t, int64(1), report["linesParsed"])
		assert.Equal(t, int64(1), report["linesRejected"])
		assert.Equal(t, int64(0), report["linesSkipped"])
		assert.Equal(t, int64(0), report["linesValid"])
		assert.Equal(t, 2, len(report["headRejected"].([]string)))
	})

	t.Run("ne comptabilise qu'une fois une ligne à la fois erronée et filtrée", func(t *testing.T) {
		var parsedLine ParsedLineResult
		parsedLine.AddRegularError(errors.New("regular error"))
		parsedLine.AddTuple(dummyTuple{})
		parsedLine.SetFilterError(errors.New("filtered"))

		tracker := NewParsingTracker()
		parseTuplesFromLine(parsedLine, &SirenFilter{}, &tracker, make(chan Tuple))
		tracker.Next()

		report := tracker.Report("", "")
		assert.Equal(t, int64(1), report["linesParsed"])
		assert.Equal(t, int64(0), report["linesRejected"])
		assert.Equal(t, int64(1), report["linesSkipped"])
		assert.Equal(t, int64(0), report["linesValid"])
		assert.Equal(t, 0, len(report["headRejected"].([]string)))
	})
}

type dummyTuple struct{}

// Key id de l'objet",
func (sirene dummyTuple) Key() string {
	return ""
}

// Type de données
func (sirene dummyTuple) Type() string {
	return "dummy"
}

// Scope de l'objet
func (sirene dummyTuple) Scope() string {
	return "etablissement"
}

type dummyParser struct {
	initError error
}

func (parser *dummyParser) GetFileType() string {
	return "dummy"
}

func (parser *dummyParser) Init(cache *Cache, batch *base.AdminBatch) error {
	return parser.initError
}

func (parser *dummyParser) Open(filePath string) (err error) {
	log.Println("opening", filePath) // TODO: supprimer cet affichage ?
	return errors.New("error from Open()")
}

func (parser *dummyParser) Close() error {
	return errors.New("error from Close()")
}

func (parser *dummyParser) ParseLines(parsedLineChan chan ParsedLineResult) {}

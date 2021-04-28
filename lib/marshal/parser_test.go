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

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
		var fatalErrors []string
		var batch = base.AdminBatch{Files: base.BatchFiles{"dummy": {"dummy.csv"}}}
		_, eventChan := ParseFilesFromBatch(NewCache(), &batch, &dummyParser{})
		// récupération des erreurs
		for event := range eventChan {
			headFatal := GetFatalErrors(event)
			for _, fatalError := range headFatal {
				fatalErrors = append(fatalErrors, fatalError.(string))
			}
		}
		// log.Println(fatalErrors)
		assert.Equal(t, 1, len(fatalErrors))
		assert.Equal(t, "Fatal: error from Init()", fatalErrors[0])
	})
}

// TODO: meme test en cas d'erreur d'ouverture et fermeture

type dummyParser struct{}

func (parser *dummyParser) GetFileType() string {
	return "dummy"
}

func (parser *dummyParser) Init(cache *Cache, batch *base.AdminBatch) error {
	return errors.New("error from Init()")
}

func (parser *dummyParser) Open(filePath string) (err error) {
	log.Println("opening", filePath)
	return errors.New("error from Open()")
}

func (parser *dummyParser) Close() error {
	return errors.New("error from Close()")
}

func (parser *dummyParser) ParseLines(parsedLineChan chan ParsedLineResult) {}

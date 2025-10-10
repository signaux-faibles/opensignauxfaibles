package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"opensignauxfaibles/lib/base"
)

func TestParseFilesFromBatch(t *testing.T) {
	ctx := context.Background()

	t.Run("interrompt le parsing en cas d'erreur d'initialisation", func(t *testing.T) {
		batch := base.MockBatch("dummy", []base.BatchFile{base.NewMockBatchFile("")})
		_, reportChan := ParseFilesFromBatch(
			ctx,
			NewEmptyCache(),
			&batch,
			&dummyParser{initError: errors.New("error from Init()")},
			NoFilter,
		)
		fatalErrors := ConsumeFatalErrors(reportChan)
		assert.Equal(t, []string{"Fatal: error from Init()"}, fatalErrors)
	})

	t.Run("ne rapporte pas d'erreur de fermeture en cas d'erreur d'ouverture", func(t *testing.T) {
		batch := base.MockBatch("dummy", []base.BatchFile{base.OpenFailsBatchFile{}})
		_, eventChan := ParseFilesFromBatch(ctx, NewEmptyCache(), &batch, &dummyParser{}, NoFilter)
		fatalErrors := ConsumeFatalErrors(eventChan)
		assert.Equal(t, []string{"Fatal: error from Open()"}, fatalErrors)
	})
}

func TestParseTuplesFromLine(t *testing.T) {
	ctx := context.Background()

	t.Run("ne comptabilise pas plus d'une erreur par ligne", func(t *testing.T) {
		var parsedLine ParsedLineResult
		parsedLine.AddRegularError(errors.New("error 1"))
		parsedLine.AddRegularError(errors.New("error 2"))
		parsedLine.AddTuple(dummyTuple{})

		tracker := NewParsingTracker()
		processTuplesFromLine(ctx, parsedLine, NoFilter, &tracker, make(chan Tuple))
		tracker.Next()

		report := tracker.Report("", "", "")
		assert.Equal(t, int64(1), report.LinesParsed)
		assert.Equal(t, int64(1), report.LinesRejected)
		assert.Equal(t, int64(0), report.LinesSkipped)
		assert.Equal(t, int64(0), report.LinesValid)
		assert.Equal(t, 2, len(report.HeadRejected))
	})

	t.Run("ne comptabilise qu'une fois une ligne à la fois erronée et filtrée", func(t *testing.T) {
		var parsedLine ParsedLineResult
		parsedLine.AddRegularError(errors.New("regular error"))
		parsedLine.AddTuple(dummyTuple{})
		parsedLine.SetFilterError(errors.New("filtered"))

		tracker := NewParsingTracker()
		processTuplesFromLine(ctx, parsedLine, NoFilter, &tracker, make(chan Tuple))
		tracker.Next()

		report := tracker.Report("", "", "")
		assert.Equal(t, int64(1), report.LinesParsed)
		assert.Equal(t, int64(0), report.LinesRejected)
		assert.Equal(t, int64(1), report.LinesSkipped)
		assert.Equal(t, int64(0), report.LinesValid)
		assert.Equal(t, 0, len(report.HeadRejected))
	})
}

type dummyTuple struct{}

func (sirene dummyTuple) Headers() []string {
	//TODO implement me
	panic("implement me")
}

func (sirene dummyTuple) Values() []string {
	//TODO implement me
	panic("implement me")
}

// Key id de l'objet",
func (sirene dummyTuple) Key() string {
	return ""
}

// Type de données
func (sirene dummyTuple) Type() base.ParserType {
	return "dummy"
}

// Scope de l'objet
func (sirene dummyTuple) Scope() string {
	return "etablissement"
}

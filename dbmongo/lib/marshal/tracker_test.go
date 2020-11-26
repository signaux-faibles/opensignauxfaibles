package marshal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracker(t *testing.T) {
	t.Run("should not keep more than `MaxParsingErrors` parse errors in memory", func(t *testing.T) {
		tracker := NewParsingTracker("", "")
		for i := 0; i < MaxParsingErrors+1; i++ {
			tracker.AddParseError(fmt.Errorf("parse error %d", i))
		}
		assert.Equal(t, MaxParsingErrors, len(tracker.firstParseErrors))
	})

	t.Run("can report more than `MaxParsingErrors` parse errors", func(t *testing.T) {
		expectedLinesRejected := MaxParsingErrors + 1
		tracker := NewParsingTracker("", "")
		for i := 0; i < expectedLinesRejected; i++ {
			tracker.AddParseError(fmt.Errorf("parse error %d", i))
			tracker.Next()
		}
		report := tracker.Report("abstract")
		assert.Equal(t, expectedLinesRejected, report["linesRejected"])
	})

	t.Run("should report just 1 rejected line, given 2 parse errors happened on that same line", func(t *testing.T) {
		expectedLinesRejected := 1
		errorsOnSameLine := 2
		tracker := NewParsingTracker("", "")
		for i := 0; i < errorsOnSameLine; i++ {
			tracker.AddParseError(fmt.Errorf("parse error %d", i))
		}
		report := tracker.Report("abstract")
		assert.Equal(t, expectedLinesRejected, report["linesRejected"])
	})

	t.Run("should report just 1 skipped line, given 2 filter errors happened on that same line", func(t *testing.T) {
		expectedLinesSkipped := 1
		errorsOnSameLine := 2
		tracker := NewParsingTracker("", "")
		for i := 0; i < errorsOnSameLine; i++ {
			tracker.AddFilterError(fmt.Errorf("filter error %d", i))
		}
		report := tracker.Report("abstract")
		assert.Equal(t, expectedLinesSkipped, report["linesSkipped"])
	})

}

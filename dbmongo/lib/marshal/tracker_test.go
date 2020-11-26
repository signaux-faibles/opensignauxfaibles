package marshal

import (
	"fmt"
	"testing"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/stretchr/testify/assert"
)

func TestTracker(t *testing.T) {
	t.Run("should not keep more than `MaxParsingErrors` parse errors in memory", func(t *testing.T) {
		tracker := NewParsingTracker("", "")
		for i := 0; i < MaxParsingErrors+1; i++ {
			tracker.Add(base.NewRegularError(fmt.Errorf("parse error %d", i)))
		}
		assert.Equal(t, MaxParsingErrors, len(tracker.parseErrors))
	})
}

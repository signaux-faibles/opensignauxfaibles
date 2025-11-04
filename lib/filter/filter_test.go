package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterReader(t *testing.T) {

	t.Run("Using a nil reader is possible and results in no filtering", func(t *testing.T) {
		var r Reader
		filter, err := r.Read()
		assert.NoError(t, err, "Reading a nil filter.Reader should succeed")
		assert.Nil(t, filter.All(), "Reading a nil filter.Reader should return a NoFilter")
	})
}

package prepareimport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchKey(t *testing.T) {

	t.Run("Should accept valid batch key", func(t *testing.T) {
		_, err := NewBatchKey("1802")
		assert.NoError(t, err)
	})

	t.Run("Should fail if batch key is invalid", func(t *testing.T) {
		_, err := NewBatchKey("")
		assert.EqualError(t, err, "la clé du batch doit respecter le format requis AAMM")
	})

	t.Run("Should return the path of a batch", func(t *testing.T) {
		batchKey, _ := NewBatchKey("1802")
		assert.Equal(t, "/1802/", batchKey.Path())
	})

	t.Run("Should return the parent of a sub-batch", func(t *testing.T) {
		batchKey, _ := NewBatchKey("1802_01")
		assert.Equal(t, "1802", batchKey.GetParentBatch())
	})
}

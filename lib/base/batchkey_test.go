package base

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
		testCases := []string{"", "18022"}

		for _, key := range testCases {
			_, err := NewBatchKey(key)
			assert.Error(t, err)
		}
	})
}

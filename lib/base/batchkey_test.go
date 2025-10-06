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

func Test_IsBatchID(t *testing.T) {
	if !IsBatchKey("1801") {
		t.Error("1801 devrait être un ID de batch")
	}

	if IsBatchKey("") {
		t.Error("'' ne devrait pas être considéré comme un ID de batch")
	}

	if IsBatchKey("190193039") {
		t.Error("'190193039' ne devrait pas être considéré comme un ID de batch")
	}
	if !IsBatchKey("1901_93039") {
		t.Error("'1901_93039'  devrait être considéré comme un ID de batch")
	}

	if IsBatchKey("abcd") {
		t.Error("'abcd' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'abcd' est bien rejeté: ")
	}
}

package base

import (
	"fmt"
	"regexp"
	"time"
)

type BatchKey string

// NewBatchKey constructs a valid batch key
func NewBatchKey(key string) (BatchKey, error) {
	if !validBatchKey.MatchString(key) {
		return "", fmt.Errorf("la clé du batch doit respecter le format requis AAMM. Reçu : %s", key)
	}
	return BatchKey(key), nil
}

var validBatchKey = regexp.MustCompile(`^[0-9]{4}$`)

func (b BatchKey) String() string {
	return string(b)
}

// IsBatchKey retourne `true` si `batchID` est un identifiant de Batch.
func IsBatchKey(batchKey string) bool {
	if len(batchKey) < 4 {
		return false
	}
	_, err := time.Parse("0601", batchKey[0:4])
	if len(batchKey) > 4 && batchKey[4] != '_' {
		return false
	}
	return err == nil
}

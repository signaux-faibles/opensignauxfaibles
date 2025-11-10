package engine

import (
	"fmt"
	"regexp"
	"time"
)

type BatchKey string

// NewBatchKey constructs a valid batch key
func NewBatchKey(key string) (BatchKey, error) {
	if !validBatchKey.MatchString(key) {
		return "", fmt.Errorf("batch key must follow required format YYMM. Received: %s", key)
	}
	return BatchKey(key), nil
}

var validBatchKey = regexp.MustCompile(`^[0-9]{4}$`)

func (b BatchKey) String() string {
	return string(b)
}

// IsBatchKey returns `true` if `batchID` is a Batch identifier.
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
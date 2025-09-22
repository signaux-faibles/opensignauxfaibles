package base

import (
	"errors"
	"regexp"
)

type BatchKey string

// NewBatchKey constructs a valid batch key
func NewBatchKey(key string) (BatchKey, error) {
	if !validBatchKey.MatchString(key) {
		return "", errors.New("la clé du batch doit respecter le format requis AAMM. Reçu :" + key)
	}
	return BatchKey(key), nil
}

var validBatchKey = regexp.MustCompile(`^[0-9]{4}$`)

func (b BatchKey) String() string {
	return string(b)
}

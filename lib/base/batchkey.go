package base

import (
	"errors"
	"regexp"
)

type BatchKey string

// NewBatchKey constructs a valid batch key
func NewBatchKey(key string) (BatchKey, error) {
	if !validBatchKey.MatchString(key) && !validSubBatchKey.MatchString(key) {
		return "", errors.New("la clé du batch doit respecter le format requis AAMM ou AAMM_bb (bb étant un numéro de sous-batch libre)")
	}
	return BatchKey(key), nil
}

var validBatchKey = regexp.MustCompile(`^[0-9]{4}`)

var validSubBatchKey = regexp.MustCompile(`^([0-9]{4})_([0-9]{2})$`)

func (b BatchKey) String() string {
	return string(b)
}

func (b BatchKey) Path() string {
	return "/" + string(b) + "/"
}

func (b BatchKey) IsSubBatch() bool {
	return validSubBatchKey.MatchString(string(b))
}

func (b BatchKey) GetParentBatch() string {
	if b.IsSubBatch() {
		return validSubBatchKey.FindStringSubmatch(string(b))[1]
	}
	return b.String()
}

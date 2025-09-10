package prepareimport

import (
	"time"
)

// DateFinEffectif is a date that can be serialized for MongoDB.
type DateFinEffectif interface {
	Date() time.Time
}

// NewDateFinEffectif creates a valid DateFinEffectif.
func NewDateFinEffectif(date time.Time) DateFinEffectif {
	// TODO: validate the value
	return dateFinEffectifType{date}
}

type dateFinEffectifType struct {
	time.Time
}

// Deprecated: now use time.Time
func (dateFinEffectif dateFinEffectifType) Date() time.Time {
	return dateFinEffectif.Time
}

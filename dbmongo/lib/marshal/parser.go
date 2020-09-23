package marshal

import (
	"encoding/json"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *base.AdminBatch) (chan Tuple, chan base.Event)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// GetJson sérialise un tuple au format JSON.
func GetJson(tuple Tuple) ([]byte, error) {
	return json.MarshalIndent(tuple, "", "  ")
}

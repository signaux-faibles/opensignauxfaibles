package marshal

import (
	"encoding/json"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *base.AdminBatch) (chan Tuple, chan Event)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// GetJSON sérialise un tuple au format JSON.
func GetJSON(tuple Tuple) ([]byte, error) {
	return json.MarshalIndent(tuple, "", "  ")
}

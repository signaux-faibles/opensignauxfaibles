package engine

import (
	"github.com/cnf/structhash"

	"opensignauxfaibles/lib/marshal"
)

type BatchKey = string

// Data objet établissement (/entreprise/)
type Data struct {
	Scope string `json:"scope" bson:"scope"`
	Key   string `json:"key" bson:"key"`
	Batch Batch  `json:"batch,omitempty" bson:"batch,omitempty"`
}

// GetMD5 returns a MD5 signature of the Tuple
func GetMD5(tuple marshal.Tuple) []byte {
	return structhash.Md5(tuple, 1)
}

// Batch ensemble des données, par type de données
type Batch map[string]marshal.Tuple

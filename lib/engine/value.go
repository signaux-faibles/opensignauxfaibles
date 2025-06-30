package engine

import (
	"github.com/cnf/structhash"

	"opensignauxfaibles/lib/marshal"
)

type BatchKey = string

// Data objet établissement (/entreprise/)
type Data struct {
	Scope string             `json:"scope" bson:"scope"`
	Key   string             `json:"key" bson:"key"`
	Batch map[BatchKey]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// GetMD5 returns a MD5 signature of the Tuple
func GetMD5(tuple marshal.Tuple) []byte {
	return structhash.Md5(tuple, 1)
}

// Batch ensemble des données
// TODO --> le 1e string est le nom du parser c'est celui qui nous intéresse
// TODO --> le 2nd string est un hash
type Batch map[string]map[string]marshal.Tuple

// Merge union de deux objets Batch
func (batch1 Batch) Merge(batch2 Batch) {
	for k := range batch2 {
		if _, ok := batch1[k]; !ok {
			batch1[k] = batch2[k]
		} else {
			for i, j := range batch2[k] {
				batch1[k][i] = j
			}
		}
	}
}

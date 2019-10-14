package engine

import (
	"errors"

	"github.com/cnf/structhash"
	"github.com/globalsign/mgo/bson"
)

// Value structure pour un établissement
type Value struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Data          `json:"value" bson:"value"`
}

// Data objet établissement (/entreprise/)
type Data struct {
	Scope string           `json:"scope" bson:"scope"`
	Key   string           `json:"key" bson:"key"`
	Batch map[string]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// GetMD5 returns a MD5 signature of the Tupe
func GetMD5(tuple Tuple) []byte {
	return structhash.Md5(tuple, 1)
}

// Batch ensemble des données
type Batch map[string]map[string]Tuple

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

// Merge union de deux objets Value
func (value1 Value) Merge(value2 Value) (Value, error) {
	if value1.Value.Key != value2.Value.Key {
		return Value{},
			errors.New("Objets non missibles: clés '" +
				value1.Value.Key + "' et '" +
				value2.Value.Key + "'")
	}
	for idBatch := range value2.Value.Batch {
		if value1.Value.Batch == nil {
			value1.Value.Batch = make(map[string]Batch)
		}
		value1.Value.Batch[idBatch].Merge(value2.Value.Batch[idBatch])
	}
	return value1, nil
}

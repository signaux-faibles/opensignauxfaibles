package engine

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// PurgeBatch permet de supprimer un batch dans les objets de RawData
func PurgeBatch(batchKey string) error {
	functions, err := loadJSFunctions("purgeBatch")
	if err != nil {
		return err
	}
	scope := bson.M{
		"currentBatch": batchKey,
		"f":            functions,
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "RawData"},
		Scope:    scope,
	}
	_, err = Db.DB.C("RawData").Find(nil).MapReduce(job, nil)
	return err
}

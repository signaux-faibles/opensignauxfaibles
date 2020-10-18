package engine

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/spf13/viper"
)

// PurgeBatchOne purge 1 batch pour 1 siren
func PurgeBatchOne(batch base.AdminBatch, key string) error {
	functions, err := loadJSFunctions("purgeBatch")
	if err != nil {
		return err
	}

	MRscope := bson.M{
		"fromBatchKey": batch.ID.Key,
		"f":            functions,
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"merge": "purgeBatch_debug"},
		Scope:    MRscope,
	}

	query := bson.M{
		"_id": bson.M{
			"$regex": bson.RegEx{Pattern: "^" + key[0:9],
				Options: "",
			},
		},
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)
	return err
}

// PurgeBatch permet de supprimer tous les batch consécutifs au un batch donné dans RawData
func PurgeBatch(batch base.AdminBatch) error {
	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return fmt.Errorf("chunkCollection a échoué: %s", err.Error())
	}
	queries := chunks.ToQueries(bson.M{}, "_id")
	spew.Dump(queries)
	// // TODO avant de changer clearTempCollections, vérifier que le nettoyage
	// // fonctionne comme attendu
	// var clearTempCollections = false

	// functions, err := loadJSFunctions("purgeBatch")
	// if err != nil {
	// 	return err
	// }
	// MRscope := bson.M{
	// 	"currentBatch": batchKey,
	// 	"f":            functions,
	// }

	// // Calculs parallélisés pour éviter un Out Of Memory qui se produit
	// // lorsqu'est lancé le map-reduce sur toute la base.
	// chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	// if err != nil {
	// 	return err
	// }

	// w := MRWait{}
	// w.init()

	// var tempDbNames []string
	// var tempDBChannel = make(chan string)

	// numQuery := 0
	// for _, chunkQuery := range chunks.ToQueries(nil, "_id") {
	// 	w.waitGroup.Add(1)

	// 	tempDbName := "purgeBatch" + strconv.Itoa(numQuery)
	// 	job := &mgo.MapReduce{
	// 		Map:      functions["map"].Code,
	// 		Reduce:   functions["reduce"].Code,
	// 		Finalize: functions["finalize"].Code,
	// 		Out:      bson.M{"merge": "RawData"},
	// 		Scope:    MRscope,
	// 	}
	// 	numQuery++

	// 	go MRroutine(
	// 		job,
	// 		chunkQuery,
	// 		tempDbName,
	// 		/*Origin collection name = */ "RawData",
	// 		/*Waitgroup = */ &w,
	// 		tempDBChannel,
	// 	)
	// }

	// // On consomme les objets dans tempDBChannel
	// go func() {
	// 	for tempDbName := range tempDBChannel {
	// 		tempDbNames = append(tempDbNames, tempDbName)
	// 	}
	// }()

	// w.waitGroup.Wait()
	// close(tempDBChannel)

	// db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	// db.SetSocketTimeout(720000 * time.Second)

	// if clearTempCollections {
	// 	for _, tempDbName := range tempDbNames {
	// 		pipeline := []bson.M{
	// 			bson.M{
	// 				"$merge": bson.M{
	// 					"into": bson.M{
	// 						"coll":        "RawData",
	// 						"db":          viper.GetString("DB"),
	// 						"whenMatched": "merge",
	// 					},
	// 				},
	// 			},
	// 		}

	// 		pipe := db.DB(tempDbName).C("TemporaryCollection").Pipe(pipeline)
	// 		var result []interface{}
	// 		err = pipe.AllowDiskUse().All(&result)
	// 		if err != nil {
	// 			w.add("errors", 1, -1)
	// 		} else {
	// 			err = db.DB(tempDbName).DropDatabase()
	// 			if err != nil {
	// 				w.add("errors", 1, -1)
	// 			}
	// 		}
	// 	}
	// }

	// db.Close()

	// errorcount, _ := w.running.Load("errors")
	// if errorcount.(int) != 0 {
	// 	return errors.New("erreurs constatées, consultez les journaux")
	// }

	return nil
}

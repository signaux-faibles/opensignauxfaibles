package engine

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/spf13/viper"
)

// PurgeBatchOne purge 1 batch pour 1 siren
func PurgeBatchOne(batch base.AdminBatch, key string) error {
	jsParams := bson.M{
		"fromBatchKey": batch.ID.Key,
	}
	mapReduceJob, err := MakeMapReduceJob("purgeBatch", jsParams)
	if err != nil {
		return err
	}

	mapReduceJob.Out = bson.M{"merge": "purgeBatch_debug"}

	query := bson.M{
		"_id": bson.M{
			"$regex": bson.RegEx{Pattern: "^" + key[0:9],
				Options: "",
			},
		},
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(mapReduceJob, nil)
	return err
}

// PurgeBatch permet de supprimer tous les batch consécutifs au un batch donné dans RawData
func PurgeBatch(batch base.AdminBatch) error {
	startDate := time.Now()

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return fmt.Errorf("chunkCollection a échoué: %s", err.Error())
	}
	queries := chunks.ToQueries(bson.M{}, "_id")
	queryChan := queriesToChan(queries)

	jsParams := bson.M{
		"fromBatchKey": batch.ID.Key,
	}
	mapReduceJob, err := MakeMapReduceJob("purgeBatch", jsParams)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for id := 0; id < viper.GetInt("MRthreads"); id++ {
		wg.Add(1)
		go MRChunks(queryChan, *mapReduceJob, "purgeBatch", id, &wg)
	}

	wg.Wait()

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for id := 0; id < viper.GetInt("MRthreads"); id++ {
		tempDB := "purgeBatch" + strconv.Itoa(id)
		pipeline := []bson.M{{
			"$merge": bson.M{"into": bson.M{"coll": "RawData", "db": viper.GetString("DB")}}}}
		pipe := db.DB(tempDB).C("TemporaryCollection").Pipe(pipeline)

		err = pipe.AllowDiskUse().All(&[]interface{}{})
		if err != nil {
			fmt.Println("quelque chose vient de se casser: " + err.Error()) // TODO: supprimer cet affichage ?
			return err
		}
		db.DB(tempDB).DropDatabase()
	}

	LogOperationEvent("PurgeBatch", startDate)

	return nil
}

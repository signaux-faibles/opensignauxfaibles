package engine

import (
	"fmt"
	"log"
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
	functions, err := loadJSFunctions("purgeBatch", bson.M{
		"fromBatchKey": batch.ID.Key,
	})
	if err != nil {
		return err
	}

	MRscope := bson.M{
		"f":            functions,
		"fromBatchKey": batch.ID.Key,
	}

	job := &mgo.MapReduce{ // TODO: laisser loadJSFunctions générer cet objet
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

func queriesToChan(queries []bson.M) chan bson.M {
	channel := make(chan bson.M)
	go func() {
		for _, query := range queries {
			channel <- query
		}
		close(channel)
	}()
	return channel
}

// MRChunks exécute un job MapReduce à partir d'un channel fournissant des queries
func MRChunks(queryChan chan bson.M, MRBaseJob mgo.MapReduce, tempDBprefix string, id int, wg *sync.WaitGroup) {
	for query := range queryChan {
		job := MRBaseJob
		job.Out = bson.M{"merge": "TemporaryCollection", "db": tempDBprefix + strconv.Itoa(id)}
		log.Println(tempDBprefix+strconv.Itoa(id)+": ", query)
		_, err := Db.DB.C("RawData").Find(query).MapReduce(&job, nil)
		if err != nil {
			fmt.Println(tempDBprefix+strconv.Itoa(id)+": error ", err.Error())
		}
	}
	wg.Done()
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
	mapReduceJob, err := makeMapReduceJob("purgeBatch", jsParams)
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
			fmt.Println("quelque chose vient de se casser: " + err.Error())
			return err
		}
		db.DB(tempDB).DropDatabase()
	}

	LogOperationEvent("PurgeBatch", startDate)

	return nil
}

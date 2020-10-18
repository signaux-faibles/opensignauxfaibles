package engine

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

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
func MRChunks(queryChan chan bson.M, MRBaseJob mgo.MapReduce, tempDBprefix string, id *int, tempDBs *[]string, idMutex *sync.Mutex, wg *sync.WaitGroup) {
	for query := range queryChan {
		idMutex.Lock()
		i := *id
		*id++
		tempDB := tempDBprefix + strconv.Itoa(i)
		*tempDBs = append(*tempDBs, tempDB)
		idMutex.Unlock()
		job := MRBaseJob
		job.Out = bson.M{"replace": "TemporaryCollection", "db": tempDB}
		log.Println(tempDBprefix+strconv.Itoa(i)+": ", query)
		_, err := Db.DB.C("RawData").Find(query).MapReduce(&job, nil)
		if err != nil {
			fmt.Println(tempDBprefix+strconv.Itoa(i)+": error ", err.Error())
		}
	}
	wg.Done()
}

// PurgeBatch permet de supprimer tous les batch consécutifs au un batch donné dans RawData
func PurgeBatch(batch base.AdminBatch) error {
	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return fmt.Errorf("chunkCollection a échoué: %s", err.Error())
	}
	queries := chunks.ToQueries(bson.M{}, "_id")
	queryChan := queriesToChan(queries)

	functions, err := loadJSFunctions("purgeBatch")
	if err != nil {
		return err
	}

	baseJob := mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Scope: bson.M{
			"fromBatchKey": batch.ID.Key,
			"f":            functions,
		},
	}

	wg := sync.WaitGroup{}
	var id = 0
	var idMutex sync.Mutex
	var tempDBs []string
	for i := 0; i < viper.GetInt("MRthreads"); i++ {
		wg.Add(1)
		go MRChunks(queryChan, baseJob, "purgeBatch", &id, &tempDBs, &idMutex, &wg)
	}

	wg.Wait()

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for _, tempDB := range tempDBs {
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

	return nil
}

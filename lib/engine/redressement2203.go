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

func Redressement2203One(batch base.AdminBatch, dateStr string) error {
	jsParams := bson.M{
		"dateStr": dateStr,
		"dateFin": batch.Params.DateFin,
	}

	mapReduceJob, err := MakeMapReduceJob("redressement2203", jsParams)
	if err != nil {
		return err
	}

	mapReduceJob.Out = bson.M{
		"merge":     "redressement2203_debug",
		"nonAtomic": true,
	}

	query := bson.M{}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(mapReduceJob, nil)
	return err
}

// PurgeBatch permet de supprimer tous les batch consécutifs au un batch donné dans RawData
func Redressement2203(batch base.AdminBatch, dateStr string) error {
	startDate := time.Now()

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return fmt.Errorf("chunkCollection a échoué: %s", err.Error())
	}
	queries := chunks.ToQueries(bson.M{"value.index.algo2": true}, "_id")
	queryChan := queriesToChan(queries)

	jsParams := bson.M{
		"dateStr": dateStr,
		"dateFin": batch.Params.DateFin,
	}
	mapReduceJob, err := MakeMapReduceJob("redressement2203", jsParams)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for id := 0; id < viper.GetInt("MRthreads"); id++ {
		wg.Add(1)
		go MRChunks(queryChan, *mapReduceJob, "redressement2203", id, &wg)
	}

	wg.Wait()

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for id := 0; id < viper.GetInt("MRthreads"); id++ {
		tempDB := "redressement2203" + strconv.Itoa(id)
		pipeline := []bson.M{{
			"$merge": bson.M{"into": bson.M{"coll": "redressement2203", "db": viper.GetString("DB")}}}}
		pipe := db.DB(tempDB).C("TemporaryCollection").Pipe(pipeline)

		err = pipe.AllowDiskUse().All(&[]interface{}{})
		if err != nil {
			fmt.Println("quelque chose vient de se casser: " + err.Error()) // TODO: supprimer cet affichage ?
			return err
		}
		db.DB(tempDB).DropDatabase()
	}

	LogOperationEvent("redressement2203", startDate)

	return nil
}

package engine

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/lib/misc"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// PublicOne traite le mapReduce public pour une clé unique (siren)
func PublicOne(batch base.AdminBatch, key string) error {

	if len(key) < 9 {
		return errors.New("key minimal length of 9")
	}

	jsParams := bson.M{
		"actual_batch":  batch.ID.Key,
		"date_fin":      batch.Params.DateFin,
		"serie_periode": misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
	}
	job, err := makeMapReduceJob("public", jsParams)
	if err != nil {
		return err
	}

	job.Out = bson.M{"replace": "Public_debug"}

	query := bson.M{
		"_id": bson.M{
			"$regex": bson.RegEx{Pattern: "^" + key[0:9],
				Options: "",
			},
		},
	}
	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// Public traite le mapReduce public pour les entreprises et établissements du perimètre "algo2".
func Public(batch base.AdminBatch) error {
	startDate := time.Now()

	jsParams := bson.M{
		"actual_batch":  batch.ID.Key,
		"date_fin":      batch.Params.DateFin,
		"serie_periode": misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
	}
	mapReduceJobTemplate, err := makeMapReduceJob("public", jsParams)
	if err != nil {
		return err
	}

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return err
	}

	w := MRWait{}
	w.init()

	var pipes []string
	var pipeChannel = make(chan string)

	filter := bson.M{} // on prend tous les objets comme on sait qu'ils font partie du filtre.

	i := 0
	for _, query := range chunks.ToQueries(filter, "_id") {
		w.waitGroup.Add(1)

		dbTemp := "purgeBatch" + strconv.Itoa(i)           // TODO: à renommer
		var chunkJob mgo.MapReduce = *mapReduceJobTemplate // on dérive mapReduceJobTemplate par copie
		chunkJob.Out = bson.M{"replace": "TemporaryCollection", "db": dbTemp}
		i++

		go MRroutine(&chunkJob, query, dbTemp, "RawData", &w, pipeChannel)

	}

	go func() {
		for pipeDB := range pipeChannel {
			pipes = append(pipes, pipeDB)
		}
	}()

	w.waitGroup.Wait()
	close(pipeChannel)

	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for _, dbTemp := range pipes {
		pipeline := []bson.M{{
			"$merge": bson.M{"into": bson.M{"coll": "Public", "db": viper.GetString("DB")}}}}
		pipe := db.DB(dbTemp).C("TemporaryCollection").Pipe(pipeline)
		var result []interface{}
		pipe.AllowDiskUse().All(&result)
		db.DB(dbTemp).DropDatabase()
	}

	db.Close()

	LogOperationEvent("Public", startDate)

	return nil
}

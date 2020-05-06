package engine

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// PublicOne traite le mapReduce public pour une clé unique (siren)
func PublicOne(batch AdminBatch, key string) error {

	if len(key) < 9 {
		return errors.New("key minimal length of 9")
	}

	functions, err := loadJSFunctions("public")

	naf, err := naf.LoadNAF()
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":      batch.Params.DateDebut,
		"date_fin":        batch.Params.DateFin,
		"serie_periode":   misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif": (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":    batch.ID.Key,
		"naf":             naf,
		"f":               functions,
		"batches":         GetBatchesID(),
		"types":           GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "Public_debug"},
		Scope:    scope,
	}

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

// Public permet de supprimer un batch dans les objets de RawData
func Public(batch AdminBatch) error {
	functions, err := loadJSFunctions("public")
	if err != nil {
		return err
	}
	scope := bson.M{
		"date_debut":      batch.Params.DateDebut,
		"date_fin":        batch.Params.DateFin,
		"serie_periode":   misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
		"offset_effectif": (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":    batch.ID.Key,
		"naf":             naf.Naf,
		"f":               functions,
		"batches":         GetBatchesID(),
		"types":           GetTypes(),
	}

	chunks, err := ChunkCollection(viper.GetString("DB"), "RawData", viper.GetInt64("chunkByteSize"))
	if err != nil {
		return err
	}

	w := MRWait{}
	w.init()

	var pipes []string
	var pipeChannel = make(chan string)

	i := 0
	for _, query := range chunks.ToQueries(bson.M{"value.index.algo2": true}, "_id") {
		w.waitGroup.Add(1)

		dbTemp := "purgeBatch" + strconv.Itoa(i)
		job := &mgo.MapReduce{
			Map:      functions["map"].Code,
			Reduce:   functions["reduce"].Code,
			Finalize: functions["finalize"].Code,
			Out:      bson.M{"replace": "TemporaryCollection", "db": dbTemp},
			Scope:    scope,
		}
		i++

		go MRroutine(job, query, dbTemp, "RawData", &w, pipeChannel)

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
	return nil
}

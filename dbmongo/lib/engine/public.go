package engine

import (
	"opensignauxfaibles/dbmongo/lib/misc"
	"opensignauxfaibles/dbmongo/lib/naf"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// Public permet de supprimer un batch dans les objets de RawData
func Public(batch AdminBatch) error {
	functions, err := loadJSFunctions("public")
	if err != nil {
		return err
	}
	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf.Naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
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
	for _, query := range chunks.ToQueries(bson.M{"value.index.algo2": true}) {
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

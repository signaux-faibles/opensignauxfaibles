package engine

import (
	"dbmongo/lib/misc"
	"dbmongo/lib/naf"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// PublicMergeAux merges collection reduced into its destination
// func PublicMergeAux() error {
// 	// job := &mgo.MapReduce{
// 	// 	Map:      "function map() { emit(this._id, {info: this.info, value: this.value}) }",
// 	// 	Reduce:   "function reduce(_, v) {return v}",
// 	// 	Finalize: "function finalize(_, v) { return v }",
// 	// 	Out:      bson.M{"merge": "Features"},
// 	// }
// 	// _, err := Db.DB.C("Features_aux").Find(bson.M{}).MapReduce(job, nil)

// 	query := []bson.M{{
// 		"$merge": bson.M{"into": "Public"},
// 	}}
// 	pipe := Db.DB.C("Public_aux").Pipe(query)
// 	resp := []bson.M{}
// 	err := pipe.All(&resp)

// 	if err != nil {
// 		return err
// 	}
// 	_, err = Db.DB.C("Public_aux").RemoveAll(nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// Public alimente la collection Public avec les objets destinés à la diffusion
// func Public(batch AdminBatch) error {

// 	functions, err := loadJSFunctions("public")

// 	scope := bson.M{
// 		"date_debut":             batch.Params.DateDebut,
// 		"date_fin":               batch.Params.DateFin,
// 		"date_fin_effectif":      batch.Params.DateFinEffectif,
// 		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
// 		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateFin.AddDate(0, -24, 0), batch.Params.DateFin),
// 		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
// 		"actual_batch":           batch.ID.Key,
// 		"naf":                    naf.Naf,
// 		"f":                      functions,
// 		"batches":                GetBatchesID(),
// 		"types":                  GetTypes(),
// 	}

// 	job := &mgo.MapReduce{
// 		Map:      functions["map"].Code,
// 		Reduce:   functions["reduce"].Code,
// 		Finalize: functions["finalize"].Code,
// 		Out:      bson.M{"replace": collection},
// 		Scope:    scope,
// 	}
// 	// exécution

// 	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

// 	if err != nil {
// 		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
// 	}
// 	return nil
// }

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

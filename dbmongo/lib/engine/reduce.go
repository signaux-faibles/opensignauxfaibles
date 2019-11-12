package engine

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// ReduceOne lance le calcul de Features pour la clé passée en argument
func ReduceOne(batch AdminBatch, algo string, key string) error {
	// éviter les noms d'algo essayant de hacker l'exploration des fonctions ci-dessous
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	if len(key) < 9 {
		return errors.New("key minimal length of 9")
	}

	functions, err := loadJSFunctions("reduce." + algo)

	naf, err := naf.LoadNAF()
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "TemporaryCollection"},
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
	pipeline := []bson.M{
		bson.M{
			"$unwind": bson.M{
				"path": "$value",
				"preserveNullAndEmptyArrays": false,
			},
		},
		bson.M{
			"$match": bson.M{
				"value.effectif": bson.M{
					"$not": bson.M{"$type": 10},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":   0.0,
				"info":  "$_id",
				"value": 1.0,
			},
		},
		bson.M{
			"$merge": bson.M{
				"into": bson.M{
					"coll": "Features_debug",
				},
			},
		},
	}

	pipe := Db.DB.C("TemporaryCollection").Pipe(pipeline)
	var result []interface{}
	err = pipe.AllowDiskUse().All(&result)

	return err
}

// Reduce alimente la base Features
func Reduce(batch AdminBatch, algo string) error {
	// éviter les noms d'algo essayant de hacker l'exploration des fonctions ci-dessous
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	functions, err := loadJSFunctions("reduce." + algo)

	naf, err := naf.LoadNAF()
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
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

	var tempDBs []string
	var tempDBChannel = make(chan string)

	// pour chaque tranche, on lance une goroutine
	// chaque goroutine essaye de lancer un mapreduce
	// voir MRthreads dans le fichier de config
	i := 0
	for _, query := range chunks.ToQueries(bson.M{"value.index." + algo: true}, "_id") {
		w.waitGroup.Add(1)
		dbTemp := "reduce" + strconv.Itoa(i)

		job := &mgo.MapReduce{
			Map:      functions["map"].Code,
			Reduce:   functions["reduce"].Code,
			Finalize: functions["finalize"].Code,
			Out:      bson.M{"replace": "TemporaryCollection", "db": dbTemp},
			Scope:    scope,
		}
		i++
		go MRroutine(job, query, dbTemp, "RawData", &w, tempDBChannel)

	}

	// collecte des noms de bases de données temporaires
	go func() {
		for tempDB := range tempDBChannel {
			tempDBs = append(tempDBs, tempDB)
		}
	}()

	w.waitGroup.Wait()
	close(tempDBChannel)

	// Merge et suppression des bases temporaires
	db, _ := mgo.Dial(viper.GetString("DB_DIAL"))
	db.SetSocketTimeout(720000 * time.Second)

	for _, dbTemp := range tempDBs {
		pipeline := []bson.M{
			bson.M{
				"$unwind": bson.M{
					"path": "$value",
					"preserveNullAndEmptyArrays": false,
				},
			},
			bson.M{
				"$match": bson.M{
					"value.effectif": bson.M{
						"$not": bson.M{"$type": 10},
					},
				},
			},
			bson.M{
				"$project": bson.M{
					"_id":   0.0,
					"info":  "$_id",
					"value": 1.0,
				},
			},
			bson.M{
				"$merge": bson.M{
					"into": bson.M{
						"coll": "Features",
						"db":   viper.GetString("DB"),
					},
				},
			},
		}

		pipe := db.DB(dbTemp).C("TemporaryCollection").Pipe(pipeline)
		var result []interface{}
		err = pipe.AllowDiskUse().All(&result)
		if err != nil {
			w.add("errors", 1, -1)
		} else {
			err = db.DB(dbTemp).DropDatabase()
			if err != nil {
				w.add("errors", 1, -1)
			}
		}
	}

	db.Close()

	errorcount, _ := w.running.Load("errors")
	if errorcount.(int) != 0 {
		return errors.New("erreurs constatées, consultez les journaux")
	}

	return nil
}

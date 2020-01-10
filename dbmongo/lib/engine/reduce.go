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
func ReduceOne(batch AdminBatch, algo string, key, from, to string, types []string) error {

	if len(key) < 9 && (from == "" && to == "") {
		return errors.New("key minimal length of 9")
	}

	scope, err := reduceDefineScope(algo, types)
	if err != nil {
		return err
	}

	var query bson.M
	if key != "" {
		query = bson.M{
			"_id": bson.M{
				"$regex": bson.RegEx{
					Pattern: "^" + key[0:9],
					Options: "",
				},
			},
		}
	} else if from != "" && to != "" {
		query = bson.M{
			"$and": []bson.M{
				bson.M{"_id": bson.M{"$gte": from}},
				bson.M{"_id": bson.M{"$lt": to}},
				query,
			},
		}
	} else {
		return fmt.Errorf("Les paramètres key, ou la paire de paramètres from et to, sont obligatoires")
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "TemporaryCollection"},
		Scope:    scope,
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}
	reduceFinalAggregation(viper.GetString("DB"), "TemporaryDatabase", viper.GetString("DB"), "Features_debug")
	return err
}

// Reduce alimente la base Features
func Reduce(batch AdminBatch, algo string, types []string) error {

	scope, err := reduceDefineScope(algo, types)
	if err != nil {
		return err
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

		reduceFinalAggregation(dbTemp, "TemporaryCollection", viper.GetString("DB"), "Features")

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

func reduceFinalAggregation(tempDatabase, tempCollection, outDatabase, outCollection string) error {
	pipeline := []bson.M{
		bson.M{
			"$unwind": bson.M{
				"path": "$value",
				"preserveNullAndEmptyArrays": false,
			},
		},
		// bson.M{
		// 	"$match": bson.M{
		// 		"value.effectif": bson.M{
		// 			"$not": bson.M{"$type": 10},
		// 		},
		// 	},
		// },
		bson.M{
			"$project": bson.M{
				"_id": bson.D{
					{"batch", "$_id.batch"},
					{"siret", "$value.siret"},
					{"periode", "$_id.periode"},
				},
				"value": 1.0,
			},
		},
		bson.M{
			"$merge": bson.M{
				"into": bson.M{
					"coll": outCollection,
					"db":   outDatabase, //  viper.GetString("DB"),
				},
				"whenMatched": []bson.M{
					bson.M{
						"$project": bson.M{
							"_id": "$_id",
							"value": bson.M{
								"$mergeObjects": []string{
									"$value",
									"$$new.value",
								},
							},
						},
					},
					bson.M{
						"$set": bson.M{
							"value.total": "test",
						},
					},
				},
			},
		},
	}

	pipe := db.DB(tempDatabase).C(tempCollection).Pipe(pipeline)
	var result []interface{}
	err = pipe.AllowDiskUse().All(&result)
	return err
}

func reduceDefineScope(algo string, types []string) (bson.M, error) {

	// Limiter les caractères de nom d'algo pour éviter de hacker la fonction
	// loadJSFunctions
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	if algo == "" {
		return nil, errors.New("Veuillez spécifier un nom d'algo (par exemple avec l'option algo=algo2)")
	}

	functions, err := loadJSFunctions("reduce." + algo)
	if err != nil {
		return nil, err
	}

	naf, err := naf.LoadNAF()
	if err != nil {
		return nil, err
	}

	includes := map[string]bool{}
	for _, data_type := range types {
		includes[data_type] = true
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
		"includes":               includes,
	}
	return scope, nil
}

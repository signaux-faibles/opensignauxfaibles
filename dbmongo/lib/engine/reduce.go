package engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/misc"
	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/naf"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// ReduceOne lance le calcul de Features pour la clé passée en argument
func ReduceOne(batch base.AdminBatch, algo string, key, from, to string, types []string) error {

	if len(key) < 9 && (from == "" && to == "") {
		return errors.New("key minimal length of 9")
	}

	scope, err := reduceDefineScope(batch, algo, types)
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
				{"_id": bson.M{"$gte": from}},
				{"_id": bson.M{"$lt": to}},
				query,
			},
		}
	} else {
		return fmt.Errorf("Les paramètres key, ou la paire de paramètres from et to, sont obligatoires")
	}

	functions := scope["f"].(map[string]bson.JavaScript)
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "TemporaryCollection"},
		Scope:    scope,
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	if err != nil {
		return err
	}
	err = reduceFinalAggregation(Db.DB, "TemporaryCollection", viper.GetString("DB"), "Features_debug")
	return err
}

// Reduce alimente la base Features
func Reduce(batch base.AdminBatch, algo string, types []string) error {

	scope, err := reduceDefineScope(batch, algo, types)
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

		functions := scope["f"].(map[string]bson.JavaScript)
		// Injection des fonctions JavaScript pour exécution par MongoDB
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

	var outCollection string
	if batch.Name != "" {
		outCollection = "Features_" + batch.Name
	} else {
		outCollection = "Features"
	}

	for _, dbTemp := range tempDBs {

		err = reduceFinalAggregation(
			db.DB(dbTemp),
			"TemporaryCollection",
			/*outDatabase = */ viper.GetString("DB"),
			outCollection,
		)

		if err != nil {
			fmt.Println(err)
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

// loadCrossComputationStages charge les étapes d'agrégation MongoDB depuis des fichiers JSON.
func reduceCrossComputations(directoryName string) (stages []bson.M, err error) {
	stages = []bson.M{}
	files, err := listCrossCompJSONFiles(filepath.Join("js", directoryName))
	if err != nil {
		return nil, err
	}
	for _, filePath := range files {
		stage, err := parseJSONObject(filePath)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	return stages, nil
}

func reduceFinalAggregation(tempDatabase *mgo.Database, tempCollection, outDatabase, outCollection string) error {

	// Création d'index sur la collection Features, pour le chargement de données depuis R
	Db.DB.C(outCollection).EnsureIndex(mgo.Index{
		Name: "_id.batch_1_value.random_order_-1__id.periode_1_value.effectif_1", // trouvé sur la db de prod
		Key:  []string{"_id.batch", "-value.random_order", "_id.periode", "value.effectif"},
	})

	setStages, err := reduceCrossComputations("reduce.algo2")
	if err != nil {
		return err
	}

	var pipeline []bson.M
	pipeline = append(pipeline, []bson.M{
		// séparation des données par établissement
		{
			"$unwind": bson.M{
				"path":                       "$value",
				"preserveNullAndEmptyArrays": false,
			},
		},
		// TODO: exclure les données concernant les entreprises non inclues dans le filtre.
		// On voudrait à nouveau exclure tous les établissements qui ont un effectif inconnu ou <10.
		// cf https://github.com/signaux-faibles/opensignauxfaibles/issues/199.
		//
		// Ancienne implémentation:
		// on ne garde que les établissements dont on connait l'effectif (non-null)
		// Commenté parce que dans le cadre de la séparation des calculs par types de données,
		// si on n'intègre pas l'effectif, cette étape filtrerait toutes les données.
		// bson.M{
		// 	"$match": bson.M{
		// 		"value.effectif": bson.M{
		// 			"$not": bson.M{"$type": 10}, // type 10 = null
		// 		},
		// 	},
		// },
		// on a plusieurs objets par clé => on génère un nouvel identifiant et on stocke la clé dans "info"
		{
			"$project": bson.M{
				"_id": bson.D{ // on utilise bson.D pour conserver cet ordre, et permettre la fusion (mergePipeline)
					bson.DocElem{Name: "batch", Value: "$_id.batch"},
					bson.DocElem{Name: "siret", Value: "$value.siret"},
					bson.DocElem{Name: "periode", Value: "$_id.periode"},
				},
				"value": 1.0,
			},
		},
	}...,
	)

	// Defining pipeline used to during merge stage
	mergePipeline := []bson.M{
		{
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
	}
	mergePipeline = append(mergePipeline, setStages...)

	// Merge stage / insertion des données dans la collection Features
	pipeline = append(pipeline,
		bson.M{
			"$merge": bson.M{
				"into": bson.M{
					"coll": outCollection,
					"db":   outDatabase,
				},
				"whenMatched": mergePipeline,
			},
		},
	)

	// Apply aggregation
	pipe := tempDatabase.C(tempCollection).Pipe(pipeline)

	var result []interface{}
	err = pipe.AllowDiskUse().All(&result)
	return err
}

func reduceDefineScope(batch base.AdminBatch, algo string, types []string) (bson.M, error) {

	// Limiter les caractères de nom d'algo pour éviter de hacker la fonction
	// loadJSFunctions
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return nil, errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
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
	if len(types) == 0 {
		includes["all"] = true
	} else {
		for _, dataType := range types {
			includes[dataType] = true
		}
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

// listCrossCompJSONFiles retourne la liste des fichiers JSON présents dans le répertoire spécifié.
func listCrossCompJSONFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && strings.Contains(filePath, ".crossComputation.json") {
			files = append(files, filePath)
		}
		return err
	})
	return files, err
}

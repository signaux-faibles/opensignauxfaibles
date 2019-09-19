package engine

// // Public runs Public MapReduce
// func Public(batch AdminBatch, siret string) error {
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
// 		Out:      bson.M{"replace": "Public"},
// 		Scope:    scope,
// 	}
// 	// exécution

// 	if siret != "" {
// 		_, err = Db.DB.C("RawData").Find(bson.M{
// 			"$or": []interface{}{
// 				bson.M{"_id": siret},
// 				bson.M{"_id": siret[0:9]},
// 			},
// 		}).MapReduce(job, nil)
// 	} else {
// 		// _, err = Db.DB.C("RawData").Find(bson.M{"value.index.algo2": true}).MapReduce(job, nil)
// 		PublicTotal(batch.ID)
// 	}
// 	if err != nil {
// 		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
// 	}
// 	return nil
// }

// // Reduce alimente la base Features
// func PublicTotal(batchKey string, query interface{}, collection string) error {

// 	// éviter les noms d'algo essayant de hacker l'exploration des fonctions ci-dessous
// 	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
// 	if !isAlphaNum(algo) {
// 		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
// 	}

// 	functions, err := loadJSFunctions("reduce." + algo)

// 	naf, err := naf.LoadNAF()
// 	if err != nil {
// 		return err
// 	}

// 	batch, err := GetBatch(batchKey)
// 	if err != nil {
// 		return err
// 	}

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
// 		//TODO merge into collection instead of replacing. Must be idempotent
// 		//transformation. Not the case now with agregation
// 		Out:   bson.M{"replace": collection}, // bson.M{"merge": collection},
// 		Scope: scope,
// 	}

// 	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

// 	if err != nil {
// 		return err
// 	}

// 	// Separating different sirets in different objects
// 	query2 := []bson.M{{
// 		"$unwind": bson.M{"path": "$value", "preserveNullAndEmptyArrays": false},
// 	},
// 		{
// 			"$project": bson.M{"_id": 0.0, "info": "$_id", "value": 1.0},
// 		},
// 		{"$out": collection}}
// 	pipe := Db.DB.C(collection).Pipe(query2)
// 	resp := []bson.M{}
// 	err = pipe.All(&resp)

// 	return err
// }

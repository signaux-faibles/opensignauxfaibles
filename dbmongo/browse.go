package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// prepareMRJob charge les fichiers MapReduce et fournit les paramètres pour l'exécution
func prepareMRJob(batchKey string, id string, typeJob string, target string, destination string) (*mgo.MapReduce, *bson.M, error) {
	query := &bson.M{
		"_id": id,
	}
	var output interface{}
	output = nil
	if id == "" {
		query = nil
		output = &bson.M{"merge": destination}
	}

	MR := MapReduceJS{}
	errEt := MR.load(typeJob, target)
	if errEt != nil {
		return &mgo.MapReduce{}, nil, fmt.Errorf("Problème d'accès aux fichiers MapReduce " + typeJob + "/" + target)
	}

	batch, _ := getBatch(batchKey)

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	naf, errNAF := loadNAF()
	if errNAF != nil {
		return &mgo.MapReduce{}, nil, fmt.Errorf("Problème d'accès aux fichiers naf")
	}

	job := &mgo.MapReduce{
		Map:      string(MR.Map),
		Reduce:   string(MR.Reduce),
		Finalize: string(MR.Finalize),
		Out:      output,
		Scope: bson.M{
			"date_debut":             dateDebut,
			"date_fin":               dateFin,
			"date_fin_effectif":      dateFinEffectif,
			"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
			"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
			"actual_batch":           batch.ID.Key,
			"naf":                    naf,
		},
	}

	return job, query, nil
}

// browserIndexHandler traite les requêtes d'indexation
func browseEtablissementHandler(c *gin.Context) {
	// var index struct {
	// 	BatchKey string `json:"batch" binding:"required"`
	// 	Siret    string `json:"siret"`
	// }

	// err := c.ShouldBind(&index)
	// if err != nil {
	// 	c.JSON(500, "Requête malformée: "+err.Error())
	// }

	siret := c.Params.ByName("siret")
	batch := c.Params.ByName("batch")

	result, err := browseEtablissement(batch, siret)

	if err != nil {
		c.JSON(500, "Erreur du traitement: "+err.Error())
		return
	}
	c.JSON(200, result)
}

func browseEtablissement(batchKey string, siret string) (interface{}, error) {
	var siren string
	if len(siret) == 14 {
		siren = siret[0:9]
	}

	// préparation des jobs
	var jobs [2]*mgo.MapReduce
	var queries [2]*bson.M
	var errMR [2]error
	jobs[0], queries[0], errMR[0] = prepareMRJob(batchKey, siret, "browser", "etablissement", "Browser")
	jobs[1], queries[1], errMR[1] = prepareMRJob(batchKey, siren, "browser", "entreprise", "Browser")

	if !allErrors(errMR[:], nil) {
		return nil, errors.New("Erreur dans la création du job MapReduce: ")
	}

	// exécution
	var resultEtablissement interface{}
	var resultEntreprise interface{}
	var err [2]error

	_, err[0] = db.DB.C("Etablissement").Find(queries[0]).MapReduce(jobs[0], &resultEtablissement)
	_, err[1] = db.DB.C("Entreprise").Find(queries[1]).MapReduce(jobs[1], &resultEntreprise)
	if !allErrors(err[:], nil) {
		errorMessage := ""
		if err[0] != nil {
			errorMessage = errorMessage + "\nEtablissement : " + err[0].Error() + " "
		}
		if err[1] != nil {
			errorMessage = errorMessage + "\nEntreprise: " + err[1].Error()
		}
		return nil, errors.New("Erreur dans l'exécution des jobs MapReduce" + errorMessage)
	}

	return map[string]interface{}{
		"etablissement": resultEtablissement,
		"entreprise":    resultEntreprise,
	}, nil
}

func publicEtablissementHandler(c *gin.Context) {
	batch := c.Params.ByName("batch")
	err := publicEtablissement(batch)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func publicEtablissement(batch string) error {
	// préparation des jobs
	job, _, err := prepareMRJob(batch, "", "public", "etablissement", "Public")

	if err != nil {
		return errors.New("Erreur dans la création du job MapReduce: " + err.Error())
	}

	// exécution

	_, err = db.DB.C("Etablissement").Find(nil).MapReduce(job, nil)

	if err != nil {
		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
	}

	return nil
}

func publicEntrepriseHandler(c *gin.Context) {
	batch := c.Params.ByName("batch")
	err := publicEntreprise(batch)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "ok")
}

func publicEntreprise(batch string) error {
	// préparation des jobs
	job, query, err := prepareMRJob(batch, "", "public", "entreprise", "Public")

	if err != nil {
		return errors.New("Erreur dans la création du job MapReduce: " + err.Error())
	}

	// exécution

	_, err = db.DB.C("Entreprise").Find(query).MapReduce(job, nil)

	if err != nil {
		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
	}

	return nil
}

func predictionBrowseHandler(c *gin.Context) {
	params := struct {
		Algo     string `json:"algo"`
		Batch    string `json:"batch"`
		Naf1     string `json:"naf1"`
		Effectif int    `json:"effectif"`
		Suivi    bool   `json:"suivi"`
		Ccsf     bool   `json:"ccsf"`
		Procol   bool   `json:"procol"`
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
	}{}

	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Bad Request: "+err.Error())
	}

	result, err := predictionBrowse(params.Batch, params.Naf1, params.Effectif, params.Suivi, params.Ccsf, params.Procol, params.Limit, params.Offset)
	if err != nil {
		c.JSON(500, err)
		return
	}
	fmt.Println(result)
	c.JSON(200, result)
}

func predictionBrowse(batch string, naf1 string, effectif int, suivi bool, ccsf bool, procol bool, limit int, offset int) (interface{}, error) {
	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id.batch": batch,
	}})

	if suivi {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"connu": false,
		}})
	}

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_idEntreprise.siren": bson.M{"$substrBytes": []interface{}{"$_id.siret", 0, 9}},
		"_idEntreprise.batch": "$_id.batch",
		"_id.siret":           "$_id.siret",
		"_id.batch":           "$_id.batch",
		"prob":                "$prob",
		"diff":                "$diff",
	}})

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"prob": -1,
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "_id",
		"foreignField": "_id",
		"as":           "etablissement"}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "_idEntreprise",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"etablissement": bson.M{"$arrayElemAt": []interface{}{"$etablissement", 0}},
		"entreprise":    bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}},
	}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"etablissement": "$etablissement.value",
		"entreprise":    "$entreprise.value",
	}})

	if naf1 != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"etablissement.sirene.ape": bson.M{"$in": naf5from1(naf1)},
		}})
	}

	if procol {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"etablissement.procol": "in_bonis",
		}})
	}

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"etablissement.effectif": bson.M{"$gt": effectif},
	}})

	pipeline = append(pipeline, bson.M{"$skip": offset})

	pipeline = append(pipeline, bson.M{"$limit": limit})
	var result = []interface{}{}
	err := db.DB.C("Prediction").Pipe(pipeline).All(&result)

	return result, err
}

func searchRaisonSociale(c *gin.Context) {
	var params struct {
		GuessRaisonSociale string `json:"guessRaisonSociale"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, err.Error())
		return
	}

	var result = make([]interface{}, 0)

	err = db.DBStatus.C("Public").Find(bson.M{"$text": bson.M{"$search": params.GuessRaisonSociale}}).Limit(15).All(&result)
	fmt.Println(err)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, result)
}

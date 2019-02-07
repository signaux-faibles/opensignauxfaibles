package engine

import (
	"dbmongo/lib/naf"

	"github.com/globalsign/mgo/bson"
)

// SearchCriteria critères de recherche d'une entreprise
type SearchCriteria struct {
	GuessRaisonSociale string `json:"guessRaisonSociale"`
}

// SearchRaisonSociale effectue une recherche texte sur la collection Public
func SearchRaisonSociale(params SearchCriteria) ([]interface{}, error) {
	var result = make([]interface{}, 0)
	err := Db.DBStatus.C("Public").Find(bson.M{"$text": bson.M{"$search": params.GuessRaisonSociale}}).Limit(15).All(&result)
	return result, err
}

// PredictionBrowse retourne la lise de prédiction filtrée pour la navigation
func PredictionBrowse(batch string, naf1 string, effectif int, suivi bool, ccsf bool, procol bool, limit int, offset int) (interface{}, error) {
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
			"etablissement.sirene.ape": bson.M{"$in": naf.Naf5from1(naf1)},
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
	err := Db.DB.C("Prediction").Pipe(pipeline).All(&result)

	return result, err
}

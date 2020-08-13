package engine

import (
	"github.com/globalsign/mgo/bson"
)

// SearchParams critères de recherche d'une entreprise
type SearchParams struct {
	Text string `json:"text"`
}

// Search effectue une recherche texte sur la collection Public
func Search(params SearchParams) ([]interface{}, error) {
	var result = make([]interface{}, 0)
	err := Db.DBStatus.C("Public").Find(
		bson.M{"$and": []interface{}{
			bson.M{"$or": []interface{}{
				bson.M{"$text": bson.M{"$search": params.Text}},
				//bson.M{"value.sirene.raison_sociale": bson.M{"$regex": params.Text}},
				bson.M{"_id.key": bson.M{"$regex": params.Text}},
			}},
			bson.M{"value.sirene": bson.M{"$exists": true}},
		},
		}).Limit(15).All(&result)
	return result, err
}

// EtablissementBrowseParams is type for params for prediction browser
type EtablissementBrowseParams struct {
	Siret string `json:"siret"`
	Batch string `json:"batch"`
}

// EtablissementBrowse retourne la lise de prédiction filtrée pour la navigation
func EtablissementBrowse(params EtablissementBrowseParams) (interface{}, error) {
	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id.scope": "etablissement",
		"_id.key":   params.Siret,
		"_id.batch": params.Batch,
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "value.idEntreprise",
		"foreignField": "value.idEntreprise",
		"as":           "etablissements"}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "value.idEntreprise",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"entreprise": bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}},
	}})

	var result = []interface{}{}
	err := Db.DB.C("Public").Pipe(pipeline).All(&result)

	return result, err
}

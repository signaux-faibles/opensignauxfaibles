package engine

import (
	"opensignauxfaibles/dbmongo/lib/naf"

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

// BrowseParams is type for params for prediction browser
type BrowseParams struct {
	Batch    string   `json:"batch"`
	Naf1     string   `json:"naf1"`
	Effectif int      `json:"effectif"`
	Suivi    *bool    `json:"suivi"`
	Ccsf     bool     `json:"ccsf"`
	Procol   bool     `json:"procol"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Zone     []string `json:"zone"`
}

// PredictionBrowse retourne la lise de prédiction filtrée pour la navigation
func PredictionBrowse(params BrowseParams) (interface{}, error) {
	var pipeline []bson.M

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id.batch": params.Batch,
	}})

	if params.Suivi != nil {
		if *params.Suivi {
			pipeline = append(pipeline, bson.M{"$match": bson.M{
				"connu": true,
			}})
		} else {
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"$and": []interface{}{
						bson.M{"connu": false},
						// bson.M{"procol": "in_bonis"},
					},
				},
			})
		}
	}

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{"inZone": bson.M{"$in": []interface{}{"$departement", params.Zone}}}})

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"inZone": true,
	}})

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id.scope":   "etablissement",
		"_id.key":     "$_id.siret",
		"_id.batch":   "$_id.batch",
		"prob":        "$prob",
		"diff":        "$diff",
		"test":        params.Zone,
		"inZone":      bson.M{"inZone": bson.M{"$in": []interface{}{"$departement", params.Zone}}},
		"departement": "$departement",
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
		"localField":   "etablissement.0.value.idEntreprise",
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

	if params.Naf1 != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"etablissement.sirene.ape": bson.M{"$in": naf.Naf5from1(params.Naf1)},
		}})
	}

	if params.Procol {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"etablissement.procol": "in_bonis",
		}})
	}

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"etablissement.dernier_effectif.effectif": bson.M{"$gt": params.Effectif},
	}})

	pipeline = append(pipeline, bson.M{"$skip": params.Offset})

	pipeline = append(pipeline, bson.M{"$limit": params.Limit})
	var result = []interface{}{}
	err := Db.DB.C("Prediction").Pipe(pipeline).All(&result)

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

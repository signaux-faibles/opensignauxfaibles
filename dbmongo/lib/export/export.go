package export

import (
	"github.com/globalsign/mgo/bson"
)

// GetEntreprisePipeline produit un pipeline pour exporter les établissements avec leur scores.
func GetEntreprisePipeline(key string) (pipeline []bson.M) {
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id": bson.RegEx{
			Pattern: "entreprise_" + key,
		},
	}})

	return pipeline
}

// GetEtablissementWithScoresPipeline produit un pipeline pour exporter les établissements avec leur scores.
func GetEtablissementWithScoresPipeline(key string) (pipeline []bson.M) {
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id": bson.RegEx{
			Pattern: "etablissement_" + key + ".*",
		},
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Scores",
		"localField":   "value.key",
		"foreignField": "siret",
		"as":           "scores"}})

	return pipeline
}

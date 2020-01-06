package exportdatapi

import (
	"github.com/globalsign/mgo/bson"
	daclient "github.com/signaux-faibles/datapi/client"
)

// GetEtablissementPipeline produit un pipeline pour exporter les établissements vers datapi
func GetEtablissementPipeline(batch string, key string) (pipeline []bson.M) {
	if key == "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"_id": bson.RegEx{
				Pattern: batch + "_etablissement_.*",
			},
		}})
	} else {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"_id": batch + "_etablissement_" + key,
		}})
	}

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Scores",
		"localField":   "value.key",
		"foreignField": "siret",
		"as":           "detection"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"detectionLength": bson.M{
			"$size": "$detection",
		},
	}})

	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"detectionLength": bson.M{
			"$gt": 0,
		},
	}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"idEntreprise": bson.M{
			"$concat": []interface{}{
				batch + "_etablissement_",
				bson.M{"$substr": []interface{}{"$value.key", 0, 9}},
			},
		},
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "idEntreprise",
		"foreignField": "_id",
		"as":           "entreprise"}})

	pipeline = append(pipeline, bson.M{"$addFields": bson.M{
		"entreprise": bson.M{"$arrayElemAt": []interface{}{"$entreprise", 0}},
	}})

	return pipeline
}

// ComputeEtablissement transforme un établissiment au format public en objet datapi
func ComputeEtablissement(data Etablissement, connus *[]string) []daclient.Object {
	var objects []daclient.Object

	urssaf := UrssafScope(data.Value.Compte.Numero, data.Value.Sirene.Departement)

	key := map[string]string{
		"siret": data.Value.Key,
		"siren": data.Value.Key[0:9],
		urssaf:  "true",
		"type":  "etablissement",
	}

	sirene := data.Value.Sirene
	if data.Entreprise != nil {
		sirene.RaisonSociale = data.Entreprise.Value.SireneUL.RaisonSociale
	}
	value := make(map[string]interface{})

	if data.Entreprise != nil {
		value["sirene"] = sirene
		value["connu"] = findString(data.Value.Key, *connus)
		value["detection"] = data.Detection

		if len(data.Entreprise.Value.Diane) > 0 {
			value["diane"] = data.Entreprise.Value.Diane
		}

		if len(data.Value.Effectif) > 0 {
			value["effectif"] = data.Value.Effectif
		}

		if len(data.Value.Procol) > 0 {
			value["procedure_collective"] = data.Value.Procol
		}
	}

	objectPublic := daclient.Object{
		Key:   key,
		Scope: []string{},
		Value: value,
	}

	scope := []string{data.Value.Sirene.Departement}
	objectPrivate := daclient.Object{
		Key:   key,
		Scope: scope,
		Value: map[string]interface{}{
			"detection": data.Detection,
		},
	}
	objects = append(objects, objectPublic, objectPrivate)

	if eligible(data) {
		send := false
		scopeURSSAF := []string{"urssaf", data.Value.Sirene.Departement}
		valueURSSAF := make(map[string]interface{})
		if len(data.Value.Debit) > 0 {
			valueURSSAF["debit"] = data.Value.Debit
			send = true
		}
		if len(data.Value.Delai) > 0 {
			valueURSSAF["delai"] = data.Value.Delai
			send = true
		}
		if len(data.Value.Cotisation) > 0 {
			valueURSSAF["cotisation"] = data.Value.Cotisation
			send = true
		}

		if send {
			objectEtablissementURSSAF := daclient.Object{
				Key:   key,
				Scope: scopeURSSAF,
				Value: valueURSSAF,
			}
			objects = append(objects, objectEtablissementURSSAF)
		}

		send = false
		scopeDGEFP := []string{"dgefp", data.Value.Sirene.Departement}
		valueDGEFP := make(map[string]interface{})
		if len(data.Value.APConso) > 0 {
			valueDGEFP["apconso"] = data.Value.APConso
			send = true
		}
		if len(data.Value.APDemande) > 0 {
			valueDGEFP["apdemande"] = data.Value.APDemande
			send = true
		}
		if send {
			objectEtablissementDGEFP := daclient.Object{
				Key:   key,
				Scope: scopeDGEFP,
				Value: valueDGEFP,
			}
			objects = append(objects, objectEtablissementDGEFP)
		}

		if data.Entreprise != nil && len(data.Entreprise.Value.BDF) > 0 {
			scopeBDF := []string{"bdf", data.Value.Sirene.Departement}
			valueBDF := map[string]interface{}{
				"bdf": data.Entreprise.Value.BDF,
			}
			objectEtablissementBDF := daclient.Object{
				Key:   key,
				Scope: scopeBDF,
				Value: valueBDF,
			}
			objects = append(objects, objectEtablissementBDF)
		}

	}

	return objects
}

func eligible(data Etablissement) bool {
	for _, i := range data.Detection {
		if i.Alert != "Pas d'alerte" {
			return true
		}
	}
	return false
}

type datapiValue map[string]interface{}

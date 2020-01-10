package exportdatapi

import (
	"errors"
	"strconv"
	"time"

	daclient "github.com/signaux-faibles/datapi/client"

	"github.com/globalsign/mgo/bson"
)

// GetDetectionPipeline construit le pipeline d'aggregation pour exporter une détection
func GetDetectionPipeline(batch, key string, algo string) (pipeline []bson.M) {
	if key == "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"algo":  algo,
			"batch": batch,
		}})
	} else {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"algo":  algo,
			"batch": batch,
			"siret": key,
		}})
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{
		"siret":     1,
		"periode":   1,
		"timestamp": -1,
	}})

	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"siret":   "$siret",
				"periode": "$periode",
				"batch":   "$batch",
				"algo":    "$algo",
			},
			"score": bson.M{
				"$first": "$score",
			},
			"alert": bson.M{
				"$first": "$alert",
			},
			"diff": bson.M{
				"$first": "$diff",
			},
			"timestamp": bson.M{
				"$first": "$timestamp",
			},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"alert": bson.M{"$ne": "Pas d'alerte"},
		},
	})

	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"siret":   "$_id.siret",
			"periode": "$_id.periode",
			"batch":   "$_id.batch",
			"algo":    "$_id.algo",
		},
	})

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_idEtablissement": bson.M{
			"$concat": []interface{}{
				"etablissement_",
				"$siret",
			},
		},

		"_idEntreprise": bson.M{
			"$concat": []interface{}{
				"entreprise_",
				bson.M{"$substr": []interface{}{"$siret", 0, 9}},
			},
		},
		"score": "$score",
		"alert": "$alert",
		"diff":  "$score_diff",
		"connu": "$connu",
		"algo":  "$algo",
	}})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "Public",
		"localField":   "_idEtablissement",
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

	return pipeline
}

func computeDetection(detection Detection, connu *[]string) (detections []daclient.Object) {
	caVal, caVar, reVal, reVar, annee := computeDiane(detection)
	dernierEffectif, variationEffectif := computeEffectif(detection)

	urssaf := UrssafScope(detection.Etablissement.Value.Compte.Numero, detection.Etablissement.Value.Sirene.Departement)

	key := map[string]string{
		"siret": detection.ID["siret"],
		"siren": detection.ID["siret"][0:9],
		"batch": detection.ID["batch"] + "." + detection.Algo,
		"type":  "detection",
		urssaf:  "true",
	}

	scopeB := []string{"detection", detection.Etablissement.Value.Sirene.Departement}
	valueB := map[string]interface{}{
		"raison_sociale":          detection.Entreprise.Value.SireneUL.RaisonSociale,
		"activite":                detection.Etablissement.Value.Sirene.Ape,
		"urssaf":                  computeUrssaf(detection),
		"activite_partielle":      computeActivitePartielle(detection),
		"dernier_effectif":        &dernierEffectif,
		"variation_effectif":      &variationEffectif,
		"annee_ca":                annee,
		"ca":                      caVal,
		"variation_ca":            caVar,
		"resultat_expl":           reVal,
		"variation_resultat_expl": reVar,
		"departement":             detection.Etablissement.Value.Sirene.Departement,
		"etat_procol":             detection.Etablissement.Value.LastProcol.Etat,
		"date_procol":             detection.Etablissement.Value.LastProcol.Date,
	}

	scopeA := []string{"detection", "score", detection.Etablissement.Value.Sirene.Departement}
	valueA := map[string]interface{}{
		"score": detection.Score,
		"alert": detection.Alert,
		"diff":  detection.Diff,
	}

	detections = append(detections, daclient.Object{
		Key:   key,
		Scope: scopeA,
		Value: valueA,
	})

	detections = append(detections, daclient.Object{
		Key:   key,
		Scope: scopeB,
		Value: valueB,
	})

	return detections
}

func computeEtablissement(detection Detection, connus *[]string) (objects []daclient.Object) {
	urssaf := UrssafScope(detection.Etablissement.Value.Compte.Numero, detection.Etablissement.Value.Sirene.Departement)

	key := map[string]string{
		"siret": detection.ID["siret"],
		"siren": detection.ID["siret"][0:9],
		"batch": detection.ID["batch"] + "." + detection.Algo,
		urssaf:  "true",
		"type":  "detail",
	}

	scope := []string{detection.Etablissement.Value.Sirene.Departement}
	sirene := detection.Etablissement.Value.Sirene
	sirene.RaisonSociale = detection.Entreprise.Value.SireneUL.RaisonSociale
	value := map[string]interface{}{
		"diane":                detection.Entreprise.Value.Diane,
		"effectif":             detection.Etablissement.Value.Effectif,
		"sirene":               sirene,
		"procedure_collective": detection.Etablissement.Value.Procol,
		"connu":                findString(detection.ID["siret"], *connus),
	}

	scopeURSSAF := []string{"urssaf", detection.Etablissement.Value.Sirene.Departement}
	valueURSSAF := map[string]interface{}{
		"debit":      detection.Etablissement.Value.Debit,
		"delai":      detection.Etablissement.Value.Delai,
		"cotisation": detection.Etablissement.Value.Cotisation,
	}

	scopeDGEFP := []string{"dgefp", detection.Etablissement.Value.Sirene.Departement}
	valueDGEFP := map[string]interface{}{
		"apconso":   detection.Etablissement.Value.APConso,
		"apdemande": detection.Etablissement.Value.APDemande,
	}

	scopeBDF := []string{"bdf", detection.Etablissement.Value.Sirene.Departement}
	valueBDF := map[string]interface{}{
		"bdf": detection.Entreprise.Value.BDF,
	}

	object := daclient.Object{
		Key:   key,
		Scope: scope,
		Value: value,
	}

	objectURSSAF := daclient.Object{
		Key:   key,
		Scope: scopeURSSAF,
		Value: valueURSSAF,
	}

	if detection.Alert != "Pas d'alerte" {
		objectDGEFP := daclient.Object{
			Key:   key,
			Scope: scopeDGEFP,
			Value: valueDGEFP,
		}
		objects = append(objects, objectDGEFP)
	}

	objectBDF := daclient.Object{
		Key:   key,
		Scope: scopeBDF,
		Value: valueBDF,
	}

	objects = append(objects, object, objectURSSAF, objectBDF)
	return objects
}

// ComputeDetection traite un objet detection pour produire les objets datapi
func ComputeDetection(detection Detection, connus *[]string) ([]daclient.Object, error) {
	if detection.Etablissement.Value.Sirene.Departement != "" {
		var objects []daclient.Object
		objects = append(objects, computeDetection(detection, connus)...)
		return objects, nil
	}

	return nil, errors.New(detection.ID["siret"] + ": " + detection.Alert + " pas d'information sirene, objet ignoré")
}

func computeEffectif(detection Detection) (dernierEffectif int, variationEffectif float64) {
	l := len(detection.Etablissement.Value.Effectif)
	if l > 2 {
		dernierEffectif := detection.Etablissement.Value.Effectif[l-1].Effectif
		variationEffectif := float64(detection.Etablissement.Value.Effectif[l-1].Effectif) / float64(detection.Etablissement.Value.Effectif[l-2].Effectif)
		return dernierEffectif, variationEffectif
	}
	return 0, 0
}

func computeUrssaf(detection Detection) bool {
	debits := detection.Etablissement.Value.Debit
	if len(debits) == 24 {
		for i := 24 - 3; i < 24; i++ {
			if (debits[i].PartOuvriere+debits[i].PartPatronale)/(debits[i-1].PartOuvriere+debits[i-1].PartPatronale) > 1.01 {
				return true
			}
		}
	}
	return false
}

func computeActivitePartielle(detection Detection) bool {
	batch := detection.ID["batch"]
	date, err := batchToTime(batch)
	if err != nil {
		return false
	}
	for _, v := range detection.Etablissement.Value.APConso {
		if v.Periode.Add(1 * time.Second).After(date) {
			return true
		}
	}
	for _, v := range detection.Etablissement.Value.APDemande {
		if v.Periode.End.Add(1 * time.Second).After(date) {
			return true
		}
	}

	return false
}

func computeDiane(detection Detection) (caVal *float64, caVar *float64, reVal *float64, reVar *float64, annee *float64) {
	for i := 1; i < len(detection.Entreprise.Value.Diane); i++ {
		if detection.Entreprise.Value.Diane[i-1].ChiffreAffaire != 0 &&
			detection.Entreprise.Value.Diane[i].ChiffreAffaire != 0 {
			d1 := detection.Entreprise.Value.Diane[i-1]
			d2 := detection.Entreprise.Value.Diane[i]
			annee = &detection.Entreprise.Value.Diane[i-1].Exercice
			cavar := d1.ChiffreAffaire / d2.ChiffreAffaire
			caVal = &d1.ChiffreAffaire
			caVar = &cavar

			if d2.ResultatExploitation*d1.ResultatExploitation != 0 {
				reVal = &d1.ResultatExploitation
				revar := d1.ResultatExploitation / d2.ResultatExploitation
				reVar = &revar
			}

			break
		}
	}

	return caVal, caVar, reVal, reVar, annee
}

func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
}

func reverseMap(input map[string]string) map[string][]string {
	var output = make(map[string][]string)
	for k, v := range input {
		output[v] = append(output[v], k)
	}
	return output
}

func findString(s string, a []string) bool {
	for _, v := range a {
		if len(v) > 9 && len(s) > 9 && s[0:9] == v[0:9] {
			return true
		}
	}
	return false
}

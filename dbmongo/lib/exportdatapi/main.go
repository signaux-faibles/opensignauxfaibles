package exportdatapi

import (
	"errors"
	"strconv"
	"time"

	daclient "github.com/signaux-faibles/datapi/client"

	"github.com/globalsign/mgo/bson"
)

// GetPipeline construit le pipeline d'aggregation
func GetPipeline(batch string) (pipeline []bson.M) {
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"batch": batch,
	}})

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id": bson.D{
			{Name: "scope", Value: "etablissement"},
			{Name: "key", Value: "$siret"},
			{Name: "batch", Value: "$batch"},
		},
		"_idEntreprise": bson.D{
			{Name: "scope", Value: "entreprise"},
			{Name: "key", Value: bson.M{"$substr": []interface{}{"$siret", 0, 9}}},
			{Name: "batch", Value: "$batch"},
		},
		"prob":  "$score",
		"diff":  "$diff",
		"connu": "$connu",
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

	return pipeline
}

// Detection correspond aux données retournées pour l'export Datapi
type Detection struct {
	ID            map[string]string `json:"_id" bson:"_id"`
	Prob          float64           `json:"prob" bson:"prob"`
	Diff          float64           `json:"diff" bson:"diff"`
	Connu         bool              `json:"connu" bson:"connu"`
	Etablissement Etablissement     `json:"etablissement" bson:"etablissement"`
	Entreprise    Entreprise        `json:"entreprise" bson:"entreprise"`
}

// Etablissement is an object
type Etablissement struct {
	ID    map[string]string `bson:"_id"`
	Value struct {
		Sirene          Sirene        `json:"sirene" bson:"sirene"`
		Cotisation      []float64     `json:"cotisation" bson:"cotisation"`
		Debit           []Debit       `json:"debit" bson:"debit"`
		APDemande       []APDemande   `json:"apdemande" bson:"apdemande"`
		APConso         []APConso     `json:"apconso" bson:"apconso"`
		Effectif        []Effectif    `json:"effectif" bson:"effectif"`
		DernierEffectif Effectif      `json:"dernier_effectif" bson:"dernier_effectif"`
		Delai           []interface{} `json:"delai" bson:"delai"`
	} `bson:"value"`
}

// Entreprise object
type Entreprise struct {
	ID    map[string]string `json:"_id" bson:"_id"`
	Value struct {
		Diane []Diane       `json:"diane" bson:"diane"`
		BDF   []interface{} `json:"bdf" bson:"bdf"`
	} `bson:"value"`
}

// Effectif detail
type Effectif struct {
	Periode  time.Time `json:"periode" bson:"periode"`
	Effectif int       `json:"effectif" bson:"effectif"`
}

// Debit detail
type Debit struct {
	PartOuvriere  float64   `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale float64   `json:"part_patronale" bson:"part_patronale"`
	Periode       time.Time `json:"periode" bson:"periode"`
}

// APConso detail
type APConso struct {
	IDConso       string    `json:"id_conso" bson:"id_conso"`
	HeureConsomme float64   `json:"heure_consomme" bson:"heure_consomme"`
	Montant       float64   `json:"montant" bson:"montant"`
	Effectif      int       `json:"int" bson:"int"`
	Periode       time.Time `json:"periode" bson:"periode"`
}

// APDemande detail
type APDemande struct {
	DateStatut time.Time `json:"date_statut" bson:"date_statut"`
	Periode    struct {
		Start time.Time `json:"start" bson:"start"`
		End   time.Time `json:"end" bson:"end"`
	} `json:"periode" bson:"periode"`
	EffectifAutorise int     `json:"effectif_autorise" bson:"effectif_autorise"`
	EffectifConsomme int     `json:"effectif_consomme" bson:"effectif_consomme"`
	IDDemande        string  `json:"id_conso" bson:"id_conso"`
	Effectif         int     `json:"int" bson:"int"`
	MTA              float64 `json:"mta" bson:"mta"`
	HTA              float64 `json:"hta" bson:"hta"`
	MotifRecoursSE   int     `json:"motif_recours_se" bson:"motif_recours_se"`
	HeureConsomme    float64 `json:"heure_consomme" bson:"heure_consomme"`
	Montant          float64 `json:"montant" bson:"montant"`
}

// Diane detail
type Diane struct {
	ChiffreAffaire       float64 `json:"ca" bson:"ca"`
	ResultatExploitation float64 `json:"benefice_ou_perte" bson:"benefice_ou_perte"`
	Exercice             float64 `json:"exercice_diane" bson:"exercice_diane"`
}

// Sirene detail
type Sirene struct {
	Region          string   `json:"region" bson:"region"`
	Commune         string   `json:"commune" bson:"commune"`
	RaisonSociale   string   `json:"raison_sociale" bson:"raison_sociale"`
	TypeVoie        string   `json:"type_voie" bson:"type_voie"`
	Siren           string   `json:"siren" bson:"siren"`
	CodePostal      string   `json:"code_postal" bson:"code_postal"`
	Lattitude       float64  `json:"lattitude" bson:"lattitude"`
	Adresse         []string `json:"adresse" bson:"adresse"`
	Departement     string   `json:"departement" bson:"departement"`
	NatureJuridique string   `json:"nature_juridique" bson:"nature_juridique"`
	NumeroVoie      string   `json:"numero_voie" bson:"numero_voie"`
	Ape             string   `json:"ape" bson:"ape"`
	Longitude       float64  `json:"longitude" bson:"longitude"`
	Nic             string   `json:"nic" bson:"nic"`
	NicSiege        string   `json:"nic_siege" bson:"nic_siege"`
}

// DatapiDetection résultat de l'aggregation
type DatapiDetection struct {
	Key   map[string]string
	Scope []string
	Value map[string]interface{}
}

func computeDetection(detection Detection) (detections []daclient.Object) {
	caVal, caVar, reVal, reVar := computeDiane(detection)
	dernierEffectif, variationEffectif := computeEffectif(detection)

	key := map[string]string{
		"siret": detection.ID["key"],
		"batch": detection.ID["batch"],
		"type":  "detection",
	}

	var acteurs []string
	if detection.Connu {
		acteurs = append(acteurs, "connu")
	}

	scopeB := []string{"detection", detection.Etablissement.Value.Sirene.Departement}
	valueB := map[string]interface{}{
		"acteurs":                 acteurs,
		"raison_sociale":          detection.Etablissement.Value.Sirene.RaisonSociale,
		"activite":                detection.Etablissement.Value.Sirene.Ape,
		"urssaf":                  computeUrssaf(detection),
		"activite_partielle":      computeActivitePartielle(detection),
		"dernier_effectif":        &dernierEffectif,
		"variation_effectif":      &variationEffectif,
		"ca":                      caVal,
		"variation_ca":            caVar,
		"resultat_expl":           reVal,
		"variation_resultat_expl": reVar,
		// "procedure_collective":    detection.Etablissement.Value.Procol,
	}

	scopeA := []string{"detection", "score", detection.Etablissement.Value.Sirene.Departement}
	valueA := map[string]interface{}{
		"prob": detection.Prob,
		"diff": detection.Diff,
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

func computeEtablissement(detection Detection) (objects []daclient.Object) {
	key := map[string]string{
		"siret": detection.Etablissement.Value.Sirene.Siren + detection.Etablissement.Value.Sirene.Nic,
		"siren": detection.Etablissement.Value.Sirene.Siren,
		"batch": detection.ID["batch"],
		"type":  "detail",
	}

	scope := []string{detection.Etablissement.Value.Sirene.Departement}
	value := map[string]interface{}{
		"diane":    detection.Entreprise.Value.Diane,
		"effectif": detection.Etablissement.Value.Effectif,
		"sirene":   detection.Etablissement.Value.Sirene,
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

	objectDGEFP := daclient.Object{
		Key:   key,
		Scope: scopeDGEFP,
		Value: valueDGEFP,
	}

	objectBDF := daclient.Object{
		Key:   key,
		Scope: scopeBDF,
		Value: valueBDF,
	}

	objects = append(objects, object, objectURSSAF, objectDGEFP, objectBDF)
	return objects
}

// Compute traite un objet detection pour produire les objets datapi
func Compute(detection Detection) ([]daclient.Object, error) {

	if detection.Etablissement.Value.Sirene.Departement != "" {
		var objects []daclient.Object
		objects = append(objects, computeDetection(detection)...)
		objects = append(objects, computeEtablissement(detection)...)
		return objects, nil
	}

	return nil, errors.New("pas d'information sirene, objet ignoré")
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

func computeDiane(detection Detection) (caVal *float64, caVar *float64, reVal *float64, reVar *float64) {
	if len(detection.Entreprise.Value.Diane) > 1 {
		d1 := detection.Entreprise.Value.Diane[0]
		d2 := detection.Entreprise.Value.Diane[1]

		if d2.ChiffreAffaire*d1.ChiffreAffaire != 0 {
			cavar := d1.ChiffreAffaire / d2.ChiffreAffaire
			caVal = &d1.ChiffreAffaire
			caVar = &cavar
		}

		if d2.ResultatExploitation*d1.ResultatExploitation != 0 {
			reVal = &d1.ResultatExploitation
			revar := d1.ResultatExploitation / d2.ResultatExploitation
			reVar = &revar
		}

	}

	return caVal, caVar, reVal, reVar
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

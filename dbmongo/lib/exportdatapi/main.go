package exportdatapi

import (
	"errors"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
)

// GetPipeline construit le pipeline d'aggregation
func GetPipeline(batch string) (pipeline []bson.M) {
	pipeline = append(pipeline, bson.M{"$match": bson.M{
		"_id.batch": batch,
	}})

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id": bson.D{
			{Name: "scope", Value: "etablissement"},
			{Name: "key", Value: "$_id.siret"},
			{Name: "batch", Value: "$_id.batch"},
		},
		"_idEntreprise": bson.D{
			{Name: "scope", Value: "entreprise"},
			{Name: "key", Value: bson.M{"$substr": []interface{}{"$_id.siret", 0, 9}}},
			{Name: "batch", Value: "$_id.batch"},
		},
		"prob": "$prob",
		"diff": "$diff",
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
	ID            map[string]string `bson:"_id"`
	Prob          float64           `bson:"prob"`
	Diff          float64           `bson:"diff"`
	Etablissement Etablissement     `bson:"etablissement"`
	Entreprise    Entreprise        `bson:"entreprise"`
}

// Etablissement is an object
type Etablissement struct {
	ID    map[string]string `bson:"_id"`
	Value struct {
		Sirene          Sirene        `bson:"sirene"`
		Cotisation      []float64     `bson:"cotisation"`
		Debit           []Debit       `bson:"debit"`
		APDemande       []APDemande   `bson:"apdemande"`
		APConso         []APConso     `bson:"apconso"`
		Effectif        []Effectif    `bson:"effectif"`
		DernierEffectif Effectif      `bson:"dernier_effectif"`
		Delai           []interface{} `bson:"delai"`
	} `bson:"value"`
}

// Entreprise object
type Entreprise struct {
	ID    map[string]string `bson:"_id"`
	Value struct {
		Diane []Diane       `bson:"diane"`
		BDF   []interface{} `bson:"bdf"`
	} `bson:"value"`
}

// Effectif detail
type Effectif struct {
	Periode  time.Time `bson:"periode"`
	Effectif int       `bson:"effectif"`
}

// Debit detail
type Debit struct {
	PartOuvriere  float64   `bson:"part_ouvriere"`
	PartPatronale float64   `bson:"part_patronale"`
	Periode       time.Time `bson:"periode"`
}

// APConso detail
type APConso struct {
	IDConso       string    `bson:"id_conso"`
	HeureConsomme float64   `bson:"heure_consomme"`
	Montant       float64   `bson:"montant"`
	Effectif      int       `bson:"int"`
	Periode       time.Time `bson:"periode"`
}

// APDemande detail
type APDemande struct {
	DateStatut time.Time `bson:"date_statut"`
	Periode    struct {
		Start time.Time `bson:"start"`
		End   time.Time `bson:"end"`
	} `bson:"periode"`
	EffectifAutorise int     `bson:"effectif_autorise"`
	EffectifConsomme int     `bson:"effectif_consomme"`
	IDDemande        string  `bson:"id_conso"`
	Effectif         int     `bson:"int"`
	MTA              float64 `bson:"mta"`
	HTA              float64 `bson:"hta"`
	MotifRecoursSE   int     `bson:"motif_recours_se"`
	HeureConsomme    float64 `bson:"heure_consomme"`
	Montant          float64 `bson:"montant"`
}

// Diane detail
type Diane struct {
	ChiffreAffaire       float64 `bson:"ca"`
	ResultatExploitation float64 `bson:"benefice_ou_perte"`
	Exercice             float64 `bson:"exercice_diane"`
}

// Sirene detail
type Sirene struct {
	Region          string   `bson:"region"`
	Commune         string   `bson:"commune"`
	RaisonSociale   string   `bson:"raison_sociale"`
	TypeVoie        string   `bson:"type_voie"`
	Siren           string   `bson:"siren"`
	CodePostal      string   `bson:"code_postal"`
	Lattitude       float64  `bson:"lattitude"`
	Adresse         []string `bson:"adresse"`
	Departement     string   `bson:"departement"`
	NatureJuridique string   `bson:"nature_juridique"`
	NumeroVoie      string   `bson:"numero_voie"`
	Ape             string   `bson:"ape"`
	Longitude       float64  `bson:"longitude"`
	Nic             string   `bson:"nic"`
	NicSiege        string   `bson:"nic_siege"`
}

// DatapiDetection résultat de l'aggregation
type DatapiDetection struct {
	Key   map[string]string
	Scope []string
	Value map[string]interface{}
}

// ComputeDetection traite un objet detection pour sortir les éléments nécessaires à une liste de détection
func ComputeDetection(detection Detection) (DatapiDetection, error) {
	caVal, caVar, reVal, reVar := computeDiane(detection)

	if detection.Etablissement.Value.Sirene.Departement != "" {
		scope := []string{"detection", detection.Etablissement.Value.Sirene.Departement}
		d := DatapiDetection{
			Key: map[string]string{
				"siret": detection.ID["key"],
				"batch": detection.ID["batch"],
				"type":  "detection",
			},
			Scope: scope,
			Value: map[string]interface{}{
				"prob":                    detection.Prob,
				"diff":                    detection.Diff,
				"raison_sociale":          detection.Etablissement.Value.Sirene.RaisonSociale,
				"activite":                detection.Etablissement.Value.Sirene.Ape,
				"urssaf":                  computeUrssaf(detection),
				"activite_partielle":      computeActivitePartielle(detection),
				"ca":                      caVal,
				"variation_ca":            caVar,
				"resultat_expl":           reVal,
				"variation_resultat_expl": reVar,
			},
		}
		return d, nil
	}
	return DatapiDetection{}, errors.New("pas d'information sirene")
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

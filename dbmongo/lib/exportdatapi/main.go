package exportdatapi

import (
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

	return pipeline
}

// Detection correspond aux données retournées pour l'export Datapi
type Detection struct {
	ID            map[string]string     `bson:"_id"`
	Prob          float64               `bson:"prob"`
	Diff          float64               `bson:"diff"`
	Etablissement []EtablissementObject `bson:"etablissement"`
	Entreprise    []interface{}         `bson:"entreprise"`
}

// EtablissementObject is an object
type EtablissementObject struct {
	ID    map[string]string `bson:"_id"`
	Value Etablissement     `bson:"etablissement"`
}

// Etablissement object
type Etablissement struct {
	Sirene          Sirene        `bson:"sirene"`
	Cotisation      []float64     `bson:"cotisation"`
	Debit           []Debit       `bson:"debit"`
	APDemande       []interface{} `bson:"apdemande"`
	APConso         []interface{} `bson:"apconso"`
	Effectif        []Effectif    `bson:"effectif"`
	DernierEffectif Effectif      `bson:"dernier_effectif"`
	Delai           []interface{} `bson:"delai"`
}

// EntrepriseObject object
type EntrepriseObject struct {
	ID    map[string]string `bson:"_id"`
	Value Entreprise        `bson:"entreprise"`
}

// Entreprise object
type Entreprise struct {
	Diane map[string]interface{} `bson:"ca"`
	BDF   map[string]interface{} `bson:"bdf"`
}

// Effectif detail
type Effectif struct {
	Periode  time.Time `bson:"periode"`
	Effectif float64   `bson:"effectif"`
}

// Debit detail
type Debit struct {
	PartOuvriere  float64   `json:"part_ouvriere"`
	PartPatronale float64   `json:"part_patronale"`
	Periode       time.Time `json:"periode"`
}

// Sirene detail
type Sirene struct {
	Region          string
	Commune         string
	RaisonSociale   string
	TypeVoie        string
	Siren           string
	CodePostal      string
	Lattitude       string
	Adresse         []string
	Departement     string
	NatureJuridique string
	NumeroVoie      string
	Ape             string
	Longitude       string
	Nic             string
	NicSiege        string
}

// DatapiDetection
type DatapiDetection struct {
	Key   map[string]string
	Scope []string
	Value map[string]interface{}
}

// ComputeDetection traite un objet detection pour sortir les éléments nécessaires à une liste de détection
func ComputeDetection(detection Detection) interface{} {
	return detection
}

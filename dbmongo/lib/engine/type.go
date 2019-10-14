package engine

import (
	"errors"
)

// Cache saves values in memory
type Cache map[string]interface{}

// Get gets a value from the cache
func (ca Cache) Get(name string) (interface{}, error) {
	if ca == nil {
		return nil, errors.New("Entry not found: " + name)
	}
	if _, ok := ca[name]; !ok {
		return nil, errors.New("Entry not found: " + name)
	}
	return ca[name], nil
}

// Set writes a value to the Cache
func (ca Cache) Set(name string, value interface{}) {
	ca[name] = value
}

// NewCache returns a new cache object
func NewCache() Cache {
	return map[string]interface{}{}
}

// GetTypes retourne la liste des types déclarés
func GetTypes() Types {
	return []Type{
		{"admin_urssaf", "Siret/Compte URSSAF", "Liste comptes"},
		{"apconso", "Consommation Activité Partielle", "conso"},
		{"bdf", "Ratios Banque de France", "bdf"},
		{"cotisation", "Cotisations URSSAF", "cotisation"},
		{"delai", "Délais URSSAF", "delais|Délais"},
		{"dpae", "Déclaration Préalable à l'embauche", "DPAE"},
		{"interim", "Base Interim", "interim"},
		{"altares", "Base Altarès", "ALTARES"},
		{"procol", "Procédures collectives", "procol"},
		{"apdemande", "Demande Activité Partielle", "dde"},
		{"ccsf", "Stock CCSF à date", "ccsf"},
		{"debit", "Débits URSSAF", "debit"},
		{"dmmo", "Déclaration Mouvement de Main d'Œuvre", "dmmo"},
		{"effectif", "Emplois URSSAF", "Emploi"},
		{"sirene", "Base GéoSirene", "sirene"},
		{"diane", "Diane", "diane"},
	}
}

// Type description des types de fichiers pris en charge
type Type struct {
	Type    string `json:"type" bson:"type"`
	Libelle string `json:"text" bson:"text"`
	Filter  string `json:"filter" bson:"filter"`
}

// Types is a Type array
type Types []Type

// ToData transforms Types to datapi compatible type
func (t Types) ToData() map[string]interface{} {
	r := make(map[string]interface{})
	for _, v := range t {
		r[v.Type] = map[string]string{
			"text":   v.Libelle,
			"filter": v.Filter,
		}
	}
	return r
}

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *AdminBatch) (chan Tuple, chan Event)

// Browseable est le type qui permet d'envoyer les objets vers le frontend
// Voir la fonction Browse
type Browseable struct {
	ID struct {
		Key   string   `json:"key" bson:"key"`
		Scope []string `json:"scope" bson:"scope"`
	}
	Value map[string]interface{} `json:"value" bson:"value"`
}

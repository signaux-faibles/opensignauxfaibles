package sirene

import (
	//"bufio"

	"opensignauxfaibles/lib/engine"
	"time"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren                string     `input:"siren"                                       json:"siren,omitempty"                 sql:"siren"                    csv:"siren"`
	Nic                  string     `input:"nic"                                         json:"nic,omitempty"                                                  csv:"nic"`
	Siret                string     `                                                    json:"-"                               sql:"siret"`
	Siege                bool       `input:"etablissementSiege"                          json:"siege,omitempty"                 sql:"siege"                    csv:"siege"`
	ComplementAdresse    string     `input:"complementAdresseEtablissement"              json:"complement_adresse,omitempty"    sql:"complement_adresse"       csv:"complement_adresse"`
	NumVoie              string     `input:"numeroVoieEtablissement"                     json:"numero_voie,omitempty"           sql:"numero_voie"              csv:"numéro_voie"`
	IndRep               string     `input:"indiceRepetitionEtablissement"               json:"indrep,omitempty"                sql:"indrep"                   csv:"indrep"`
	TypeVoie             string     `input:"typeVoieEtablissement"                       json:"type_voie,omitempty"             sql:"type_voie"                csv:"type_voie"`
	Voie                 string     `input:"libelleVoieEtablissement"                    json:"voie,omitempty"                  sql:"voie"                     csv:"voie"`
	Commune              string     `input:"libelleCommuneEtablissement"                 json:"commune,omitempty"               sql:"commune"                  csv:"commune"`
	CommuneEtranger      string     `input:"libelleCommuneEtrangerEtablissement"         json:"commune_etranger,omitempty"      sql:"commune_etranger"         csv:"commune_étranger"`
	DistributionSpeciale string     `input:"distributionSpecialeEtablissement"           json:"distribution_speciale,omitempty" sql:"distribution_speciale"    csv:"distribution_speciale"`
	CodeCommune          string     `input:"codeCommuneEtablissement"                    json:"code_commune,omitempty"          sql:"code_commune"             csv:"code_commune"`
	CodeCedex            string     `input:"codeCedexEtablissement"                      json:"code_cedex,omitempty"            sql:"code_cedex"               csv:"code_cedex"`
	Cedex                string     `input:"libelleCedexEtablissement"                   json:"cedex,omitempty"                 sql:"cedex"                    csv:"cedex"`
	CodePaysEtranger     string     `input:"codePaysEtrangerEtablissement"               json:"code_pays_etranger,omitempty"    sql:"code_pays_etranger"       csv:"code_pays_étranger"`
	PaysEtranger         string     `input:"libellePaysEtrangerEtablissement"            json:"pays_etranger,omitempty"         sql:"pays_etranger"            csv:"pays_étranger"`
	CodePostal           string     `input:"codePostalEtablissement"                     json:"code_postal,omitempty"           sql:"code_postal"              csv:"code_postal"`
	Departement          string     `                                                    json:"departement,omitempty"           sql:"departement"              csv:"département"`
	APE                  string     `                                                    json:"ape,omitempty"                   sql:"ape"                      csv:"ape"`
	CodeActivite         string     `input:"activitePrincipaleEtablissement"             json:"code_activite,omitempty"         sql:"code_activite"            csv:"code_activité"`
	NomenActivite        string     `input:"nomenclatureActivitePrincipaleEtablissement" json:"nomen_activite,omitempty"        sql:"nomenclature_activite"    csv:"nomenclature_activité"`
	Creation             *time.Time `input:"dateCreationEtablissement"                   json:"date_creation,omitempty"         sql:"date_creation"            csv:"création"`
	Longitude            float64    `input:"longitude"                                   json:"longitude,omitempty"             sql:"longitude"                csv:"longitude"`
	Latitude             float64    `input:"latitude"                                    json:"latitude,omitempty"              sql:"latitude"                 csv:"latitude"`
	EstActif             bool       `input:"etatAdministratifEtablissement"              json:"est_actif"                       sql:"est_actif"                csv:"est_actif"`
}

// Key id de l'objet",
func (sirene Sirene) Key() string {
	return sirene.Siren + sirene.Nic
}

// Type de données
func (sirene Sirene) Type() engine.ParserType {
	return engine.Sirene
}

// Scope de l'objet
func (sirene Sirene) Scope() engine.Scope {
	return engine.ScopeEtablissement
}

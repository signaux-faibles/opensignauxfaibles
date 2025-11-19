package engine

// ParserType is the type used by all constants like ADMIN_URSSAF, APCONSO, etc...
type ParserType string

// These constants represent types supported by our data integration process.
// See https://documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	// There is actually no parser for AdminUrssaf, but it is used for an ID to
	// Siret Mapping
	AdminUrssaf ParserType = "admin_urssaf"
	Apconso     ParserType = "apconso"
	Apdemande   ParserType = "apdemande"
	Bdf         ParserType = "bdf"
	Ccsf        ParserType = "ccsf"
	Cotisation  ParserType = "cotisation"
	Debit       ParserType = "debit"
	Delai       ParserType = "delai"
	Diane       ParserType = "diane"
	Effectif    ParserType = "effectif"
	EffectifEnt ParserType = "effectif_ent"
	Filter      ParserType = "filter"
	Procol      ParserType = "procol"
	Sirene      ParserType = "sirene"
	SireneUl    ParserType = "sirene_ul"
	SireneHisto ParserType = "sirene_histo"
	Paydex      ParserType = "paydex"
	Ellisphere  ParserType = "ellisphere"
)

// Scope represents the type of entity: company (entreprise) or establishment (etablissement)
type Scope string

const (
	ScopeEntreprise    Scope = "entreprise"
	ScopeEtablissement Scope = "etablissement"
)


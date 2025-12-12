package engine

// ParserType defines a specific file format that can be parsed
type ParserType string

// These constants represent types supported by our data integration process.
// See https://documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	// There is actually no parser for AdminUrssaf, but it is used for an ID to
	// Siret Mapping
	Apconso     ParserType = "apconso"
	Apdemande   ParserType = "apdemande"
	Ccsf        ParserType = "ccsf"
	Cotisation  ParserType = "cotisation"
	Debit       ParserType = "debit"
	Delai       ParserType = "delai"
	Effectif    ParserType = "effectif"
	EffectifEnt ParserType = "effectif_ent"
	Filter      ParserType = "filter"
	Procol      ParserType = "procol"
	Sirene      ParserType = "sirene"
	SireneUl    ParserType = "sirene_ul"
	SireneHisto ParserType = "sirene_histo"
)

// Scope represents the type of entity: company (entreprise) or establishment (etablissement)
type Scope string

const (
	ScopeEntreprise    Scope = "entreprise"
	ScopeEtablissement Scope = "etablissement"
)

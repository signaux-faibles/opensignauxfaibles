package base

// ValidFileType is the type used by all constants like ADMIN_URSSAF, APCONSO, etc...
type ValidFileType string

// These constants represent types supported by our data integration process.
// See https://documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	AdminUrssaf ValidFileType = "admin_urssaf"
	Apconso     ValidFileType = "apconso"
	Apdemande   ValidFileType = "apdemande"
	Bdf         ValidFileType = "bdf"
	Ccsf        ValidFileType = "ccsf"
	Cotisation  ValidFileType = "cotisation"
	Debit       ValidFileType = "debit"
	Delai       ValidFileType = "delai"
	Diane       ValidFileType = "diane"
	Effectif    ValidFileType = "effectif"
	EffectifEnt ValidFileType = "effectif_ent"
	Filter      ValidFileType = "filter"
	Procol      ValidFileType = "procol"
	Sirene      ValidFileType = "sirene"
	SireneUl    ValidFileType = "sirene_ul"
	Paydex      ValidFileType = "paydex"
	Ellisphere  ValidFileType = "ellisphere"
)

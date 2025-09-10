package prepareimport

import (
	"regexp"
)

// ExtractFileTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractFileTypeFromFilename(filename string) ValidFileType {
	possiblyGzFilename := regexp.MustCompile(`^(.*)\.gz$`).ReplaceAllString(filename, `$1`)
	switch {
	case filename == "consommation_ap.csv":
		return apconso
	case filename == "demande_ap.csv":
		return apdemande
	case possiblyGzFilename == "sigfaible_etablissement_utf8.csv":
		return adminUrssaf
	case possiblyGzFilename == "sigfaible_effectif_siren.csv":
		return effectifEnt
	case possiblyGzFilename == "sigfaible_pcoll.csv":
		return procol
	case possiblyGzFilename == "sigfaible_cotisdues.csv":
		return cotisation
	case possiblyGzFilename == "sigfaible_delais.csv":
		return delai
	case possiblyGzFilename == "sigfaible_ccsf.csv":
		return ccsf
	case filename == "sireneUL.csv":
		return sireneUl
	case filename == "StockEtablissement_utf8_geo.csv":
		return sirene
	case mentionsDebits.MatchString(filename):
		return debit
	case hasDianePrefix.MatchString(filename):
		return diane
	case mentionsEffectif.MatchString(filename):
		return effectif
	case hasFilterPrefix.MatchString(filename):
		return filter
	case isRetroPaydex.MatchString(filename):
		return paydex
	case isEllisphere.MatchString(filename):
		return ellisphere
	default:
		return ""
	}
}

// These constants represent types supported by our data integration process.
// See https://documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
const (
	adminUrssaf ValidFileType = "admin_urssaf"
	apconso     ValidFileType = "apconso"
	apdemande   ValidFileType = "apdemande"
	bdf         ValidFileType = "bdf"
	ccsf        ValidFileType = "ccsf"
	cotisation  ValidFileType = "cotisation"
	debit       ValidFileType = "debit"
	delai       ValidFileType = "delai"
	diane       ValidFileType = "diane"
	effectif    ValidFileType = "effectif"
	effectifEnt ValidFileType = "effectif_ent"
	filter      ValidFileType = "filter"
	procol      ValidFileType = "procol"
	sirene      ValidFileType = "sirene"
	sireneUl    ValidFileType = "sirene_ul"
	paydex      ValidFileType = "paydex"
	ellisphere  ValidFileType = "ellisphere"
)

// ValidFileType is the type used by all constants like ADMIN_URSSAF, APCONSO, etc...
type ValidFileType string

var hasDianePrefix = regexp.MustCompile(`^[Dd]iane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)
var isRetroPaydex = regexp.MustCompile(`^E_[0-9]{12}_Retro-Paydex_[0-9]{8}.csv$`)
var isEllisphere = regexp.MustCompile(`^Ellisphère-Tête de groupe-[^.]*.xlsx$`)

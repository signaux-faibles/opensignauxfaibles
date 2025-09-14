package prepareimport

import (
	"opensignauxfaibles/lib/base"
	"regexp"
)

// ExtractFileTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractFileTypeFromFilename(filename string) base.ValidFileType {
	possiblyGzFilename := regexp.MustCompile(`^(.*)\.gz$`).ReplaceAllString(filename, `$1`)
	switch {
	case filename == "consommation_ap.csv":
		return base.Apconso
	case filename == "demande_ap.csv":
		return base.Apdemande
	case possiblyGzFilename == "sigfaible_etablissement_utf8.csv":
		return base.AdminUrssaf
	case possiblyGzFilename == "sigfaible_effectif_siren.csv":
		return base.EffectifEnt
	case possiblyGzFilename == "sigfaible_pcoll.csv":
		return base.Procol
	case possiblyGzFilename == "sigfaible_cotisdues.csv":
		return base.Cotisation
	case possiblyGzFilename == "sigfaible_delais.csv":
		return base.Delai
	case possiblyGzFilename == "sigfaible_ccsf.csv":
		return base.Ccsf
	case filename == "sireneUL.csv":
		return base.SireneUl
	case filename == "StockEtablissement_utf8_geo.csv":
		return base.Sirene
	case mentionsDebits.MatchString(filename):
		return base.Debit
	case hasDianePrefix.MatchString(filename):
		return base.Diane
	case mentionsEffectif.MatchString(filename):
		return base.Effectif
	case hasFilterPrefix.MatchString(filename):
		return base.Filter
	case isRetroPaydex.MatchString(filename):
		return base.Paydex
	case isEllisphere.MatchString(filename):
		return base.Ellisphere
	default:
		return ""
	}
}

var hasDianePrefix = regexp.MustCompile(`^[Dd]iane`)
var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter_`)
var isRetroPaydex = regexp.MustCompile(`^E_[0-9]{12}_Retro-Paydex_[0-9]{8}.csv$`)
var isEllisphere = regexp.MustCompile(`^Ellisphère-Tête de groupe-[^.]*.xlsx$`)

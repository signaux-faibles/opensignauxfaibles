package prepareimport

import (
	"opensignauxfaibles/lib/engine"
	"regexp"
)

// ExtractParserTypeFromFilename returns a file type from filename, or empty string for unsupported file names
func ExtractParserTypeFromFilename(filename string) engine.ParserType {
	possiblyGzFilename := regexp.MustCompile(`^(.*)\.gz$`).ReplaceAllString(filename, `$1`)
	switch {
	case filename == "consommation_ap.csv":
		return engine.Apconso
	case filename == "demande_ap.csv":
		return engine.Apdemande
	case possiblyGzFilename == "sigfaible_etablissement_utf8.csv":
		return engine.AdminUrssaf
	case possiblyGzFilename == "sigfaible_effectif_siren.csv":
		return engine.EffectifEnt
	case possiblyGzFilename == "sigfaible_pcoll.csv":
		return engine.Procol
	case possiblyGzFilename == "sigfaible_cotisdues.csv":
		return engine.Cotisation
	case possiblyGzFilename == "sigfaible_delais.csv":
		return engine.Delai
	case possiblyGzFilename == "sigfaible_ccsf.csv":
		return engine.Ccsf
	case filename == "sireneUL.csv":
		return engine.SireneUl
	case filename == "StockEtablissement_utf8_geo.csv":
		return engine.Sirene
	case mentionsDebits.MatchString(filename):
		return engine.Debit
	case hasDianePrefix.MatchString(filename):
		return engine.Diane
	case mentionsEffectif.MatchString(filename):
		return engine.Effectif
	case hasFilterPrefix.MatchString(filename):
		return engine.Filter
	case isRetroPaydex.MatchString(filename):
		return engine.Paydex
	case isEllisphere.MatchString(filename):
		return engine.Ellisphere
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

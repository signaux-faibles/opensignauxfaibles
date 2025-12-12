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
	case filename == "StockEtablissementHistorique_utf8.csv":
		return engine.SireneHisto
	case mentionsDebits.MatchString(filename):
		return engine.Debit
	case mentionsEffectif.MatchString(filename):
		return engine.Effectif
	case hasFilterPrefix.MatchString(filename):
		return engine.Filter
	default:
		return ""
	}
}

var mentionsEffectif = regexp.MustCompile(`effectif_`)
var mentionsDebits = regexp.MustCompile(`_debits`)
var hasFilterPrefix = regexp.MustCompile(`^filter`)

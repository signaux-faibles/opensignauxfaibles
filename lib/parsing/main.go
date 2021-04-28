// Ce paquet centralise les parseurs de fichiers pour les mettre à
// disposition de `engine`.

package parsing

import (
	"errors"

	"github.com/signaux-faibles/opensignauxfaibles/lib/apconso"
	"github.com/signaux-faibles/opensignauxfaibles/lib/apdemande"
	"github.com/signaux-faibles/opensignauxfaibles/lib/bdf"
	"github.com/signaux-faibles/opensignauxfaibles/lib/diane"
	"github.com/signaux-faibles/opensignauxfaibles/lib/ellisphere"
	"github.com/signaux-faibles/opensignauxfaibles/lib/marshal"
	"github.com/signaux-faibles/opensignauxfaibles/lib/paydex"
	"github.com/signaux-faibles/opensignauxfaibles/lib/sirene"
	sireneul "github.com/signaux-faibles/opensignauxfaibles/lib/sirene_ul"
	"github.com/signaux-faibles/opensignauxfaibles/lib/urssaf"
)

// RegisteredParsers liste des parsers disponibles
// Note: penser à tenir à jour la table des formats, dans la documentation:
// https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
var registeredParsers = map[string]marshal.Parser{
	"debit":        urssaf.ParserDebit,
	"ccsf":         urssaf.ParserCCSF,
	"cotisation":   urssaf.ParserCotisation,
	"admin_urssaf": urssaf.ParserCompte,
	"delai":        urssaf.ParserDelai,
	"effectif":     urssaf.ParserEffectif,
	"effectif_ent": urssaf.ParserEffectifEnt,
	"procol":       urssaf.ParserProcol,
	"apconso":      apconso.Parser,
	"apdemande":    apdemande.Parser,
	"bdf":          bdf.Parser,
	"sirene":       sirene.Parser,
	"sirene_ul":    sireneul.Parser,
	"diane":        diane.Parser,
	"ellisphere":   ellisphere.Parser,
	"paydex":       paydex.ParserPaydex,
}

// IsSupportedParser retourne true si un parseur est défini pour le fileType spécifié
// ou si le type est "filter". (cf issue #354)
func IsSupportedParser(fileType string) bool {
	return fileType == "filter" || registeredParsers[fileType] != nil
}

// ResolveParsers sélectionne, vérifie et charge les parsers.
func ResolveParsers(parserNames []string) ([]marshal.Parser, error) {
	var parsers []marshal.Parser
	if parserNames == nil {
		for _, fileParser := range registeredParsers {
			parsers = append(parsers, fileParser)
		}
	} else {
		for _, fileType := range parserNames {
			if fileParser, ok := registeredParsers[fileType]; ok {
				parsers = append(parsers, fileParser)
			} else {
				return parsers, errors.New(fileType + " n'est pas un parser reconnu.")
			}
		}
	}
	return parsers, nil
}

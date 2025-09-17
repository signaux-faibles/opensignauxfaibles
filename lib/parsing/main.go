// Ce paquet centralise les parseurs de fichiers pour les mettre à
// disposition de `engine`.

package parsing

import (
	"errors"

	"opensignauxfaibles/lib/apconso"
	"opensignauxfaibles/lib/apdemande"
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/marshal"
	"opensignauxfaibles/lib/sirene"
	sireneul "opensignauxfaibles/lib/sirene_ul"
	"opensignauxfaibles/lib/urssaf"
)

// RegisteredParsers liste des parsers disponibles
// Note: penser à tenir à jour la table des formats, dans la documentation:
// https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
var registeredParsers = map[base.ParserType]marshal.Parser{
	base.Debit:       urssaf.ParserDebit,
	base.Ccsf:        urssaf.ParserCCSF,
	base.Cotisation:  urssaf.ParserCotisation,
	base.AdminUrssaf: urssaf.ParserCompte,
	base.Delai:       urssaf.ParserDelai,
	base.Effectif:    urssaf.ParserEffectif,
	base.EffectifEnt: urssaf.ParserEffectifEnt,
	base.Procol:      urssaf.ParserProcol,
	base.Apconso:     apconso.Parser,
	base.Apdemande:   apdemande.Parser,
	base.Sirene:      sirene.Parser,
	base.SireneUl:    sireneul.Parser,
}

// IsSupportedParser retourne true si un parseur est défini pour le fileType spécifié
// ou si le type est "filter". (cf issue #354)
func IsSupportedParser(fileType base.ParserType) bool {
	return fileType == "filter" || registeredParsers[fileType] != nil
}

// ResolveParsers sélectionne, vérifie et charge les parsers.
func ResolveParsers(parserNames []base.ParserType) ([]marshal.Parser, error) {
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
				return parsers, errors.New(string(fileType) + " n'est pas un parser reconnu.")
			}
		}
	}
	return parsers, nil
}

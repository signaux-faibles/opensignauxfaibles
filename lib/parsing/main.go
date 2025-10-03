// Ce paquet centralise les parseurs de fichiers pour les mettre à
// disposition de `engine`.

package parsing

import (
	"opensignauxfaibles/lib/base"
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing/apconso"
	"opensignauxfaibles/lib/parsing/apdemande"
	"opensignauxfaibles/lib/parsing/sirene"
	sireneul "opensignauxfaibles/lib/parsing/sirene_ul"
	"opensignauxfaibles/lib/parsing/urssaf"
)

// Implements engine.ParserRegistry
type ParserRegistry map[base.ParserType]engine.Parser

func (pr ParserRegistry) Resolve(parserType base.ParserType) engine.Parser {
	if parser, ok := pr[parserType]; ok {
		return parser
	}
	return nil
}

func (pr ParserRegistry) All() []engine.Parser {
	var parsers []engine.Parser
	for _, parser := range pr {
		parsers = append(parsers, parser)
	}
	return parsers
}

// DefaultParsers liste des parsers disponibles
// Note: penser à tenir à jour la table des formats, dans la documentation:
// https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport
var DefaultParsers = ParserRegistry{
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

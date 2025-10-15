// Package registry centralise les parseurs de fichiers pour les mettre à
// disposition de `engine`.
package registry

import (
	"opensignauxfaibles/lib/engine"
	"opensignauxfaibles/lib/parsing/apconso"
	"opensignauxfaibles/lib/parsing/apdemande"
	"opensignauxfaibles/lib/parsing/effectif"
	"opensignauxfaibles/lib/parsing/sirene"
	sireneul "opensignauxfaibles/lib/parsing/sirene_ul"
	"opensignauxfaibles/lib/parsing/urssaf"
)

// ParserRegistry implements engine.ParserRegistry
type ParserRegistry map[engine.ParserType]engine.Parser

func (pr ParserRegistry) Resolve(parserType engine.ParserType) engine.Parser {
	if pr == nil {
		return nil
	}

	if parser, ok := pr[parserType]; ok {
		return parser
	}
	return nil
}

func (pr ParserRegistry) All() []engine.Parser {
	if pr == nil {
		return []engine.Parser{}
	}

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
	engine.Debit:       urssaf.NewDebitParser(),
	engine.Ccsf:        urssaf.NewCCSFParser(),
	engine.Cotisation:  urssaf.NewCotisationParser(),
	engine.Delai:       urssaf.NewDelaiParser(),
	engine.Effectif:    effectif.NewEffectifParser(),
	engine.EffectifEnt: effectif.NewEffectifEntParser(),
	engine.Procol:      urssaf.NewProcolParser(),
	engine.Apconso:     apconso.NewApconsoParser(),
	engine.Apdemande:   apdemande.NewApdemandeParser(),
	engine.Sirene:      sirene.NewSireneParser(),
	engine.SireneUl:    sireneul.NewSireneULParser(),
}

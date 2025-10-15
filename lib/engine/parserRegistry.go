package engine

import (
	"fmt"
)

// ParserRegistry est une interface pour accéder à l'ensemble des parsers supportés
type ParserRegistry interface {
	// All retourne tous les parsers disponibles
	All() []Parser

	// Resolve retourne le parser associé à un type donné
	//
	// Retourne `nil` si ce parser n'existe pas
	Resolve(ParserType) Parser
}

// ResolveParsers récupère les parsers associés aux types donnés en entrée.
//
// Si aucun type n'est renseigné, tous les parsers sont retournés.
//
// Si un type absent du registre est requêté, une erreur est retournée.
func ResolveParsers(registry ParserRegistry, types []ParserType) ([]Parser, error) {
	var parsers []Parser

	if len(types) == 0 {
		return registry.All(), nil
	}

	for _, parserType := range types {
		parser := registry.Resolve(parserType)
		if parser == nil {
			return nil, fmt.Errorf("parser registry could not resolve parser \"%s\"", parserType)
		}
		parsers = append(parsers, parser)
	}
	return parsers, nil
}

package engine

import (
	"fmt"
	"opensignauxfaibles/lib/base"
)

type ParserRegistry interface {
	Resolve(base.ParserType) Parser
	All() []Parser
}

// ResolveParsers gets the appropriate parsers from a types array.
// If "types" is empty, all parsers are returned.
// A provided type that is missing from the registry triggers an error.
func ResolveParsers(registry ParserRegistry, types []base.ParserType) ([]Parser, error) {
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

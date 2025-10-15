package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DummyRegistry stores parsers in a map, with keys associated with the type
type DummyRegistry struct {
	parsers map[ParserType]Parser
}

func (r DummyRegistry) Resolve(t ParserType) Parser {
	return r.parsers[t]
}

func (r DummyRegistry) All() []Parser {
	var allParsers []Parser
	for _, v := range r.parsers {
		allParsers = append(allParsers, v)
	}

	return allParsers
}

func TestResolveParsers(t *testing.T) {
	type1 := ParserType("one")
	type2 := ParserType("two")
	parser1 := &dummyParser{parserType: type1}
	parser2 := &dummyParser{parserType: type2}

	registry := DummyRegistry{
		map[ParserType]Parser{
			parser1.Type(): parser1,
			parser2.Type(): parser2,
		},
	}

	testCases := []struct {
		name              string
		types             []ParserType
		expectErr         bool
		expectParserTypes []ParserType
	}{
		{
			"Aucun type renseigné (empty slice), tous les parsers sont retournés",
			[]ParserType{},
			false,
			[]ParserType{type1, type2},
		},
		{
			"Aucun type renseigné (`nil`), tous les parsers sont retournés",
			nil,
			false,
			[]ParserType{type1, type2},
		},
		{
			"Type 1 renseigné, parser 1 retourné",
			[]ParserType{type1},
			false,
			[]ParserType{type1},
		},
		{
			"Types 1 et 2 renseignés, parsers 1 et 2 retournés",
			[]ParserType{type1, type2},
			false,
			[]ParserType{type1, type2},
		},
		{
			"Type inconnu renseigné, erreur retournée",
			[]ParserType{"nonsense"},
			true,
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsers, err := ResolveParsers(registry, tc.types)

			if tc.expectErr {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
				assert.Len(t, parsers, len(tc.expectParserTypes))
				for _, p := range parsers {
					assert.Contains(t, tc.expectParserTypes, p.Type())
				}
			}
		})
	}
}

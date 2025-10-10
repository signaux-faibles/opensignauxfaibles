package engine

import (
	"opensignauxfaibles/lib/base"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DummyRegistry stores parsers in a map, with keys associated with the type
type DummyRegistry struct {
	parsers map[base.ParserType]Parser
}

func (r DummyRegistry) Resolve(t base.ParserType) Parser {
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
	type1 := base.ParserType("one")
	type2 := base.ParserType("two")
	parser1 := &dummyParser{parserType: type1}
	parser2 := &dummyParser{parserType: type2}

	registry := DummyRegistry{
		map[base.ParserType]Parser{
			parser1.Type(): parser1,
			parser2.Type(): parser2,
		},
	}

	testCases := []struct {
		name              string
		types             []base.ParserType
		expectErr         bool
		expectParserTypes []base.ParserType
	}{
		{
			"Aucun type renseigné (empty slice), tous les parsers sont retournés",
			[]base.ParserType{},
			false,
			[]base.ParserType{type1, type2},
		},
		{
			"Aucun type renseigné (`nil`), tous les parsers sont retournés",
			nil,
			false,
			[]base.ParserType{type1, type2},
		},
		{
			"Type 1 renseigné, parser 1 retourné",
			[]base.ParserType{type1},
			false,
			[]base.ParserType{type1},
		},
		{
			"Types 1 et 2 renseignés, parsers 1 et 2 retournés",
			[]base.ParserType{type1, type2},
			false,
			[]base.ParserType{type1, type2},
		},
		{
			"Type inconnu renseigné, erreur retournée",
			[]base.ParserType{"nonsense"},
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

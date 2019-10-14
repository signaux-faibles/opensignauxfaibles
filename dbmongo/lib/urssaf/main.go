package urssaf

import (
	"opensignauxfaibles/dbmongo/lib/engine"
)

// Parser fournit le contenu des fichiers urssaf
func Parser(cache engine.Cache, batch *engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	go func() {
		defer close(outputChannel)
		// TODO close properly eventChannel
		// 77fbb9b6-e2ff-438d-9396-2ce6277414c8

		functions := []engine.Parser{
			parseCCSF,
			parseCotisation,
			parseDebit,
			parseDelai,
			parseEffectif,
			parseProcol,
		}

		for _, f := range functions {
			outputs, events := f(cache, batch)
			go engine.PlugEvents(events, eventChannel)
			for o := range outputs {
				outputChannel <- o
			}
		}
	}()
	return outputChannel, eventChannel
}

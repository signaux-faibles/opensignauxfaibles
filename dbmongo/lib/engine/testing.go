package engine

import "github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"

// DiscardTuple supprime les évènements
func DiscardTuple(tuples chan base.Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

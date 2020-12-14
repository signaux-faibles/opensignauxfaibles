package engine

import "github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/marshal"

// DiscardTuple supprime les évènements
func DiscardTuple(tuples chan marshal.Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

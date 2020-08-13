package engine

// DiscardTuple supprime les évènements
func DiscardTuple(tuples chan Tuple) {
	go func() {
		for range tuples {
		}
	}()
}

package marshal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/signaux-faibles/opensignauxfaibles/dbmongo/lib/base"
)

// Parser fonction de traitement de données en entrée
type Parser func(Cache, *base.AdminBatch) (chan Tuple, chan Event)

// Tuple unité de donnée à insérer dans un type
type Tuple interface {
	Key() string
	Scope() string
	Type() string
}

// GetJSON sérialise un tuple au format JSON.
func GetJSON(tuple Tuple) ([]byte, error) {
	return json.MarshalIndent(tuple, "", "  ")
}

// LogProgress affiche le numéro de ligne en cours de parsing, toutes les 2s.
func LogProgress(lineNumber *int) (stop context.CancelFunc) {
	ctx, stop := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for range time.Tick(time.Second * 2) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			fmt.Printf("Reading csv line %d\n", *lineNumber)
		}
	}(ctx)
	return stop
}

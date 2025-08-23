// Ce fichier est responsable de collecter les messages et de les ajouter
// dans la collection Journal.

package engine

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"opensignauxfaibles/lib/marshal"
)

const EventTable string = "import_logs"

type EventSink interface {
	Process(ch chan marshal.Event) error
}

type PostgresEventSink struct {
	conn *pgxpool.Pool

	// Name of the table to which to write
	table string

	// Attach the processed events to a specific command
	command string

	// time of creation
	timestamp time.Time
}

func NewPostgresEventSink(conn *pgxpool.Pool, command string) EventSink {
	return &PostgresEventSink{conn, EventTable, command, time.Now()}
}

func (sink *PostgresEventSink) Process(ch chan marshal.Event) error {
	for range ch {
	}
	return nil
}

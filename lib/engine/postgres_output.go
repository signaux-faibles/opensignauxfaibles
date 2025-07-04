package engine

import (
	"opensignauxfaibles/lib/marshal"

	"github.com/jackc/pgx/v5"
)

// PostgresOutputStreamer writes the output to csv files. It implements `OutputStreamer`
// If writer is nil, it will stream into csv files in the "relativeDirPath"
// directory.
// Otherwise, it will stream to the io.Writer.
type PostgresOutputStreamer struct {
	conn *pgx.Conn
}

func NewPostgresOutputStreamer(conn *pgx.Conn) *PostgresOutputStreamer {
	return &PostgresOutputStreamer{conn}
}

func (out *PostgresOutputStreamer) Stream(ch chan marshal.Tuple) error {
	firstTuple, ok := <-ch // to retrieve the type of data
	if !ok {
		return nil // no data to process
	}

}

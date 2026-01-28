package export

import (
	"context"
	"fmt"
	"opensignauxfaibles/lib/db"
)

var viewsToExport = []string{
	db.ViewCleanSirene,
	db.ViewCleanSireneUL,
}

type Exporter struct {
	path string
	conn db.Pool
}

func NewExporter(path string, conn db.Pool) *Exporter {
	return &Exporter{path, conn}
}

// CleanViews exports all clean views
func (exp *Exporter) CleanViews() error {
	ctx := context.TODO()

	for _, view := range viewsToExport {
		ToCsv(ctx, exp.path, view, view+".csv", exp.conn)
	}

	return nil
}

func ToCsv(ctx context.Context, path string, view string, tableName string, conn db.Pool) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`
    COPY (
      SELECT *
      FROM %s
    ) TO %s
    WITH (FORMAT CSV, HEADER, DELIMITER ',');`, view, tableName),
	)
	return err
}

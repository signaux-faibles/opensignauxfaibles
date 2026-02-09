// Package export defines the export of the postgres views.
// It differs from the sinks in that it exports refined data (clean views)
// instead of staging data (`stg_*` tables).
package export

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"opensignauxfaibles/lib/db"

	"golang.org/x/sync/errgroup"
)

var viewsToExport = []string{
	db.ViewSireneHisto,
	db.ViewEffectifEnt,
	db.ViewCotisation,
	db.ViewDebit,
	db.ViewDelai,
	db.ViewProcol,
	db.ViewAp,
}

type Exporter struct {
	path string
	conn db.Pool
}

// NewExporter initialise la fonctionnalité d'export de la base de donnée
// dont on fournit une connexion `conn` vers un répertoire *sur le serveur ou
// conteneur de base de données* `path`.
// Il est supposé que le répertoire existe, sans vérification.
func NewExporter(path string, conn db.Pool) *Exporter {
	return &Exporter{path, conn}
}

// CleanViews exports all clean views
func (exp *Exporter) CleanViews() error {
	g, ctx := errgroup.WithContext(context.Background())

	err := os.MkdirAll(exp.path, 0644)
	if err != nil {
		return err
	}
	dirAbsPath, err := filepath.Abs(exp.path)
	if err != nil {
		return err
	}

	for _, view := range viewsToExport {
		// run export in parallel
		g.Go(func() error {
			fileAbsPath := filepath.Join(dirAbsPath, view+".csv")
			// Fail if file already exist
			_, err := os.Stat(fileAbsPath)
			if err == nil {
				slog.Warn(fmt.Sprintf("Le fichier %s existe déjà et n'est pas réexporté", fileAbsPath))
				return nil
			}

			f, err := os.Create(fileAbsPath)
			if err != nil {
				return err
			}

			if err = toWriter(ctx, f, view, exp.conn); err != nil {
				return fmt.Errorf("export of %s failed: %w", view, err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("Erreur:", err)
	}

	return nil
}

// toWriter copies the whole content of the `view` (on database with
// connection `conn`) to `w`
func toWriter(ctx context.Context, w io.Writer, view string, conn db.Pool) error {
	poolConn, err := conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer poolConn.Release()
	_, err = poolConn.Conn().PgConn().CopyTo(
		ctx,
		w,
		fmt.Sprintf(`COPY (SELECT * FROM %s) TO STDOUT WITH (FORMAT CSV, HEADER, DELIMITER ',');`, view),
	)
	return err
}

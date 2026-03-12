package filter

import (
	"context"
	"fmt"
	"log/slog"
	"opensignauxfaibles/lib/sinks"
	"os"
	"path/filepath"
	"time"
)

// FilterStats tracks the number of companies at each filtering stage
type FilterStats struct {
	TotalCompanies             int
	AfterEffectifFilter        int
	ExcludedByEffectifFilter   int
	AfterRecentDataFilter      int
	ExcludedByRecentDataFilter int
}

// SQLFilterStats contains statistics about SQL-based filtering
type SQLFilterStats struct {
	BeforeSQLFilter   int
	AfterSQLFilter    int
	ExcludedBySQL     int
	ExcludedByJuridic int
	ExcludedByNAF     int
	ExcludedByForeign int
}

// generateFilterReport creates a detail-filtres.md file with detailed filter statistics
func generateFilterReport(stats FilterStats, batchKey string) error {
	exportPath := filepath.Join(sinks.DefaultExportPath, batchKey)
	filename := filepath.Join(exportPath, "detail-filtres.md")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create filter report file: %w", err)
	}
	defer f.Close()

	// Write header
	fmt.Fprintf(f, "# Détail des filtres appliqués\n\n")
	fmt.Fprintf(f, "Généré le : %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "## Paramètres\n\n")
	fmt.Fprintf(f, "- **Fenêtre d'observation effectifs** : %d mois\n", DefaultNbMois)
	fmt.Fprintf(f, "- **Seuil minimum d'effectif** : %d employés\n", DefaultMinEffectif)
	fmt.Fprintf(f, "- **Nombre de mois avec données récentes requis** : %d mois\n\n", DefaultNbMoisRecentData)

	// Write filter statistics
	fmt.Fprintf(f, "## Résultats du filtrage\n\n")
	fmt.Fprintf(f, "| Étape | Entreprises | Exclues | Taux d'exclusion |\n")
	fmt.Fprintf(f, "|-------|-------------|---------|------------------|\n")
	fmt.Fprintf(f, "| Périmètre initial | %d | - | - |\n", stats.TotalCompanies)

	// Filter 1: Effectif threshold
	effectifExclusionRate := 0.0
	if stats.TotalCompanies > 0 {
		effectifExclusionRate = float64(stats.ExcludedByEffectifFilter) / float64(stats.TotalCompanies) * 100
	}
	fmt.Fprintf(f, "| Suppression des entreprises ne dépassant pas %d employés sur les %d derniers mois | %d | %d | %.2f%% |\n",
		DefaultMinEffectif, DefaultNbMois, stats.AfterEffectifFilter, stats.ExcludedByEffectifFilter, effectifExclusionRate)

	// Filter 2: Recent data
	recentDataExclusionRate := 0.0
	if stats.AfterEffectifFilter > 0 {
		recentDataExclusionRate = float64(stats.ExcludedByRecentDataFilter) / float64(stats.AfterEffectifFilter) * 100
	}
	fmt.Fprintf(f, "| Suppression des entreprises n'ayant pas de données d'effectif sur les %d derniers mois | %d | %d | %.2f%% |\n",
		DefaultNbMoisRecentData, stats.AfterRecentDataFilter, stats.ExcludedByRecentDataFilter, recentDataExclusionRate)

	// Note about subsequent filters
	fmt.Fprintf(f, "\n## Filtres appliqués ultérieurement (en SQL)\n\n")
	fmt.Fprintf(f, "Les filtres suivants sont appliqués au niveau de la base de données (entre les couches `stg_*` et `clean_*`) :\n\n")
	fmt.Fprintf(f, "1. **Suppression de certaines formes juridiques**\n")
	fmt.Fprintf(f, "   - Établissements publics nationaux et locaux\n")
	fmt.Fprintf(f, "   - Communes, départements, associations\n")
	fmt.Fprintf(f, "   - Syndicats et groupements de collectivités territoriales\n\n")
	fmt.Fprintf(f, "2. **Suppression de certaines activités principales** (codes NAF)\n")
	fmt.Fprintf(f, "   - Administration publique (84.XX)\n")
	fmt.Fprintf(f, "   - Enseignement (85.XX)\n")
	fmt.Fprintf(f, "   - Organisations associatives (94.XX)\n")
	fmt.Fprintf(f, "   - Services financiers et assurance (64.XX, 65.XX, 66.XX)\n")
	fmt.Fprintf(f, "   - Organisations extraterritoriales (99.XX)\n\n")
	fmt.Fprintf(f, "3. **Suppression des entreprises domiciliées hors de France**\n")
	fmt.Fprintf(f, "   - Exclusion des sièges avec département vide\n\n")

	// Final note
	fmt.Fprintf(f, "## Périmètre final\n\n")
	fmt.Fprintf(f, "Le périmètre final est disponible dans :\n")
	fmt.Fprintf(f, "- **Base de données** : table `sfdata.clean_filter`\n")
	fmt.Fprintf(f, "- **Export CSV** : fichier `clean_filter.csv`\n")

	slog.Info("filter report generated", "filename", filename)
	return nil
}

// UpdateFilterReportWithSQLStats appends SQL filter statistics to the existing report
func UpdateFilterReportWithSQLStats(connGetter func() (interface{}, error), batchKey string) error {
	conn, err := connGetter()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Type assertion to get the actual connection type
	type queryExecutor interface {
		QueryRow(ctx context.Context, query string, args ...interface{}) interface{ Scan(...interface{}) error }
	}

	db, ok := conn.(queryExecutor)
	if !ok {
		return fmt.Errorf("connection does not support QueryRow")
	}

	ctx := context.Background()

	// Query SQL filter statistics
	var stats SQLFilterStats

	// Count before SQL filters (stg_filter_import)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM stg_filter_import").Scan(&stats.BeforeSQLFilter)
	if err != nil {
		return fmt.Errorf("failed to count stg_filter_import: %w", err)
	}

	// Count excluded by SQL filters (siren_blacklist)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM siren_blacklist").Scan(&stats.ExcludedBySQL)
	if err != nil {
		return fmt.Errorf("failed to count siren_blacklist: %w", err)
	}

	// Count after SQL filters (clean_filter)
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM clean_filter").Scan(&stats.AfterSQLFilter)
	if err != nil {
		return fmt.Errorf("failed to count clean_filter: %w", err)
	}

	// Count detailed exclusions by category
	detailQuery := `
		WITH excluded_categories AS (
			SELECT ARRAY[
				'4110', '4120', '4140', '4160', '7210', '7220', '7346', '7348', 
				'7366', '7373', '7379', '7383', '7389', '7410', '7430', '7470', '7490'
			] AS categories
		),
		excluded_activities AS (
			SELECT ARRAY[
				'84.11Z', '84.12Z', '84.13Z', '84.21Z', '84.22Z', '84.23Z', '84.24Z', '84.25Z',
				'84.30A', '84.30B', '84.30C', '85.10Z', '85.20Z', '85.31Z', '85.32Z', '85.41Z',
				'85.42Z', '94.11Z', '94.12Z', '94.20Z', '94.91Z', '94.92Z', '64.11Z', '64.19Z',
				'64.30Z', '64.91Z', '64.92Z', '64.99Z', '65.11Z', '65.12Z', '65.20Z', '65.30Z',
				'66.11Z', '66.12Z', '66.19A', '66.19B', '66.21Z', '66.22Z', '66.29Z', '66.30Z', '99.00Z'
			] AS activities
		)
		SELECT 
			COUNT(*) FILTER (WHERE sirene_ul.categorie_juridique::text = ANY (ec.categories)) as juridic_count,
			COUNT(*) FILTER (WHERE sirene_ul.activite_principale::text = ANY (ea.activities)) as naf_count,
			COUNT(*) FILTER (WHERE sirene.siren IS NULL) as foreign_count
		FROM stg_filter_import fp
		LEFT JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren::text = fp.siren::text
		LEFT JOIN stg_sirene sirene ON sirene.siren = fp.siren 
			AND sirene.siege = true
			AND sirene.departement <> ''
		CROSS JOIN excluded_categories ec
		CROSS JOIN excluded_activities ea
		WHERE fp.siren IN (SELECT siren FROM siren_blacklist)
	`

	err = db.QueryRow(ctx, detailQuery).Scan(&stats.ExcludedByJuridic, &stats.ExcludedByNAF, &stats.ExcludedByForeign)
	if err != nil {
		return fmt.Errorf("failed to get detailed exclusion counts: %w", err)
	}

	// Open report file in append mode
	exportPath := filepath.Join(sinks.DefaultExportPath, batchKey)
	filename := filepath.Join(exportPath, "detail-filtres.md")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open filter report for update: %w", err)
	}
	defer f.Close()

	// Append SQL filter statistics
	fmt.Fprintf(f, "\n---\n\n")
	fmt.Fprintf(f, "**Mise à jour après application des filtres SQL** (généré le : %s)\n\n", time.Now().Format("2006-01-02 15:04:05"))

	sqlExclusionRate := 0.0
	if stats.BeforeSQLFilter > 0 {
		sqlExclusionRate = float64(stats.ExcludedBySQL) / float64(stats.BeforeSQLFilter) * 100
	}

	fmt.Fprintf(f, "| Suppression via filtres SQL (formes juridiques, NAF, domiciliation) | %d | %d | %.2f%% |\n",
		stats.AfterSQLFilter, stats.ExcludedBySQL, sqlExclusionRate)

	fmt.Fprintf(f, "\n### Détail des exclusions SQL\n\n")
	fmt.Fprintf(f, "- **Formes juridiques exclues** : %d entreprises\n", stats.ExcludedByJuridic)
	fmt.Fprintf(f, "- **Codes NAF exclus** : %d entreprises\n", stats.ExcludedByNAF)
	fmt.Fprintf(f, "- **Domiciliation hors France** : %d entreprises\n", stats.ExcludedByForeign)
	fmt.Fprintf(f, "\n*Note : Une entreprise peut être comptée dans plusieurs catégories si elle remplit plusieurs critères d'exclusion.*\n")

	fmt.Fprintf(f, "\n### Périmètre final\n\n")
	fmt.Fprintf(f, "**%d entreprises** dans le périmètre final après tous les filtres.\n\n", stats.AfterSQLFilter)

	totalExcluded := stats.BeforeSQLFilter - stats.AfterSQLFilter
	finalExclusionRate := 0.0
	if stats.BeforeSQLFilter > 0 {
		finalExclusionRate = float64(totalExcluded) / float64(stats.BeforeSQLFilter) * 100
	}
	fmt.Fprintf(f, "**Taux d'exclusion global** : %.2f%% (%d entreprises exclues sur %d initiales)\n",
		finalExclusionRate, totalExcluded, stats.BeforeSQLFilter)

	slog.Info("filter report updated with SQL statistics", "filename", filename)
	return nil
}

package db

// Table names
const (
	// Regular tables
	TableImportLogs      = "import_logs"
	TableLabelsMotif     = "labels_motif_recours"
	TableStgApconso      = "stg_apconso"
	TableStgApdemande    = "stg_apdemande"
	TableStgCotisation   = "stg_cotisation"
	TableStgDebit        = "stg_debit"
	TableStgDelai        = "stg_delai"
	TableStgEffectif     = "stg_effectif"
	TableStgEffectifEnt  = "stg_effectif_ent"
	TableStgSirene       = "stg_sirene"
	TableStgSireneUl     = "stg_sirene_ul"
	TableStgFilterImport = "stg_filter_import"
)

// Materialized views
const (
	ViewStgApdemandePeriod = "stg_apdemande_by_period"
	ViewCleanFilter        = "clean_filter"
)

// AllTables returns all table names that should be truncated during test cleanup
func AllTables() []string {
	return []string{
		VersionTable, // migrations table

		TableImportLogs,
		TableLabelsMotif,
		TableStgApconso,
		TableStgApdemande,
		TableStgCotisation,
		TableStgDebit,
		TableStgDelai,
		TableStgEffectif,
		TableStgEffectifEnt,
		TableStgSirene,
		TableStgSireneUl,
	}
}

// AllMaterializedViews returns all materialized view names that should be refreshed during test cleanup
func AllMaterializedViews() []string {
	return []string{
		ViewStgApdemandePeriod,
		ViewCleanFilter,
	}
}

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
	TableStgProcol       = "stg_procol"
	TableStgSireneHisto  = "stg_sirene_histo"
)

// Views
const (
	ViewSirene      = "clean_sirene"
	ViewSireneUL    = "clean_sirene_ul"
	ViewSireneHisto = "clean_sirene_histo"
	ViewEffectif    = "clean_effectif"
	ViewEffectifEnt = "clean_effectif_ent"
	ViewCotisation  = "clean_cotisation"
	ViewDelai       = "clean_delai"
	ViewProcol      = "clean_procol"
	ViewAp          = "clean_ap"
	ViewFilter      = "clean_filter"
)

// Materialized views
const (
	ViewStgApdemandePeriod = "stg_apdemande_by_period"
	ViewCleanAp            = "clean_ap"
	ViewSirenBlacklist     = "siren_blacklist"
	IntermediateViewDebits = "stg_tmp_debits_simplified"
	ViewDebit              = "clean_debit"
)

// AllTables returns all table names that should be truncated during test cleanup
func AllTables() []string {
	return []string{
		TableImportLogs,
		TableStgApconso,
		TableStgApdemande,
		TableStgCotisation,
		TableStgDebit,
		TableStgDelai,
		TableStgEffectif,
		TableStgEffectifEnt,
		TableStgSirene,
		TableStgSireneUl,
		TableStgFilterImport,
		TableStgProcol,
		TableStgSireneHisto,
	}
}

// AllMaterializedViews returns all materialized view names that should be refreshed during test cleanup
func AllMaterializedViews() []string {
	return []string{
		ViewStgApdemandePeriod,
		ViewSirenBlacklist,
	}
}

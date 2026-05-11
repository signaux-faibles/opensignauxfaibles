# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Open Signaux Faibles (`sfdata`) — a Go-based data pipeline for importing, filtering, and exporting socioeconomic data (URSSAF, activité partielle, SIRENE, effectifs) into PostgreSQL, as part of a French early warning system for company financial distress. Written in Go 1.24, uses `cosiner/flag` for CLI, `pgx/v5` for PostgreSQL, `tern` for migrations.

## Build & Test Commands

```bash
# Build
make build                  # Compiles ./sfdata binary (embeds git commit hash)
make build-prod             # Cross-compile for Linux AMD64, CGO disabled

# Database (required for e2e tests)
make start-postgres         # Starts PostgreSQL 17 container (testuser/testpass on testdb)
make stop-postgres

# Tests
go test ./...               # Unit tests only
go test -tags=e2e ./...     # E2E tests only (requires running postgres + built ./sfdata binary)
./test-all.sh               # All tests (unit + e2e)
./test-all.sh --update      # All tests + update golden file snapshots

# Single test
go test ./lib/engine -run TestTrackerReports

# E2E tests with coverage (as in CI)
go test -v -tags=e2e -coverprofile=go-coverage.out -covermode=atomic -coverpkg=./... ./...
```

Note: there is no linter configured. CI (.github/workflows/ci.yml) runs unit tests with coverage and e2e tests separately.

## Architecture

### Data Pipeline

```
Raw CSV/GZ files → Parser → SirenFilter → DataSink (PostgreSQL / CSV) → Clean Views → Parquet Export
```

Three CLI commands: **import** (batch pipeline), **parseFile** (single file to stdout), **export** (clean views to parquet).

### Core Interfaces (`lib/engine/`)

- **Parser** — transforms raw files into typed Tuples. Implementations in `lib/parsing/*/`.
- **Tuple** — a parsed row with table name and SQL-tagged fields.
- **SirenFilter** — SIREN-based perimeter check (accept/reject rows).
- **DataSink** — write destination (PostgreSQL, CSV, stdout, or composite).
- **SinkFactory** — creates DataSink per parser type.

### Parser Registry (`lib/registry/main.go`)

Central registration of all 13 parser types: `apconso`, `apdemande`, `ccsf`, `cotisation`, `debit`, `delai`, `effectif`, `effectif_ent`, `filter`, `procol`, `sirene`, `sirene_histo`, `sirene_ul`. Parser type constants defined in `lib/engine/parser_types.go`.

### Two-Layer Database

- **`stg_*` tables** — raw imported data (staging), populated by sinks.
- **`clean_*` views** — enriched/cleaned data for consumption, auto-refreshed when dependencies are imported.
- Tables and materialized views are registered in `lib/db/tables.go` (used for test cleanup and view refresh logic).
- Migrations in `lib/db/migrations/` (numbered SQL files) run automatically at import start via `tern`. The `migrations` table tracks the last applied migration.
- **Workaround to run only migrations** (no full import): `./sfdata import --batch 1802 --parsers delai` with a small file.

### Materialized View Refresh (`lib/sinks/postgresSink.go`)

Each parser type maps to specific materialized views via `viewsToRefresh`. After data is written, only the views that depend on the imported parser type are refreshed. This is a switch-based explicit dependency, not automatic.

### Import Flow (`importHandler.go`)

1. Batch preparation — discover files by directory or explicit JSON config.
2. Drop indexes (saved in `tmp_saved_indexes` for recovery).
3. Parse files in parallel → filter by SIREN perimeter → write to sinks.
4. Rebuild indexes, refresh materialized views.
5. Track results in `import_logs` table.

### Filter System (3-stage)

Filtering restricts imported data by SIREN to control volume. `sirene` and `sirene_ul` data bypass filtering (imported in full).

1. **Compute perimeter** (`computePerimeter` command): reads `effectif_ent`, writes the requested perimeter into `stg_filter_import` (companies that reached 10+ employees in the last 120 months), then refreshes `siren_blacklist`.
2. **Import-time SIREN filtering**: `import` reads the filter from `stg_filter_import` (or an explicit `filter_*` file) and rejects rows outside the perimeter as they are parsed.
3. **SQL refinement**: `siren_blacklist` (materialized) excludes from the perimeter any SIREN whose juridic category, NAF activity code, or domiciliation matches a business-driven exclusion list. `clean_filter` = `stg_filter_import - siren_blacklist`. All `clean_*` views downstream filter on `siren_blacklist`.

To evolve filtering logic: `CREATE OR REPLACE VIEW siren_blacklist_logic`, then `REFRESH MATERIALIZED VIEW siren_blacklist`.

Without a filter, import fails by default. Use `--no-filter` to import everything.

### Responsibility split: `computePerimeter` vs `import`

- `computePerimeter` produces the **requested** perimeter (`stg_filter_import`) from `effectif_ent` and refreshes `siren_blacklist`.
- `import` imports all other staging tables, applies the SIREN filter at parse time, and refreshes the `clean_*` views via parser-specific rules in `lib/sinks/postgresSink.go::viewsToRefresh`.

### `clean_filter` requires SIRENE data

`siren_blacklist_logic` (migration `042_change_perimeter.sql`) joins `stg_filter_import` against `stg_sirene_ul` (juridic category, NAF activity) and `stg_sirene` (siège, département). None of these attributes exist in `effectif_ent`, so `clean_filter` is only complete once `stg_sirene` and `stg_sirene_ul` have been populated.

**Required execution order for a from-scratch pipeline**:

```
sfdata import --parsers sirene,sirene_ul …    # populate SIRENE staging tables
sfdata computePerimeter …                      # compute perimeter, refresh siren_blacklist
sfdata import …                                # import remaining files using the filter
```

`computePerimeter` refuses to run if `stg_sirene` or `stg_sirene_ul` is empty.

## Adding a New Parser

1. Define parser type constant in `lib/engine/parser_types.go`.
2. Create `lib/parsing/<name>/parser.go` — implement `Parser` interface (`New()`, `Type()`).
3. Define tuple struct with `csv:"column_name"` and `sql:"column_name"` tags.
4. Register in `lib/registry/main.go`.
5. Add `stg_<name>` table migration in `lib/db/migrations/`, then register table in `lib/db/tables.go`.
6. Add `clean_<name>` view migration. For materialized views, also register in `lib/db/tables.go`.
7. Add filename pattern in `lib/prepare-import/parsertypes.go`.
8. Add sink support in `lib/sinks/postgresSink.go`. For materialized views, add refresh conditions in `viewsToRefresh`/`CreateSink`.

## Key Conventions

- **Batch key format**: YYMM (e.g., `"1802"` for February 2018).
- **CSV parsing**: header-based column indexing via `HeaderIndexer`, non-fatal parse errors collected (capped at 200 per file in `MaxParsingErrors`).
- **Testing**: `stretchr/testify` for assertions. Golden file pattern with `compareWithGoldenFileOrUpdate()`. Mock batches via `MockBatch()` and helpers in `lib/engine/testing.go`.
- **E2E tests**: tagged `//go:build e2e`, run the compiled `./sfdata` binary, compare stdout against `tests/output-snapshots/*.golden.txt`.
- **Compressed files**: `.csv.gz` supported transparently.
- **Performance**: batch insert size 100,000 rows; `work_mem` and `maintenance_work_mem` set to 512MB during import. Indexes dropped before import and rebuilt after.
- **SIREN/SIRET queries**: to search by SIREN in a SIRET-indexed table, use `siret LIKE '123456789%'` to leverage the index.

## Configuration

Config loaded from (in order): env vars → `/etc/opensignauxfaibles/config.toml` → `~/.opensignauxfaibles/config.toml` → `./config.toml`. Key values: `POSTGRES_DB_URL`, `APP_DATA`, `BATCH_CONFIG_FILE`, `export.path`. For env vars, use uppercase with `_` replacing `.` (e.g., `LOG_LEVEL` for `log.level`).

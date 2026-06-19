//go:build e2e

package main

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type procolAtDateRow struct {
	Siren         string
	DateEffet     time.Time
	ActionProcol  string
	StadeProcol   string
	LibelleProcol *string
}

func TestProcolAtDateExclusions(t *testing.T) {
	cleanDB := setupDBTest(t)
	defer cleanDB()

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, suite.PostgresURI)
	require.NoError(t, err)
	defer conn.Close()

	_, err = conn.Exec(ctx, `
		INSERT INTO stg_procol (siret, date_effet, action_procol, stade_procol) VALUES
		-- Entreprise uniquement en sauvegarde, terminée
		('10000000100001', '2023-04-13', 'sauvegarde', 'ouverture'),
		('10000000100001', '2026-01-23', 'sauvegarde', 'fin_procedure'),
		-- Sauvegarde terminée + redressement ouvert
		('52493493200001', '2023-04-13', 'sauvegarde', 'ouverture'),
		('52493493200001', '2024-04-26', 'sauvegarde', 'plan_continuation'),
		('52493493200001', '2025-12-19', 'redressement', 'ouverture'),
		('52493493200001', '2026-01-23', 'sauvegarde', 'fin_procedure'),
		-- Liquidation ouverte + redressement clos par inclusion autre procédure
		('30000000300001', '2020-01-01', 'liquidation', 'ouverture'),
		('30000000300001', '2021-06-01', 'redressement', 'ouverture'),
		('30000000300001', '2022-03-01', 'redressement', 'inclusion_autre_procedure')
	`)
	require.NoError(t, err)

	queryDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	t.Run("sauvegarde terminée sans autre procédure active", func(t *testing.T) {
		rows := queryProcolAtDate(t, conn, queryDate, "100000001")
		assert.Empty(t, rows)
	})

	t.Run("redressement ouvert malgré une sauvegarde récemment terminée", func(t *testing.T) {
		rows := queryProcolAtDate(t, conn, queryDate, "524934932")
		require.Len(t, rows, 1)
		assert.Equal(t, "redressement", rows[0].ActionProcol)
		assert.Equal(t, "ouverture", rows[0].StadeProcol)
		require.NotNil(t, rows[0].LibelleProcol)
		assert.Equal(t, "Redressement judiciaire", *rows[0].LibelleProcol)
	})

	t.Run("liquidation ouverte malgré un redressement clos par inclusion autre procédure", func(t *testing.T) {
		rows := queryProcolAtDate(t, conn, time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC), "300000003")
		require.Len(t, rows, 1)
		assert.Equal(t, "liquidation", rows[0].ActionProcol)
		require.NotNil(t, rows[0].LibelleProcol)
		assert.Equal(t, "Liquidation judiciaire", *rows[0].LibelleProcol)
	})
}

func queryProcolAtDate(t *testing.T, conn *pgxpool.Pool, date time.Time, siren string) []procolAtDateRow {
	t.Helper()

	rows, err := conn.Query(
		context.Background(),
		`SELECT siren, date_effet, action_procol, stade_procol, libelle_procol
		 FROM procol_at_date($1::date)
		 WHERE siren = $2
		 ORDER BY date_effet`,
		date,
		siren,
	)
	require.NoError(t, err)
	defer rows.Close()

	var result []procolAtDateRow
	for rows.Next() {
		var row procolAtDateRow
		err := rows.Scan(&row.Siren, &row.DateEffet, &row.ActionProcol, &row.StadeProcol, &row.LibelleProcol)
		require.NoError(t, err)
		result = append(result, row)
	}
	require.NoError(t, rows.Err())

	return result
}

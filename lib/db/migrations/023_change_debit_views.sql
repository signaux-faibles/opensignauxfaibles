DROP MATERIALIZED VIEW IF EXISTS clean_debit;

-- Compared to previous migration :
--   - extends time frame compared to previous migration for data science
--   - add column "is_latest" that tags last available data
--   - inverse cross join in order to avoid combinatory explosion
CREATE MATERIALIZED VIEW IF NOT EXISTS clean_debit AS
WITH calendar AS (
    SELECT generate_series(
        '2016-01-01'::date,
        date_trunc('month', current_date)::date,
        '1 month'::interval
    ) AS periode
),
debits AS (
  SELECT
    c.periode,
    d.siret,
    d.periode_debut,
    d.periode_fin,
    d.numero_compte,
    d.numero_ecart_negatif,
    d.part_ouvriere,
    d.part_patronale
  FROM calendar c
  CROSS JOIN LATERAL (
    -- for each period, get last available data for all sirets (inside the perimeter)

    -- DISTINCT ON in conjunction with ORDER BY keeps last available data
    SELECT DISTINCT ON (siret, periode_debut, periode_fin, numero_compte, numero_ecart_negatif)
      siret,
      periode_debut,
      periode_fin,
      numero_compte,
      numero_ecart_negatif,
      part_ouvriere,
      part_patronale
    FROM stg_debit
    WHERE date_traitement <= c.periode - interval '1 month' + interval '20 days'
    AND LEFT(siret, 9) IN (SELECT siren FROM clean_filter)
    ORDER BY siret, periode_debut, periode_fin, numero_compte, numero_ecart_negatif, numero_historique_ecart_negatif DESC
  ) d
),

aggregated AS (
    SELECT
        siret,
        periode,
        sum(part_ouvriere) AS part_ouvriere,
        sum(part_patronale) AS part_patronale
    FROM debits
    GROUP BY siret, periode
)
SELECT
    *,
    periode = (SELECT MAX(a2.periode)
               FROM aggregated a2
               WHERE a2.siret = aggregated.siret) AS is_latest
FROM aggregated;

CREATE INDEX IF NOT EXISTS idx_clean_debit_siret ON clean_debit(siret);
CREATE INDEX IF NOT EXISTS idx_clean_debit_periode ON clean_debit(periode);

---- create above / drop below ----

DROP MATERIALIZED VIEW IF EXISTS clean_debit;

-- Rollback : restaure l'ancienne vue
-- (migration 012_create_debit_views.sql)
CREATE MATERIALIZED VIEW IF NOT EXISTS clean_debit AS
WITH calendar AS (
    SELECT date_trunc('month', current_date) - generate_series(1, 24) * '1 month'::interval AS periode
),
debits AS (
    SELECT DISTINCT ON (c.periode, d.siret, d.periode_debut, d.periode_fin, d.numero_compte, d.numero_ecart_negatif)
        c.periode,
        d.siret,
        d.periode_debut,
        d.periode_fin,
        d.numero_compte,
        d.numero_ecart_negatif,
        d.part_ouvriere,
        d.part_patronale
    FROM stg_debit d
    INNER JOIN calendar c ON d.date_traitement <= c.periode - '1 month'::interval + '20 days'::interval
                         AND LEFT(d.siret, 9) IN (SELECT siren FROM clean_filter)
    ORDER BY c.periode, d.siret, d.periode_debut, d.periode_fin, d.numero_compte, d.numero_ecart_negatif, d.numero_historique_ecart_negatif DESC
)
SELECT
    siret,
    periode,
    sum(part_ouvriere) AS part_ouvriere,
    sum(part_patronale) AS part_patronale
FROM debits
GROUP BY siret, periode;

CREATE INDEX IF NOT EXISTS idx_clean_debit_siret ON clean_debit(siret);
CREATE INDEX IF NOT EXISTS idx_clean_debit_periode ON clean_debit(periode);

DROP MATERIALIZED VIEW IF EXISTS clean_debit;

-- Compared to previous migration :
--   - extends time frame compared to previous migration for data science
--   - add column "is_latest" that tags last available data
CREATE MATERIALIZED VIEW clean_debit AS
WITH debits_simplified AS (
  SELECT DISTINCT ON (stg_debit.periode_prise_en_compte, stg_debit.siret, stg_debit.periode_debut, stg_debit.periode_fin, stg_debit.numero_ecart_negatif) stg_debit.siret,
     stg_debit.periode_debut,
     stg_debit.periode_fin,
     stg_debit.numero_ecart_negatif,
     stg_debit.part_ouvriere,
     stg_debit.part_patronale,
     stg_debit.periode_prise_en_compte
  FROM stg_debit
  WHERE NOT (LEFT(stg_debit.siret::text, 9) IN (SELECT siren_blacklist.siren FROM siren_blacklist))
   ORDER BY stg_debit.periode_prise_en_compte, stg_debit.siret, stg_debit.periode_debut, stg_debit.periode_fin, stg_debit.numero_ecart_negatif, stg_debit.numero_historique_ecart_negatif DESC
        ), debit_summed AS (
         SELECT debits_simplified.siret,
            debits_simplified.periode_prise_en_compte,
            sum(debits_simplified.part_ouvriere) AS part_ouvriere,
            sum(debits_simplified.part_patronale) AS part_patronale
           FROM debits_simplified
          GROUP BY debits_simplified.siret, debits_simplified.periode_prise_en_compte
        )
 SELECT debit_summed.siret,
    debit_summed.periode_prise_en_compte AS periode,
    debit_summed.part_ouvriere,
    debit_summed.part_patronale,
    debit_summed.periode_prise_en_compte = max(debit_summed.periode_prise_en_compte) OVER (PARTITION BY debit_summed.siret) AS is_last
   FROM debit_summed
WITH DATA;

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

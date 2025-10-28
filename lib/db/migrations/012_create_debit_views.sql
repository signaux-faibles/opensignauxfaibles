
CREATE MATERIALIZED VIEW clean_debit AS
WITH calendar AS (
    SELECT date_trunc('month', current_date) - generate_series(1, 24) * '1 month'::interval AS periode
),
-- DISTINCT ON keeps the row with the highest numero_historique_ecart_negatif for each group
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

---- create above / drop below ----

DROP MATERIALIZED VIEW clean_debit;

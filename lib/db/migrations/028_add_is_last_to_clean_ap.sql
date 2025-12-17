-- Ajoute la colonne is_last à la vue clean_ap
-- is_last vaut TRUE si c'est la dernière période disponible pour un SIRET donné

DROP MATERIALIZED VIEW IF EXISTS clean_ap;

CREATE MATERIALIZED VIEW clean_ap AS
WITH aggregated AS (
    SELECT
        siret,
        LEFT(siret, 9) as siren,
        periode,
        SUM(ETP_autorise) as ETP_autorise,
        SUM(ETP_consomme) as ETP_consomme,
        STRING_AGG(DISTINCT motif_recours, '; ' ORDER BY motif_recours) as motif_recours
    FROM (SELECT * FROM stg_apdemande_by_period UNION ALL SELECT * FROM stg_apconso_by_period) tmp
    WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = LEFT(tmp.siret, 9))
    GROUP BY siret, periode
)
SELECT
    siret,
    siren,
    periode,
    ETP_autorise,
    ETP_consomme,
    motif_recours,
    (periode = MAX(periode) OVER (PARTITION BY siret)) as is_last
FROM aggregated;

---- create above / drop below ----

DROP MATERIALIZED VIEW IF EXISTS clean_ap;

CREATE MATERIALIZED VIEW clean_ap AS
SELECT
    siret,
    LEFT(siret, 9) as siren,
    periode,
    SUM(ETP_autorise) as ETP_autorise,
    SUM(ETP_consomme) as ETP_consomme,
    STRING_AGG(DISTINCT motif_recours, '; ' ORDER BY motif_recours) as motif_recours
FROM (SELECT * FROM stg_apdemande_by_period UNION ALL SELECT * FROM stg_apconso_by_period) tmp
WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter)
GROUP BY siret, periode;

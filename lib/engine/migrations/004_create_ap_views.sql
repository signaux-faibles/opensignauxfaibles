CREATE MATERIALIZED VIEW IF NOT EXISTS clean_ap AS

  WITH stg_apdemande_by_period AS (
    -- Vue intermédiaire des demandes
    -- * elle décompose (=répète pour chaque période) les intervalles concernés par
    -- la demdande.
    -- * elle calcule en "ETP équivalent" le nombre d'heures demandés.
    SELECT
        d.siret,
        DATE_TRUNC('month', period_series.period)::date as periode,
        CASE
            WHEN d.heures IS NOT NULL THEN d.heures / 151.67 -- moyenne heures mensuelles
            ELSE 0
        END as ETP_autorise,
        0 as ETP_consomme,
        l.label as motif_recours
    FROM stg_apdemande d
    LEFT JOIN labels_motif_recours l ON d.motif_recours = l.id
    CROSS JOIN LATERAL (
        SELECT generate_series(
            DATE_TRUNC('month', d.periode_debut),
            DATE_TRUNC('month', d.periode_fin),
            '1 month'::interval
        ) as period
    ) period_series
    WHERE d.siret IS NOT NULL
      AND d.periode_debut IS NOT NULL
      AND d.periode_fin IS NOT NULL
    ),

    stg_apconso_by_period  AS (
      -- Vue intermédiaire des consommations.
      -- * calcule les ETP consommés,
      -- * formatée dans un format identique que les demandes pour permettre un UNION ALL
        SELECT
            c.siret,
            periode::date,
            0 as ETP_autorise,
            CASE
                WHEN c.heures IS NOT NULL THEN c.heures / 151.67
                ELSE 0
            END as ETP_consomme,
            NULL as motif_recours  -- pas de motif pour la consommation, uniquement pour les demandes
        FROM stg_apconso c
        WHERE c.siret   IS NOT NULL
          AND c.periode IS NOT NULL
       )

    -- on agrége par siret x période les consommations et les demandes
    SELECT
        siret,
        LEFT(siret, 9)::VARCHAR(9) as siren,
        periode,
        SUM(ETP_autorise) as ETP_autorise,
        SUM(ETP_consomme) as ETP_consomme,
        -- Pour motif_recours, on concatène les valeurs uniques
        STRING_AGG(DISTINCT motif_recours, '; ' ORDER BY motif_recours) as motif_recours

    FROM (SELECT * FROM stg_apdemande_by_period UNION ALL SELECT * FROM stg_apconso_by_period) tmp
    GROUP BY siret, periode;

CREATE INDEX IF NOT EXISTS idx_clean_ap_period ON clean_ap(periode);
CREATE INDEX IF NOT EXISTS idx_clean_ap_siret ON clean_ap(siret);
CREATE INDEX IF NOT EXISTS idx_clean_ap_siren ON clean_ap(siren);

---- create above / drop below ----
DROP MATERIALIZED VIEW IF EXISTS clean_ap CASCADE;

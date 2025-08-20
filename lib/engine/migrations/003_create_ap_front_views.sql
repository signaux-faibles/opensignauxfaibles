-- Vue 1: agrégation au siret
CREATE OR REPLACE VIEW frontend_ap_etab AS
SELECT
    siret,
    LEFT(siret, 9) as siren,
    periode,
    SUM(ETP_autorise) as ETP_autorise,
    SUM(ETP_consomme) as ETP_consomme,
    -- Pour motif_recours, on concatène les valeurs uniques
    STRING_AGG(DISTINCT motif_recours, '; ' ORDER BY motif_recours) as motif_recours

FROM (
    -- Subquery qui calcule les ETP autorisés, décomposés (=répétés) sur
    -- toutes les périodes de chaque demande
    SELECT
        d.siret,
        DATE_TRUNC('month', period_series.period) as periode,
        CASE
            WHEN d.heures IS NOT NULL THEN d.heures / 151.67
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
      AND DATE_TRUNC('month', period_series.period) >= CURRENT_DATE - INTERVAL '24 months'

    UNION ALL

    -- Subquery qui calcule les ETP consommés
    SELECT
        c.siret,
        DATE_TRUNC('month', c.periode) as periode,
        0 as ETP_autorise,
        CASE
            WHEN c.heures IS NOT NULL THEN c.heures / 151.67
            ELSE 0
        END as ETP_consomme,
        NULL as motif_recours  -- pas de motif pour la consommation, uniquement pour les demandes
    FROM stg_apconso c
    WHERE c.siret IS NOT NULL
      AND c.periode IS NOT NULL
      AND DATE_TRUNC('month', c.periode) >= CURRENT_DATE - INTERVAL '24 months'
) ETP
GROUP BY siret, periode;

-- Vue 2: agrégation au siren
CREATE OR REPLACE VIEW frontend_ap_entr AS
SELECT
    siren,
    periode,
    SUM(ETP_autorise) as ETP_autorise,
    SUM(ETP_consomme) as ETP_consomme,
    STRING_AGG(DISTINCT motif_recours, '; ' ORDER BY motif_recours) as motif_recours
FROM frontend_ap_etab
GROUP BY siren, periode;

---- create above / drop below ----

DROP VIEW IF EXISTS frontend_ap_entr CASCADE;
DROP VIEW IF EXISTS frontend_ap_etab CASCADE;



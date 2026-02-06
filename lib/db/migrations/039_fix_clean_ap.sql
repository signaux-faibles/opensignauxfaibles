-- is_last is now computed as data for the last available period in
-- `stg_apconso`
--
-- some refactorings along the way of intermediate views (mock columns are only setup at the
-- UNION ALL stage, otherwise it is quite confusing)
BEGIN;
DROP materialized VIEW clean_ap;
DROP VIEW stg_apconso_by_period;

CREATE OR REPLACE VIEW stg_apconso_by_period
AS SELECT c.siret,
    c.periode,
    CASE
      WHEN c.heures IS NOT NULL THEN c.heures / 151.67::double precision
        ELSE 0::double precision
      END AS etp_consomme,
    c.periode = max(c.periode) over () as is_last
   FROM stg_apconso c
  WHERE c.siret IS NOT NULL AND c.periode IS NOT NULL;

drop materialized view stg_apdemande_by_period;

CREATE MATERIALIZED VIEW stg_apdemande_by_period
TABLESPACE pg_default
AS SELECT d.siret,
    date_trunc('month'::text, period_series.period)::date AS periode,
        CASE
            WHEN d.heures IS NOT NULL THEN d.heures / 151.67::double precision
            ELSE 0::double precision
        END AS etp_autorise,
    l.label AS motif_recours
   FROM stg_apdemande d
     LEFT JOIN labels_motif_recours l ON d.motif_recours = l.id
     CROSS JOIN LATERAL ( SELECT generate_series(date_trunc('month'::text, d.periode_debut::timestamp with time zone), date_trunc('month'::text, d.periode_fin::timestamp with time zone), '1 mon'::interval) AS period) period_series
  WHERE d.siret IS NOT NULL AND d.periode_debut IS NOT NULL AND d.periode_fin IS NOT NULL
WITH DATA;

CREATE INDEX idx_stg_apdemande_by_period_period ON stg_apdemande_by_period USING btree (periode);
CREATE INDEX idx_stg_apdemande_by_period_siren ON stg_apdemande_by_period USING btree ("left"((siret)::text, 9));
CREATE INDEX idx_stg_apdemande_by_period_siret ON stg_apdemande_by_period USING btree (siret);


CREATE MATERIALIZED VIEW clean_ap
TABLESPACE pg_default
AS WITH aggregated AS (
         SELECT tmp.siret,
            "left"(tmp.siret::text, 9) AS siren,
            tmp.periode,
            sum(tmp.etp_autorise) AS etp_autorise,
            sum(tmp.etp_consomme) AS etp_consomme,
            string_agg(DISTINCT tmp.motif_recours, '; '::text ORDER BY tmp.motif_recours) AS motif_recours,
            coalesce(bool_or(tmp.is_last), false) as is_last
           FROM ( SELECT stg_apdemande_by_period.siret,
                    stg_apdemande_by_period.periode,
                    stg_apdemande_by_period.etp_autorise,
                    0 AS etp_consomme,
                    stg_apdemande_by_period.motif_recours,
                    null as is_last
                   FROM stg_apdemande_by_period
                UNION ALL
                 SELECT stg_apconso_by_period.siret,
                    stg_apconso_by_period.periode,
                    0 AS etp_autorise,
                    stg_apconso_by_period.etp_consomme,
                    NULL::text AS motif_recours,
                    stg_apconso_by_period.is_last
                   FROM stg_apconso_by_period
                 ) tmp
          WHERE NOT (EXISTS ( SELECT b.siren
                   FROM siren_blacklist b
                  WHERE b.siren::text = "left"(tmp.siret::text, 9)))
          GROUP BY tmp.siret, tmp.periode
        )
 SELECT aggregated.siret,
    aggregated.siren,
    aggregated.periode,
    aggregated.etp_autorise,
    aggregated.etp_consomme,
    aggregated.motif_recours,
    aggregated.is_last
   FROM aggregated
WITH DATA;
COMMIT;

CREATE MATERIALIZED VIEW IF NOT EXISTS stg_tmp_debits_simplified AS
SELECT DISTINCT ON (siret, periode_prise_en_compte, debit_id)
    siret,
    periode_prise_en_compte,
    debit_id,
    part_ouvriere,
    part_patronale
    FROM stg_debit
  ORDER BY siret, debit_id, periode_prise_en_compte, numero_historique_ecart_negatif DESC
WITH NO DATA;

CREATE INDEX IF NOT EXISTS idx_debits_simplified_siret_debit_periode
  ON stg_tmp_debits_simplified(siret, debit_id, periode_prise_en_compte DESC)
  INCLUDE (part_ouvriere, part_patronale);

DROP MATERIALIZED VIEW clean_debit;

CREATE MATERIALIZED VIEW clean_debit AS
  WITH periodes_uniques AS (
      SELECT DISTINCT
        siret,
        periode_prise_en_compte as periode
      FROM stg_tmp_debits_simplified
      WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(p.siret, 9) = b.siren)
  )
  SELECT
    p.siret,
    p.periode as periode,
    SUM(sub.part_ouvriere) as part_ouvriere,
    SUM(sub.part_patronale) as part_patronale
  FROM periodes_uniques p
    CROSS JOIN LATERAL (
     SELECT DISTINCT ON (siret, debit_id)
       d.part_patronale,
       d.part_ouvriere
       FROM stg_tmp_debits_simplified d
       WHERE d.siret= p.siret
         AND d.periode_prise_en_compte <= p.periode
       ORDER BY siret, debit_id, periode_prise_en_compte DESC
  ) sub
  GROUP BY p.siret, p.periode
WITH NO DATA;

CREATE INDEX IF NOT EXISTS idx_clean_debit_siren ON sfdata.clean_debit USING btree ("left"((siret)::text, 9));
CREATE INDEX IF NOT EXISTS idx_clean_debit_period ON sfdata.clean_debit USING btree (periode);
CREATE INDEX IF NOT EXISTS idx_clean_debit_siret ON sfdata.clean_debit USING btree (siret);


---- create above / drop below ----

-- reset to previous view
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

DROP MATERIALIZED VIEW IF EXISTS stg_tmp_debits_simplified;

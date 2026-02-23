-- This migration is more tricky than it seems : because `siren_blacklist` is a
-- materialized view, changing the query means to DROP every view depending on
-- it (and there are many)
--
-- Replacing the materialized view with a normal view is not an option because
-- of performance.
--
-- To avoid this difficulty for further perimeter changes, we introduce a
-- pattern with an intermediate (normal) view (named
-- `siren_blacklist_logic`), that can be changed with a simple `CREATE OR REPLACE VIEW`,
-- on which depends the materialized (cached) view `siren_blacklist`.
SET LOCAL work_mem TO '512MB';
SET LOCAL maintenance_work_mem TO '512MB';

CREATE INDEX IF NOT EXISTS idx_sirene_siren_siege
  ON stg_sirene (siren)
  WHERE siege AND departement <> '';

DROP MATERIALIZED VIEW siren_blacklist CASCADE;

-- this view definition can be replaced with `CREATE OR REPLACE VIEW`
CREATE OR REPLACE VIEW siren_blacklist_logic AS
  WITH excluded_categories AS (
 -- Excluded Catégories Juridiques:
    SELECT ARRAY[
      '4110', -- Établissement public national à caractère industriel ou commercial
      '4120', -- Établissement public national à caractère administratif
      '4140', -- Établissement public local à caractère industriel ou commercial
      '4160', -- Établissement public local à caractère administratif
      '7210', -- Commune et commune nouvelle
      '7220', -- Département
      '7346', -- Association de droit local (Bas-Rhin, Haut-Rhin et Moselle)
      '7348', -- Association intermédiaire
      '7366', -- Syndicat mixte fermé
      '7373', -- Association syndicale libre
      '7379', -- Autre groupement de droit privé non doté de la personnalité morale
      '7383', -- Syndicat mixte ouvert
      '7389', -- Autre groupement de collectivités territoriales
      '7410', -- Établissement public national d'enseignement
      '7430', -- Établissement public local d'enseignement
      '7470', -- Groupement de coopération sanitaire à gestion publique
      '7490'  -- Autre établissement public local d'enseignement
    ] AS categories
  ),
  excluded_activities AS (
    SELECT ARRAY[
      -- 84.XX: Administration publique et défense ; sécurité sociale obligatoire
      '84.11Z', -- Administration publique générale
      '84.12Z', -- Administration publique (tutelle) de la santé, de la formation, de la culture et des services sociaux, autre que sécurité sociale
      '84.13Z', -- Administration publique (tutelle) des activités économiques
      '84.21Z', -- Affaires étrangères
      '84.22Z', -- Défense
      '84.23Z', -- Justice
      '84.24Z', -- Activités d'ordre public et de sécurité
      '84.25Z', -- Services du feu et de secours
      '84.30A', -- Activités générales de sécurité sociale
      '84.30B', -- Gestion des retraites complémentaires
      '84.30C', -- Distribution sociale de revenus

      -- 85.XX: Enseignement (codes spécifiques)
      '85.10Z', -- Enseignement pré-primaire
      '85.20Z', -- Enseignement primaire
      '85.31Z', -- Enseignement secondaire général
      '85.32Z', -- Enseignement secondaire technique ou professionnel
      '85.41Z', -- Enseignement post-secondaire non supérieur
      '85.42Z', -- Enseignement supérieur

      -- 94.XX: Activités des organisations associatives (codes spécifiques)
      '94.11Z', -- Activités des organisations patronales et consulaires
      '94.12Z', -- Activités des organisations professionnelles
      '94.20Z', -- Activités des syndicats de salariés
      '94.91Z', -- Activités des organisations religieuses
      '94.92Z', -- Activités des organisations politiques

      -- 64.XX: Activités des services financiers, hors assurance et caisses de retraite
      '64.11Z', -- Activités de banque centrale
      '64.19Z', -- Autres intermédiations monétaires
      '64.30Z', -- Fonds de placement et entités financières similaires
      '64.91Z', -- Crédit-bail
      '64.92Z', -- Autre distribution de crédit
      '64.99Z', -- Autres activités des services financiers, hors assurance et caisses de retraite, n.c.a.

      -- 65.XX: Assurance
      '65.11Z', -- Assurance vie
      '65.12Z', -- Autres assurances
      '65.20Z', -- Réassurance
      '65.30Z', -- Caisses de retraite

      -- 66.XX: Activités auxiliaires de services financiers et d'assurance
      '66.11Z', -- Administration de marchés financiers
      '66.12Z', -- Courtage de valeurs mobilières et de marchandises
      '66.19A', -- Supports juridiques de gestion de patrimoine mobilier
      '66.19B', -- Autres activités auxiliaires de services financiers, hors assurance et caisses de retraite, n.c.a.
      '66.21Z', -- Évaluation des risques et dommages
      '66.22Z', -- Activités des agents et courtiers d'assurances
      '66.29Z', -- Autres activités auxiliaires d'assurance et de caisses de retraite
      '66.30Z', -- Gestion de fonds

      -- 99.XX: Activités des organisations et organismes extraterritoriaux
      '99.00Z'  -- Activités des organisations et organismes extraterritoriaux
    ] AS activities
  )
 SELECT fp.siren
   FROM stg_filter_import fp
     LEFT JOIN stg_sirene_ul sirene_ul
       ON sirene_ul.siren::text = fp.siren::text
     LEFT JOIN stg_sirene sirene
      ON sirene.siren = fp.siren
        AND sirene.siege = true
        -- exclure les sièges à l'étranger (département vide)
        AND sirene.departement <> ''
     CROSS JOIN excluded_categories ec
     CROSS JOIN excluded_activities ea
  WHERE (sirene_ul.categorie_juridique::text = ANY (ec.categories))
    OR sirene_ul.activite_principale::text = ANY (ea.activities)
    OR sirene.siren IS NULL;

-- This materialized view serves as a cache for siren_blacklist_logic
CREATE MATERIALIZED VIEW IF NOT EXISTS siren_blacklist AS
SELECT * FROM siren_blacklist_logic
WITH DATA;

CREATE UNIQUE INDEX IF NOT EXISTS siren_blacklist_siren_index
  ON siren_blacklist (siren);

-- recreate dependent views...
-- clean_ap
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

-- clean_cotisation
CREATE OR REPLACE VIEW clean_cotisation AS
SELECT
  siret,
  periode_debut as periode,
  sum(du) as du
FROM stg_cotisation c
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(c.siret, 9) = b.siren)
GROUP BY siret, periode_debut;

-- clean_debit
CREATE MATERIALIZED VIEW clean_debit AS
  WITH periodes_uniques AS (
      SELECT DISTINCT
        siret,
        periode_prise_en_compte as periode
      FROM stg_tmp_debits_simplified
      WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(siret, 9) = b.siren)
  ),
  aggregated AS (
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
  )
  SELECT
    siret,
    periode,
    part_ouvriere,
    part_patronale,
    periode = MAX(periode) OVER (PARTITION BY siret) AS is_last
  FROM aggregated
WITH NO DATA;

CREATE INDEX IF NOT EXISTS idx_clean_debit_siren ON clean_debit USING btree ("left"((siret)::text, 9));
CREATE INDEX IF NOT EXISTS idx_clean_debit_period ON clean_debit USING btree (periode);
CREATE INDEX IF NOT EXISTS idx_clean_debit_siret ON clean_debit USING btree (siret);

-- clean_delai
CREATE OR REPLACE VIEW clean_delai
  AS SELECT *
  FROM stg_delai d
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(d.siret, 9) = b.siren);

-- clean_effectif
CREATE OR REPLACE VIEW clean_effectif AS
  SELECT
    *,
    -- Nouvelle colonne qui taggue la dernière valeur disponible pour chaque siret
    periode = (SELECT MAX(e2.periode)
               FROM stg_effectif e2
               WHERE e2.siret = e.siret) AS is_latest
  FROM stg_effectif e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = LEFT(e.siret, 9));

-- clean_effectif_ent
CREATE OR REPLACE VIEW clean_effectif_ent AS
  SELECT
    *,
    -- Nouvelle colonne qui taggue la dernière valeur disponible pour chaque siren
    periode = (SELECT MAX(e2.periode)
               FROM stg_effectif_ent e2
               WHERE e2.siren = e.siren) AS is_latest -- nouvelle colonne
  FROM stg_effectif_ent e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = e.siren);

-- clean_filter
CREATE OR REPLACE VIEW clean_filter AS
  SELECT f.siren
  FROM  stg_filter_import f
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = f.siren);

-- clean_procol
CREATE OR REPLACE VIEW clean_procol AS
SELECT
  LEFT(siret, 9) as siren,
  date_effet,
  action_procol,
  stade_procol,
  CASE
    WHEN stade_procol = 'fin_procedure' THEN 'In bonis'
    WHEN action_procol = 'redressement' AND stade_procol = 'plan_continuation' THEN 'Plan de redressement'
    WHEN action_procol = 'redressement' THEN 'Redressement judiciaire'
    WHEN action_procol = 'liquidation' THEN 'Liquidation judiciaire'
    WHEN action_procol = 'sauvegarde' AND stade_procol = 'plan_continuation' THEN 'Plan de sauvegarde'
    WHEN action_procol = 'sauvegarde' THEN 'Sauvegarde'
  END AS libelle_procol
FROM stg_procol p
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(p.siret, 9) = b.siren)
  -- On ignore la clôture administrative de la procédure
  AND stade_procol != 'solde_procedure'
GROUP BY LEFT(siret, 9), date_effet, action_procol, stade_procol;

-- clean_sirene_histo

CREATE OR REPLACE VIEW clean_sirene_histo AS
WITH ranked_changes AS (
  SELECT *,
         rank() OVER (PARTITION BY siret ORDER BY date_debut ASC) as rank
  FROM stg_sirene_histo sh
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(sh.siret, 9) = b.siren)
)
SELECT
  siret,
  date_debut,
  LEAD(date_debut) OVER (PARTITION BY siret ORDER BY rank) - 1 as date_fin,
  est_actif
FROM ranked_changes
WHERE rank = 1 OR changement_statut_actif;

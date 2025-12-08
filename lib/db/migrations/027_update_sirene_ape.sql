-- Remove code_activite and nomenclature_activite from stg_sirene
-- as we now only accept NAFRev2 nomenclature
ALTER TABLE stg_sirene DROP COLUMN IF EXISTS code_activite;
ALTER TABLE stg_sirene DROP COLUMN IF EXISTS nomenclature_activite;

-- Rename activite_principale to ape in stg_sirene_ul
-- as we now only accept NAFRev2 nomenclature
ALTER TABLE stg_sirene_ul RENAME COLUMN activite_principale TO ape;

-- Update clean_sirene view to use ape instead of code_activite and nomenclature_activite
CREATE OR REPLACE VIEW clean_sirene AS
  SELECT
    siren,
    siret,
    siege,
    -- adresse
    TRIM(
      CONCAT_WS(' ',
        COALESCE(numero_voie, ''),
        COALESCE(indrep, ''),
        COALESCE(type_voie, ''),
        COALESCE(voie, ''),
        CASE WHEN complement_adresse IS NOT NULL THEN CONCAT('(', complement_adresse, ')') END
        ) || E'\n' ||
      CONCAT_WS(' ',
        COALESCE(code_postal, ''),
        COALESCE(commune, commune_etranger, ''),
        CASE WHEN pays_etranger IS NOT NULL THEN CONCAT(E'\n', pays_etranger) END
      )
    ) AS adresse,
    code_commune,
    departement,
    ape,
    date_creation,
    longitude,
    latitude,
    est_actif
  FROM stg_sirene
  WHERE siren IN (SELECT siren FROM clean_filter);

-- Update clean_sirene_ul view to add ape column
CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    siren,
    raison_sociale,
    statut_juridique,
    ape,
    creation,
    est_actif
  FROM stg_sirene_ul
  WHERE siren IN (SELECT siren FROM clean_filter);

-- Update clean_filter to use ape instead of activite_principale
DROP MATERIALIZED VIEW IF EXISTS clean_filter;

CREATE MATERIALIZED VIEW clean_filter AS
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
  )
  SELECT fp.siren
  FROM stg_filter_import fp
  INNER JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren = fp.siren
  CROSS JOIN excluded_categories ec
  WHERE
    -- Exclude if statut_juridique is in the excluded categories list
    NOT (sirene_ul.statut_juridique = ANY(ec.categories))
    -- Exclude Activity Codes:
    -- 84.XX: Administration publique et défense ; sécurité sociale obligatoire
    -- 85.XX: Enseignement
    AND NOT (sirene_ul.ape LIKE '84%' OR sirene_ul.ape LIKE '85%');

CREATE UNIQUE INDEX clean_filter_siren_index
    ON clean_filter(siren);

---- create above / drop below ----

-- Restore code_activite and nomenclature_activite columns
ALTER TABLE stg_sirene ADD COLUMN code_activite VARCHAR(5);
ALTER TABLE stg_sirene ADD COLUMN nomenclature_activite VARCHAR(10);

-- Rename ape back to activite_principale
ALTER TABLE stg_sirene_ul RENAME COLUMN ape TO activite_principale;

-- Revert to previous clean_sirene view with code_activite and nomenclature_activite
CREATE OR REPLACE VIEW clean_sirene AS
  SELECT
    siren,
    siret,
    siege,
    -- adresse
    TRIM(
      CONCAT_WS(' ',
        COALESCE(numero_voie, ''),
        COALESCE(indrep, ''),
        COALESCE(type_voie, ''),
        COALESCE(voie, ''),
        CASE WHEN complement_adresse IS NOT NULL THEN CONCAT('(', complement_adresse, ')') END
        ) || E'\n' ||
      CONCAT_WS(' ',
        COALESCE(code_postal, ''),
        COALESCE(commune, commune_etranger, ''),
        CASE WHEN pays_etranger IS NOT NULL THEN CONCAT(E'\n', pays_etranger) END
      )
    ) AS adresse,
    code_commune,
    departement,
    code_activite,
    nomenclature_activite,
    date_creation,
    longitude,
    latitude,
    est_actif
  FROM stg_sirene
  WHERE siren IN (SELECT siren FROM clean_filter);

-- Revert clean_sirene_ul view
CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    siren,
    raison_sociale,
    statut_juridique,
    activite_principale,
    creation,
    est_actif
  FROM stg_sirene_ul
  WHERE siren IN (SELECT siren FROM clean_filter);

-- Revert clean_filter to use activite_principale
DROP MATERIALIZED VIEW IF EXISTS clean_filter;

CREATE MATERIALIZED VIEW clean_filter AS
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
  )
  SELECT fp.siren
  FROM stg_filter_import fp
  INNER JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren = fp.siren
  CROSS JOIN excluded_categories ec
  WHERE
    -- Exclude if statut_juridique is in the excluded categories list
    NOT (sirene_ul.statut_juridique = ANY(ec.categories))
    -- Exclude Activity Codes:
    -- 84.XX: Administration publique et défense ; sécurité sociale obligatoire
    -- 85.XX: Enseignement
    AND NOT (sirene_ul.activite_principale LIKE '84%' OR sirene_ul.activite_principale LIKE '85%');

CREATE UNIQUE INDEX clean_filter_siren_index
    ON clean_filter(siren);

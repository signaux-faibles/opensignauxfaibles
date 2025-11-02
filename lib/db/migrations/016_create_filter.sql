CREATE TABLE IF NOT EXISTS filter_partial (
    siren VARCHAR(9) PRIMARY KEY
);

CREATE UNIQUE INDEX filter_partial_siren_index
    ON filter_partial(siren);

CREATE MATERIALIZED VIEW IF NOT EXISTS filter AS
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
  FROM filter_partial fp
  INNER JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren = fp.siren
  CROSS JOIN excluded_categories ec
  WHERE
    -- Exclude if statut_juridique is in the excluded categories list
    NOT (sirene_ul.statut_juridique = ANY(ec.categories))
    -- Exclude Activity Codes:
    -- 84.XX: Administration publique et défense ; sécurité sociale obligatoire
    -- 85.XX: Enseignement
    AND NOT (sirene_ul.activite_principale LIKE '84%' OR sirene_ul.activite_principale LIKE '85%');

CREATE UNIQUE INDEX filter_siren_index
    ON filter(siren);

---- create above / drop below ----

DROP TABLE IF EXISTS filter_partial;
DROP MATERIALIZED VIEW IF EXISTS filter;

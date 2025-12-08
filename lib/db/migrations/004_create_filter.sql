-- filter_import est le périmètre de l'import des données.
-- Ce n'est pas le filtrage définitif, qui croise plusieurs données, mais
-- un filtrage sur la seule donnée de l'effectif qui limite déjà
-- considérablement le volume des données importées.
--
-- Un filtrage plus fin sera réalisé via la vue ci-dessous pour la couche de
-- données propres "clean_xxx"
CREATE TABLE IF NOT EXISTS stg_filter_import (
    siren VARCHAR(9) PRIMARY KEY
);

-- clean_filter est le périmètre définitif des données distribuées par la
-- couche de données propres "clean_xxx"
CREATE MATERIALIZED VIEW IF NOT EXISTS clean_filter AS
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

---- create above / drop below ----

DROP TABLE IF EXISTS stg_filter_import;
DROP MATERIALIZED VIEW IF EXISTS clean_filter;

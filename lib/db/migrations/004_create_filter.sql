-- stg_filter_import est le périmètre de l'import des données.
-- Ce n'est pas le filtrage définitif, qui croise plusieurs données, mais
-- un filtrage sur la seule donnée de l'effectif qui limite déjà
-- considérablement le volume des données importées.
--
-- Un filtrage plus fin sera réalisé via la vue ci-dessous pour la couche de
-- données propres "clean_xxx"
CREATE TABLE IF NOT EXISTS stg_filter_import (
    siren VARCHAR(9) PRIMARY KEY
);

-- siren_blacklist est une liste de siren à exclure du périmètre final.
CREATE MATERIALIZED VIEW siren_blacklist
AS WITH excluded_categories AS (
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
     JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren::text = fp.siren::text
     CROSS JOIN excluded_categories ec
  -- Exclude Activity Codes:
  -- 84.XX: Administration publique et défense ; sécurité sociale obligatoire
  -- 85.XX: Enseignement
  WHERE (sirene_ul.statut_juridique::text = ANY (ec.categories))
    OR sirene_ul.activite_principale::text ~~ '84%'::text
    OR sirene_ul.activite_principale::text ~~ '85%'::text;

CREATE UNIQUE INDEX siren_blacklist_siren_index ON siren_blacklist(siren);

-- clean_filter représente le périmètre final de Signaux Faibles
-- (stg_filter_import - siren_blacklist)
CREATE OR REPLACE VIEW clean_filter AS
  SELECT f.siren
  FROM  stg_filter_import f
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = f.siren);


---- create above / drop below ----

DROP VIEW IF EXISTS clean_filter;
DROP MATERIALIZED VIEW IF EXISTS siren_blacklist;
DROP TABLE IF EXISTS stg_filter_import;

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
TABLESPACE pg_default
AS WITH excluded_categories AS (
         SELECT ARRAY['4110'::text, '4120'::text, '4140'::text, '4160'::text, '7210'::text, '7220'::text, '7346'::text, '7348'::text, '7366'::text, '7373'::text, '7379'::text, '7383'::text, '7389'::text, '7410'::text, '7430'::text, '7470'::text, '7490'::text] AS categories
        )
 SELECT fp.siren
   FROM stg_filter_import fp
     JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren::text = fp.siren::text
     CROSS JOIN excluded_categories ec
  WHERE (sirene_ul.statut_juridique::text = ANY (ec.categories)) OR sirene_ul.activite_principale::text ~~ '84%'::text OR sirene_ul.activite_principale::text ~~ '85%'::text;

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

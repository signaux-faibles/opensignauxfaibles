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
     JOIN stg_sirene_ul sirene_ul ON sirene_ul.siren::text = fp.siren::text
     CROSS JOIN excluded_categories ec
     CROSS JOIN excluded_activities ea
  WHERE (sirene_ul.statut_juridique::text = ANY (ec.categories))
    OR sirene_ul.activite_principale::text = ANY (ea.activities)
    -- Exclure les entreprises dont le siège est à l'étranger
    -- Cela se manifeste par un departement vide ('')
    OR NOT EXISTS (
      SELECT 1 FROM stg_sirene s
      WHERE s.siren = fp.siren
        AND s.siege = true
        AND s.departement <> ''
    );

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

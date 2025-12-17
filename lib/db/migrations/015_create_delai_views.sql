-- Cette vue permet de ne pas modifier le code des utilisateurs de données
-- le jour ou des transformations sont nécessaires dans ces données.
-- Par ailleurs, les utilisateurs savent qu'ils doivent uniquement requêter
-- les tables nommées 'clean_...'
CREATE OR REPLACE VIEW clean_delai
  AS SELECT *
  FROM stg_delai d
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(d.siret, 9) = b.siren);

---- create above / drop below ----

DROP VIEW IF EXISTS clean_delai;

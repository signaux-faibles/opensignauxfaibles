-- Cette vue permet de ne pas modifier le code des utilisateurs de données
-- le jour ou des transformations sont nécessaires dans ces données.
-- Par ailleurs, les utilisateurs savent qu'ils doivent uniquement requêter
-- les tables nommées 'clean_...'
CREATE OR REPLACE VIEW sfdata_clean_delai AS SELECT * FROM sfdata_stg_delai;

---- create above / drop below ----

DROP VIEW sfdata_clean_delai;

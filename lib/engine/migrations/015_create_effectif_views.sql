-- Ces vues permettent de ne pas modifier le code des utilisateurs de données
-- le jour ou des transformations sont nécessaires dans ces données.
-- Par ailleurs, les utilisateurs savent qu'ils doivent uniquement requêter
-- les tables nommées 'clean_...'
CREATE OR REPLACE VIEW clean_effectif AS (SELECT * FROM stg_effectif);
CREATE OR REPLACE VIEW clean_effectif_ent AS (SELECT * FROM stg_effectif_ent);

---- create above / drop below ----

DROP VIEW clean_effectif;
DROP VIEW clean_effectif_ent;


-- Migrations remplacées par 022_change_effectif_views.sql
CREATE OR REPLACE VIEW clean_effectif AS
  SELECT *
  FROM stg_effectif
  WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter);

CREATE OR REPLACE VIEW clean_effectif_ent AS
  SELECT *
  FROM stg_effectif_ent
  WHERE siren IN (SELECT siren FROM clean_filter);

---- create above / drop below ----

DROP VIEW IF EXISTS clean_effectif;
DROP VIEW IF EXISTS clean_effectif_ent;

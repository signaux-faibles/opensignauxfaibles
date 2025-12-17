
-- Migrations remplacées par 022_change_effectif_views.sql
CREATE OR REPLACE VIEW clean_effectif AS
  SELECT *
  FROM stg_effectif e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(e.siret, 9) = b.siren);

CREATE OR REPLACE VIEW clean_effectif_ent AS
  SELECT *
  FROM stg_effectif_ent e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE e.siren = b.siren);

---- create above / drop below ----

DROP VIEW IF EXISTS clean_effectif;
DROP VIEW IF EXISTS clean_effectif_ent;

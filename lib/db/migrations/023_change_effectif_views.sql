
CREATE OR REPLACE VIEW clean_effectif AS
  SELECT
    *,
    -- Nouvelle colonne qui taggue la dernière valeur disponible pour chaque siret
    periode = (SELECT MAX(e2.periode)
               FROM stg_effectif e2
               WHERE e2.siret = stg_effectif.siret) AS is_latest
  FROM stg_effectif e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = LEFT(e.siret, 9));

CREATE OR REPLACE VIEW clean_effectif_ent AS
  SELECT
    *,
    -- Nouvelle colonne qui taggue la dernière valeur disponible pour chaque siren
    periode = (SELECT MAX(e2.periode)
               FROM stg_effectif_ent e2
               WHERE e2.siren = stg_effectif_ent.siren) AS is_latest -- nouvelle colonne
  FROM stg_effectif_ent e
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE b.siren = e.siren);

---- create above / drop below ----

-- Rollback : restaure l'ancienne vue
-- (migration 016_create_effectif_views.sql)
CREATE OR REPLACE VIEW clean_effectif AS
  SELECT *
  FROM stg_effectif
  WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter);

CREATE OR REPLACE VIEW clean_effectif_ent AS
  SELECT *
  FROM stg_effectif_ent
  WHERE siren IN (SELECT siren FROM clean_filter);

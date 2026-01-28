DROP VIEW IF EXISTS clean_sirene_histo;

CREATE VIEW clean_sirene_histo AS
WITH ranked_changes AS (
  SELECT *,
         rank() OVER (PARTITION BY siret ORDER BY date_debut ASC) as rank
  FROM stg_sirene_histo sh
  WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(sh.siret, 9) = b.siren)
)
SELECT
  siret,
  date_debut,
  LEAD(date_debut) OVER (PARTITION BY siret ORDER BY rank) - 1 as date_fin,
  est_actif
FROM ranked_changes
WHERE rank = 1 OR changement_statut_actif;

---- create above / drop below ----

-- Rollback
DROP VIEW IF EXISTS clean_sirene_histo;


CREATE VIEW clean_sirene_histo AS
SELECT *
FROM stg_sirene_histo sh
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(sh.siret, 9) = b.siren);

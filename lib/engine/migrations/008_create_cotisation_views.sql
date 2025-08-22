CREATE OR REPLACE VIEW clean_cotisations AS
SELECT
  siret,
  periode_debut as periode,
  sum(du) as du
FROM stg_cotisation
GROUP BY siret, periode_debut;

---- create above / drop below ----

DROP VIEW clean_cotisations;

CREATE OR REPLACE VIEW sfdata_clean_cotisation AS
SELECT
  siret,
  periode_debut as periode,
  sum(du) as du
FROM sfdata_stg_cotisation
GROUP BY siret, periode_debut;

---- create above / drop below ----

DROP VIEW sfdata_clean_cotisation;

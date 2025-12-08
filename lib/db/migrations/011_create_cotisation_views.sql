CREATE OR REPLACE VIEW clean_cotisation AS
SELECT
  siret,
  periode_debut as periode,
  sum(du) as du
FROM stg_cotisation
WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter)
GROUP BY siret, periode_debut;

---- create above / drop below ----

DROP VIEW clean_cotisation;

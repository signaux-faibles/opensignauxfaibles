CREATE OR REPLACE VIEW clean_cotisation AS
SELECT
  siret,
  periode_debut as periode,
  sum(du) as du
FROM stg_cotisation c
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(c.siret, 9) = b.siren)
GROUP BY siret, periode_debut;

---- create above / drop below ----

DROP VIEW IF EXISTS clean_cotisation;

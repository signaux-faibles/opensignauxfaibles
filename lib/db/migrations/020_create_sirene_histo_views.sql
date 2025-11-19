CREATE OR REPLACE VIEW clean_sirene_histo AS
SELECT *
FROM stg_sirene_histo
WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter);

---- create above / drop below ----

DROP VIEW clean_sirene_histo;

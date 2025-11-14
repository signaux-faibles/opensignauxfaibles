CREATE OR REPLACE VIEW clean_procol AS
SELECT *
FROM stg_procol
WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter);

---- create above / drop below ----

DROP VIEW clean_procol;

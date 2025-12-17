CREATE OR REPLACE VIEW clean_sirene_histo AS
SELECT *
FROM stg_sirene_histo sh
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(sh.siret, 9) = b.siren);

---- create above / drop below ----

DROP VIEW IF EXISTS clean_sirene_histo;

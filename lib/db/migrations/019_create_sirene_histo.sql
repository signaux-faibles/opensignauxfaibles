CREATE TABLE IF NOT EXISTS stg_sirene_histo (
);

CREATE INDEX IF NOT EXISTS idx_stg_sirene_histo_siret ON stg_sirene_histo(siret);
CREATE INDEX IF NOT EXISTS idx_stg_sirene_histo_siren ON stg_sirene_histo(LEFT(siret, 9));

---- create above / drop below ----

DROP TABLE stg_sirene_histo;

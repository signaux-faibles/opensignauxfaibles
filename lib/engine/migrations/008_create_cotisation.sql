CREATE TABLE IF NOT EXISTS sfdata_stg_cotisation (
    siret           VARCHAR(14),
    periode_debut   DATE,
    periode_fin     DATE,
    du              FLOAT
);

CREATE INDEX IF NOT EXISTS idx_stg_cotisation_siret ON sfdata_stg_cotisation(siret);
CREATE INDEX IF NOT EXISTS idx_stg_cotisation_siren ON sfdata_stg_cotisation(LEFT(siret, 9));
CREATE INDEX IF NOT EXISTS idx_stg_cotisation_periode ON sfdata_stg_cotisation(periode_debut);

---- create above / drop below ----

DROP TABLE sfdata_stg_cotisation;

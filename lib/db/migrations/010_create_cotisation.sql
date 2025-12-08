CREATE TABLE IF NOT EXISTS stg_cotisation (
    siret           VARCHAR(14),
    periode_debut   DATE,
    periode_fin     DATE,
    du              FLOAT
);

CREATE INDEX IF NOT EXISTS idx_stg_cotisation_siret ON stg_cotisation(siret);
CREATE INDEX IF NOT EXISTS idx_stg_cotisation_siren ON stg_cotisation(LEFT(siret, 9));
CREATE INDEX IF NOT EXISTS idx_stg_cotisation_periode ON stg_cotisation(periode_debut);

---- create above / drop below ----

DROP TABLE stg_cotisation;

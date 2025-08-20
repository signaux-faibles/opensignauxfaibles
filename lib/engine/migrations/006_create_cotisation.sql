CREATE TABLE IF NOT EXISTS stg_cotisation (
    siret           VARCHAR(14),
    numero_compte   TEXT,
    periode_debut   DATE,
    periode_fin     DATE,
    du              FLOAT
);

CREATE INDEX IF NOT EXISTS idx_stg_cotisation_key ON stg_cotisation(siret);

---- create above / drop below ----

DROP TABLE stg_cotisation;

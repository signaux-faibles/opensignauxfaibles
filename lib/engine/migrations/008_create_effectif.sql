CREATE TABLE IF NOT EXISTS stg_effectif (
    siret           VARCHAR(14),
    numero_compte   TEXT,
    periode         DATE,
    effectif        INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_siret ON stg_effectif(siret);

---- create above / drop below ----

DROP TABLE stg_effectif;

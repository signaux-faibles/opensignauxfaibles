CREATE TABLE IF NOT EXISTS stg_effectif_ent (
    siren           VARCHAR(9),
    periode         DATE,
    effectif_ent    INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_ent_siren ON stg_effectif_ent(siren);

---- create above / drop below ----

DROP TABLE stg_effectif_ent;

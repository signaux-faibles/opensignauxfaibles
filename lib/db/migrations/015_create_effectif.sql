CREATE TABLE IF NOT EXISTS stg_effectif (
    siret           VARCHAR(14),
    periode         DATE,
    effectif        INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_siret ON stg_effectif(siret);
CREATE INDEX IF NOT EXISTS idx_stg_effectif_siren ON stg_effectif(LEFT(siret, 9));
CREATE INDEX IF NOT EXISTS idx_stg_effectif_periode ON stg_effectif (periode) WHERE (periode >= '2023-01-01':: date);


CREATE TABLE IF NOT EXISTS stg_effectif_ent (
    siren           VARCHAR(9),
    periode         DATE,
    effectif_ent    INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_ent_siren ON stg_effectif_ent(siren);
CREATE INDEX IF NOT EXISTS idx_stg_effectif_ent_periode ON stg_effectif_ent (periode) WHERE (periode >= '2023-01-01':: date);

---- create above / drop below ----


DROP TABLE stg_effectif;
DROP TABLE stg_effectif_ent;

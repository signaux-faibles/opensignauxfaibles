CREATE TABLE IF NOT EXISTS sfdata_stg_effectif (
    siret           VARCHAR(14),
    periode         DATE,
    effectif        INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_siret ON sfdata_stg_effectif(siret);


CREATE TABLE IF NOT EXISTS sfdata_stg_effectif_ent (
    siren           VARCHAR(9),
    periode         DATE,
    effectif_ent    INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_effectif_ent_siren ON sfdata_stg_effectif_ent(siren);

---- create above / drop below ----

DROP TABLE sfdata_stg_effectif;
DROP TABLE sfdata_stg_effectif_ent;

CREATE TABLE IF NOT EXISTS stg_procol (
    siret         VARCHAR(14),
    date_effet    DATE,
    action_procol TEXT,
    stade_procol  TEXT
);

CREATE INDEX IF NOT EXISTS idx_stg_procol_siret ON stg_procol(siret);
CREATE INDEX IF NOT EXISTS idx_stg_procol_siren ON stg_procol(LEFT(siret, 9));

---- create above / drop below ----

DROP TABLE IF EXISTS stg_procol;

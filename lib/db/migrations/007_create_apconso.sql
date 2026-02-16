CREATE TABLE  IF NOT EXISTS stg_apconso  (
    siret               VARCHAR(14),
    id_demande          VARCHAR(11),
    heures              FLOAT,
    montant             FLOAT,
    effectif            INTEGER,
    periode             DATE,
    PRIMARY KEY (siret, periode, id_demande)
);

CREATE INDEX IF NOT EXISTS idx_stg_apconso_id_demande ON stg_apconso(id_demande);
CREATE INDEX IF NOT EXISTS idx_stg_apconso_siret ON stg_apconso(siret);
CREATE INDEX IF NOT EXISTS idx_stg_apconso_siren ON stg_apconso(LEFT(siret, 9));
CREATE INDEX IF NOT EXISTS idx_stg_apconso_period ON stg_apconso(periode);

---- create above / drop below ----

DROP TABLE IF EXISTS stg_apconso;

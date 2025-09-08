CREATE TABLE  IF NOT EXISTS stg_apconso  (
    siret               VARCHAR(14),
    id_demande          VARCHAR(11),
    heures_consommees   FLOAT,
    montant             FLOAT,
    effectif            INTEGER,
    periode             DATE
);

CREATE INDEX IF NOT EXISTS idx_stg_apconso_id_demande ON stg_apconso(id_demande);

---- create above / drop below ----

DROP TABLE stg_apconso;



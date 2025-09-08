CREATE TABLE IF NOT EXISTS sfdata_stg_sirene (
    siren                   VARCHAR(9),
    siret                   VARCHAR(14),
    siege                   BOOLEAN,
    complement_adresse      TEXT,
    numero_voie             VARCHAR(10),
    indrep                  VARCHAR(10),
    type_voie               VARCHAR(20),
    voie                    TEXT,
    commune                 TEXT,
    commune_etranger        TEXT,
    distribution_speciale   TEXT,
    code_commune            VARCHAR(5),
    code_cedex              VARCHAR(5),
    cedex                   VARCHAR(100),
    code_pays_etranger      VARCHAR(10),
    pays_etranger           VARCHAR(100),
    code_postal             VARCHAR(10),
    departement             VARCHAR(10),
    ape                     VARCHAR(100),
    code_activite           VARCHAR(5),
    nomenclature_activite   VARCHAR(10),
    date_creation           DATE,
    longitude               FLOAT,
    latitude                FLOAT
);

CREATE INDEX IF NOT EXISTS idx_stg_sirene_siren ON sfdata_stg_sirene(siren);
CREATE INDEX IF NOT EXISTS idx_stg_sirene_siret ON sfdata_stg_sirene(siret);

---- create above / drop below ----

DROP TABLE sfdata_stg_sirene;

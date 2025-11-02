CREATE TABLE IF NOT EXISTS stg_sirene_ul (
    siren                   VARCHAR(9),
    raison_sociale          TEXT,
    prenom1_unite_legale    VARCHAR(100),
    prenom2_unite_legale    VARCHAR(100),
    prenom3_unite_legale    VARCHAR(100),
    prenom4_unite_legale    VARCHAR(100),
    nom_unite_legale        VARCHAR(100),
    nom_usage_unite_legale  VARCHAR(100),
    statut_juridique        VARCHAR(10),
    activite_principale     VARCHAR(10),
    creation                DATE
);

CREATE INDEX IF NOT EXISTS idx_stg_sirene_ul_siren ON stg_sirene_ul(siren);

---- create above / drop below ----

DROP TABLE stg_sirene_ul;

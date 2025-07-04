CREATE TABLE stg_apdemande (
    ID VARCHAR(255),
    Siret VARCHAR(14) PRIMARY KEY,
    EffectifEntreprise INTEGER,
    Effectif INTEGER,
    DateStatut DATE,
    Periode DATE,
    HTA FLOAT,
    MTA FLOAT,
    EffectifAutorise INTEGER,
    MotifRecoursSE INTEGER,
    HeureConsommee FLOAT,
    MontantConsomme FLOAT,
    EffectifConsomme INTEGER,
    Perimetre INTEGER
);

CREATE TABLE IF NOT EXISTS stg_apdemande (
    ID VARCHAR(255),
    Siret VARCHAR(14) PRIMARY KEY,
    EffectifEntreprise INTEGER,
    Effectif INTEGER,
    DateStatut DATE,
    PeriodeStart DATE,
    PeriodeEnd DATE,
    HTA FLOAT,
    MTA FLOAT,
    EffectifAutorise INTEGER,
    MotifRecoursSE INTEGER,
    HeureConsommee FLOAT,
    MontantConsomme FLOAT,
    EffectifConsomme INTEGER,
    Perimetre INTEGER
);

CREATE TABLE  IF NOT EXISTS stg_apconso  (
    Siret VARCHAR(14) PRIMARY KEY,
    ID VARCHAR(255),
    HeureConsommee FLOAT,
    Montant FLOAT,
    Effectif INTEGER,
    Periode DATE
);

---- create above / drop below ----

DROP TABLE stg_apconso;



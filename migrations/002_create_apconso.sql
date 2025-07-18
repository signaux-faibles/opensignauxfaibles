CREATE TABLE  IF NOT EXISTS stg_apconso  (
    siret VARCHAR(14) PRIMARY KEY,
    id_conso VARCHAR(255),
    heures_consommees FLOAT,
    montant FLOAT,
    effectif INTEGER,
    periode DATE
);

---- create above / drop below ----

DROP TABLE stg_apconso;



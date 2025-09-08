CREATE TABLE IF NOT EXISTS sfdata_stg_delai (
    siret               VARCHAR(14),
    date_creation       DATE,
    date_echeance       DATE,
    duree_delai         INTEGER,
    montant_echeancier  FLOAT,
    stade               VARCHAR(50),
    action              VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_stg_delai_siret ON sfdata_stg_delai(siret);
CREATE INDEX IF NOT EXISTS idx_stg_delai_siret ON sfdata_stg_delai(LEFT(siret, 9));

---- create above / drop below ----

DROP TABLE sfdata_stg_delai;

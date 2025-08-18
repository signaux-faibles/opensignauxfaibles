CREATE TABLE IF NOT EXISTS stg_delai (
    key                 VARCHAR(14),
    numero_compte       TEXT,
    numero_contentieux  VARCHAR(50),
    date_creation       DATE,
    date_echeance       DATE,
    duree_delai         INTEGER,
    denomination        TEXT,
    indic_6m            VARCHAR(10),
    annee_creation      INTEGER,
    montant_echeancier  FLOAT,
    stade               VARCHAR(50),
    action              VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_stg_delai_key ON stg_delai(key);
CREATE INDEX IF NOT EXISTS idx_stg_delai_date_creation ON stg_delai(date_creation);
CREATE INDEX IF NOT EXISTS idx_stg_delai_date_echeance ON stg_delai(date_echeance);

---- create above / drop below ----

DROP TABLE stg_delai;

CREATE TABLE IF NOT EXISTS stg_apdemande (
            id_demande           VARCHAR(11) PRIMARY KEY,
            siret                VARCHAR(14),
            effectif_entreprise  INTEGER,
            effectif             INTEGER,
            date_statut          DATE,
            periode_debut        DATE,
            periode_fin          DATE,
            hta                  FLOAT,
            mta                  FLOAT,
            effectif_autorise    INTEGER,
            motif_recours_se     INTEGER,
            heures_consommees    FLOAT,
            montant_consomme     FLOAT,
            effectif_consomme    INTEGER,
            perimetre            INTEGER
);

CREATE INDEX IF NOT EXISTS idx_stg_apdemande_siret ON stg_apdemande(siret);

---- create above / drop below ----

DROP TABLE stg_apdemande;

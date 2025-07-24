CREATE TABLE IF NOT EXISTS stg_apdemande (
            id_demande           VARCHAR(255),
            siret                VARCHAR(14) PRIMARY KEY,
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

---- create above / drop below ----

DROP TABLE stg_apdemande;

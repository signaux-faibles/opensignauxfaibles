CREATE TABLE IF NOT EXISTS stg_apdemande (
            id                   VARCHAR(255),
            siret                VARCHAR(14)     PRIMARY   KEY,
            effectif_entreprise   INTEGER,
            effectif             INTEGER,
            date_statut           DATE,
            periode_start         DATE,
            periode_end           DATE,
            hta                  FLOAT,
            mta                  FLOAT,
            effectif_autorise     INTEGER,
            motif_recours_se       INTEGER,
            heureconsommee       FLOAT,
            montantconsomme      FLOAT,
            effectifconsomme     INTEGER,
            perimetre            INTEGER
);

---- create above / drop below ----

DROP TABLE stg_apdemande;

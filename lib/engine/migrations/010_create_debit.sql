CREATE TABLE IF NOT EXISTS sfdata_stg_debit (
    siret                           VARCHAR(14),
    numero_compte                   TEXT,
    numero_ecart_negatif            VARCHAR(50),
    date_traitement                 DATE,
    part_ouvriere                   FLOAT,
    part_patronale                  FLOAT,
    numero_historique_ecart_negatif INTEGER,
    etat_compte                     INTEGER,
    code_procedure_collective       VARCHAR(10),
    periode_debut                   DATE,
    periode_fin                     DATE,
    code_operation_ecart_negatif    VARCHAR(10),
    code_motif_ecart_negatif        VARCHAR(10),
    recours_en_cours                BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_stg_debit_siret ON sfdata_stg_debit(siret);
CREATE INDEX IF NOT EXISTS idx_stg_debit_date_traitement ON sfdata_stg_debit(date_traitement);
CREATE INDEX IF NOT EXISTS idx_stg_debit_periode_debut ON sfdata_stg_debit(periode_debut);

---- create above / drop below ----

DROP TABLE sfdata_stg_debit;

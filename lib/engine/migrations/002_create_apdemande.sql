CREATE TABLE IF NOT EXISTS sfdata_labels_motif_recours (
  id    INTEGER PRIMARY KEY,
  label TEXT
);

INSERT INTO sfdata_labels_motif_recours (id, label) VALUES
(1, 'Conjoncture économique'),
(2, 'Difficultés d''approvisionnement en matières premières ou en énergie'),
(3, 'Sinistre ou intempéries de caractère exceptionnel'),
(4, 'Transformation, restructuration ou modernisation des installations et des bâtiments'),
(5, 'Autres circonstances exceptionnelles');

CREATE TABLE IF NOT EXISTS sfdata_stg_apdemande (
            id_demande           VARCHAR(11) PRIMARY KEY,
            siret                VARCHAR(14),
            date_statut          DATE,
            periode_debut        DATE,
            periode_fin          DATE,
            heures               FLOAT,
            montant              FLOAT,
            effectif             INTEGER,
            motif_recours        INTEGER REFERENCES sfdata_labels_motif_recours(id)
);

CREATE INDEX IF NOT EXISTS idx_stg_apdemande_siret ON sfdata_stg_apdemande(siret);
CREATE INDEX IF NOT EXISTS idx_stg_apdemande_siren ON sfdata_stg_apdemande(LEFT(siret, 9));

---- create above / drop below ----

DROP TABLE sfdata_labels_motif_recours;
DROP TABLE sfdata_stg_apdemande;

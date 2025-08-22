CREATE TABLE IF NOT EXISTS labels_motif_recours (
  id    INTEGER PRIMARY KEY,
  label TEXT
);

INSERT INTO labels_motif_recours (id, label) VALUES
(1, 'Conjoncture économique'),
(2, 'Difficultés d''approvisionnement en matières premières ou en énergie'),
(3, 'Sinistre ou intempéries de caractère exceptionnel'),
(4, 'Transformation, restructuration ou modernisation des installations et des bâtiments'),
(5, 'Autres circonstances exceptionnelles');

CREATE TABLE IF NOT EXISTS stg_apdemande (
            id_demande           VARCHAR(11) PRIMARY KEY,
            siret                VARCHAR(14),
            date_statut          DATE,
            periode_debut        DATE,
            periode_fin          DATE,
            heures               FLOAT,
            montant              FLOAT,
            effectif             INTEGER,
            motif_recours        INTEGER REFERENCES labels_motif_recours(id)
);

CREATE INDEX IF NOT EXISTS idx_stg_apdemande_siret ON stg_apdemande(siret);

---- create above / drop below ----

DROP TABLE labels_motif_recours;
DROP TABLE stg_apdemande;

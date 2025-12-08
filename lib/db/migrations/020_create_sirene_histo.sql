CREATE TABLE IF NOT EXISTS stg_sirene_histo (
  siret                   VARCHAR(14),
  date_debut              DATE,
  date_fin                DATE,
  est_actif               BOOLEAN,
  changement_statut_actif BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_stg_sirene_histo_siret ON stg_sirene_histo(siret);
CREATE INDEX IF NOT EXISTS idx_stg_sirene_histo_siren ON stg_sirene_histo(LEFT(siret, 9));

---- create above / drop below ----

DROP TABLE stg_sirene_histo;

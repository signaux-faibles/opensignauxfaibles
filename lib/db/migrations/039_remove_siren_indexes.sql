-- Allow for `siret ~ '^123456789'` queries using the index
-- https://www.postgresql.org/docs/current/indexes-types.html#INDEXES-TYPES-BTREE

DROP INDEX IF EXISTS idx_stg_cotisation_siret;
DROP INDEX IF EXISTS idx_stg_debit_siret;
DROP INDEX IF EXISTS idx_stg_delai_siret;
DROP INDEX IF EXISTS idx_stg_effectif_siret;
DROP INDEX IF EXISTS idx_stg_sirene_siret;
CREATE INDEX idx_stg_cotisation_siret ON stg_cotisation(siret varchar_pattern_ops);
CREATE INDEX idx_stg_debit_siret ON stg_debit(siret varchar_pattern_ops);
CREATE INDEX idx_stg_delai_siret ON stg_delai(siret varchar_pattern_ops);
CREATE INDEX idx_stg_effectif_siret ON stg_effectif(siret varchar_pattern_ops);
CREATE INDEX idx_stg_sirene_siret ON stg_sirene(siret varchar_pattern_ops);

-- DROP obsolete SIREN indexes
DROP INDEX IF EXISTS idx_stg_cotisation_siren;
DROP INDEX IF EXISTS idx_stg_debit_siren;
DROP INDEX IF EXISTS idx_stg_delai_siren;
DROP INDEX IF EXISTS idx_stg_effectif_siren;
DROP INDEX IF EXISTS idx_stg_sirene_siren;

---- create above / drop below ----

-- DROP new SIRET indexes with pattern_ops
DROP INDEX IF EXISTS idx_stg_cotisation_siret;
DROP INDEX IF EXISTS idx_stg_debit_siret;
DROP INDEX IF EXISTS idx_stg_delai_siret;
DROP INDEX IF EXISTS idx_stg_effectif_siret;
DROP INDEX IF EXISTS idx_stg_sirene_siret;

-- Recreate old SIRET indexes without pattern_ops
CREATE INDEX idx_stg_cotisation_siret ON stg_cotisation(siret);
CREATE INDEX idx_stg_debit_siret ON stg_debit(siret);
CREATE INDEX idx_stg_delai_siret ON stg_delai(siret);
CREATE INDEX idx_stg_effectif_siret ON stg_effectif(siret);
CREATE INDEX idx_stg_sirene_siret ON stg_sirene(siret);

-- Recreate old SIREN indexes
CREATE INDEX idx_stg_cotisation_siren ON stg_cotisation(LEFT(siret, 9));
CREATE INDEX idx_stg_debit_siren ON stg_debit(LEFT(siret, 9));
CREATE INDEX idx_stg_delai_siren ON stg_delai(LEFT(siret, 9));
CREATE INDEX idx_stg_effectif_siren ON stg_effectif(LEFT(siret, 9));
CREATE INDEX idx_stg_sirene_siren ON stg_sirene(siren);

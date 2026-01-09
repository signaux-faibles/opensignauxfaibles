ALTER TABLE stg_debit ADD COLUMN IF NOT EXISTS debit_id VARCHAR(33);

DROP INDEX IF EXISTS idx_stg_debit_ecart_negatif;


---- create above / drop below ----

ALTER TABLE stg_debit DROP COLUMN IF EXISTS debit_id;

CREATE INDEX IF NOT EXISTS idx_stg_debit_ecart_negatif ON stg_debit(siret, periode_debut, periode_fin, numero_ecart_negatif, periode_prise_en_compte, numero_historique_ecart_negatif);

ALTER TABLE stg_sirene ADD COLUMN est_actif BOOLEAN;
ALTER TABLE stg_sirene_ul ADD COLUMN est_actif BOOLEAN;

---- create above / drop below ----

ALTER TABLE stg_sirene DROP COLUMN est_actif;
ALTER TABLE stg_sirene_ul DROP COLUMN est_actif;

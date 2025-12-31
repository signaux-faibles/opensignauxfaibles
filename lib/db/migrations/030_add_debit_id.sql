ALTER TABLE stg_debit ADD COLUMN IF NOT EXISTS debit_id VARCHAR(33);

---- create above / drop below ----

ALTER TABLE stg_debit DROP COLUMN IF EXISTS debit_id;

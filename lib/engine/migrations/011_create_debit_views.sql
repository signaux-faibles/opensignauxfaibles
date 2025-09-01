
CREATE OR REPLACE VIEW sfdata_clean_debit AS SELECT * FROM sfdata_stg_debit;

---- create above / drop below ----

DROP VIEW sfdata_clean_debit;

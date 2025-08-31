
CREATE OR REPLACE VIEW clean_debit AS SELECT * FROM stg_debit;

---- create above / drop below ----

DROP VIEW clean_debit;

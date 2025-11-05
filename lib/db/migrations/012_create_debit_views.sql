
CREATE OR REPLACE VIEW clean_debit AS
  SELECT *
  FROM stg_debit
  WHERE LEFT(siret, 9) IN (SELECT siren FROM clean_filter);

---- create above / drop below ----

DROP VIEW clean_debit;

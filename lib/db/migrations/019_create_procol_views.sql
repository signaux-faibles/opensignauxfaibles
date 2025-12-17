CREATE OR REPLACE VIEW clean_procol AS
SELECT
  LEFT(siret, 9) as siren,
  date_effet,
  action_procol,
  stade_procol
FROM stg_procol p
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(p.siret, 9) = b.siren)
  -- On ignore la clôture administrative de la procédure
  AND stade_procol != 'solde_procedure'
GROUP BY LEFT(siret, 9), date_effet, action_procol, stade_procol;

---- create above / drop below ----

DROP VIEW IF EXISTS clean_procol;

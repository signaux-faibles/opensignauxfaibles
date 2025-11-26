--
CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siret VARCHAR(14), date_effet DATE, action_procol TEXT, stade_procol TEXT) AS $$
  WITH last_action_procol AS (
    SELECT DISTINCT ON (siret, action_procol)
      siret, date_effet, action_procol, stade_procol
    FROM clean_procol
    WHERE date_effet <= date_param
      -- On ignore le stade "solde_procedure" qui est une régularisation
      -- administrative et comptable
      AND stade_procol != 'solde_procedure'
    ORDER BY siret, action_procol, date_effet DESC
  )
  SELECT siret, date_effet, action_procol, stade_procol
  FROM last_action_procol
  WHERE action_procol != 'fin_procedure';
$$ LANGUAGE SQL;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns établissements that have a procédure collective in progress on a given date. A single établissement may have several simultaneous proceedings. Completed proceedings are not counted (action_procol = "fin_procedure") — closed établissements are nevertheless displayed.';

---- create above / drop below ----

DROP FUNCTION IF EXISTS procol_at_date;

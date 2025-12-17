--
CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT) AS $$
  WITH last_action_procol AS (
    SELECT DISTINCT ON (siren, action_procol)
      siren, date_effet, action_procol, stade_procol
    FROM clean_procol
    WHERE date_effet <= date_param
    ORDER BY siren, action_procol, date_effet DESC
  )
  SELECT siren, date_effet, action_procol, stade_procol
  FROM last_action_procol
  -- On retire les procédures collectives qui se sont terminées
  WHERE action_procol != 'fin_procedure' AND action_procol != 'inclusion_autre_procedure';
$$ LANGUAGE SQL;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns entreprises that have a procédure collective in progress on a given date. A single entreprise may have several simultaneous proceedings. Completed proceedings are not counted (action_procol = "fin_procedure" or action_procol = "inclusion_autre_procedure") — closed entreprises are nevertheless displayed.';

---- create above / drop below ----

DROP FUNCTION IF EXISTS procol_at_date;

-- Exclude terminated proceedings by stade_procol, matching the URSSAF parser and
-- clean_procol libelle logic (fin_procedure / inclusion_autre_procedure are stades,
-- not action_procol values).
-- CREATE OR REPLACE only: clean_procol_at_date depends on this function.
CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT, libelle_procol TEXT)
LANGUAGE plpgsql
SET search_path FROM CURRENT
AS $$
BEGIN
  RETURN QUERY
  WITH last_action_procol AS (
    SELECT DISTINCT ON (cp.siren, cp.action_procol)
      cp.siren::VARCHAR(9), cp.date_effet, cp.action_procol, cp.stade_procol, cp.libelle_procol
    FROM clean_procol cp
    WHERE cp.date_effet <= date_param
    ORDER BY cp.siren, cp.action_procol, cp.date_effet DESC
  )
  SELECT DISTINCT ON (lap.siren)
    lap.siren, lap.date_effet, lap.action_procol, lap.stade_procol, lap.libelle_procol
  FROM last_action_procol lap
  WHERE lap.stade_procol NOT IN ('fin_procedure', 'inclusion_autre_procedure')
    AND NOT (lap.stade_procol = 'plan_continuation' AND AGE(date_param, lap.date_effet) >= INTERVAL '10 years')
  ORDER BY lap.siren,
    CASE lap.action_procol
      WHEN 'liquidation' THEN 1
      WHEN 'redressement' THEN 2
      WHEN 'sauvegarde' THEN 3
      ELSE 4
    END;
END;
$$;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns one row per entreprise with an active procédure collective on a given date, choosing the most severe proceeding in case of conflict (liquidation > redressement > sauvegarde). Completed proceedings are excluded when stade_procol is fin_procedure or inclusion_autre_procedure.';

---- create above / drop below ----

CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT, libelle_procol TEXT)
LANGUAGE plpgsql
SET search_path FROM CURRENT
AS $$
BEGIN
  RETURN QUERY
  WITH last_action_procol AS (
    SELECT DISTINCT ON (cp.siren, cp.action_procol)
      cp.siren::VARCHAR(9), cp.date_effet, cp.action_procol, cp.stade_procol, cp.libelle_procol
    FROM clean_procol cp
    WHERE cp.date_effet <= date_param
    ORDER BY cp.siren, cp.action_procol, cp.date_effet DESC
  )
  SELECT DISTINCT ON (lap.siren)
    lap.siren, lap.date_effet, lap.action_procol, lap.stade_procol, lap.libelle_procol
  FROM last_action_procol lap
  WHERE lap.action_procol != 'fin_procedure'
    AND lap.action_procol != 'inclusion_autre_procedure'
    AND NOT (lap.stade_procol = 'plan_continuation' AND AGE(date_param, lap.date_effet) >= INTERVAL '10 years')
  ORDER BY lap.siren,
    CASE lap.action_procol
      WHEN 'liquidation' THEN 1
      WHEN 'redressement' THEN 2
      WHEN 'sauvegarde' THEN 3
      ELSE 4
    END;
END;
$$;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns one row per entreprise with an active procédure collective on a given date, choosing the most severe proceeding in case of conflict (liquidation > redressement > sauvegarde). Completed proceedings are excluded (action_procol = "fin_procedure" or "inclusion_autre_procedure").';

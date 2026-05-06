-- Update procol_at_date to return a single row per siren (most severe active procol),
-- and add libelle_procol to the return type.
-- Uses plpgsql to avoid body validation at creation time (clean_procol may have been
-- dropped and recreated by migration 042's CASCADE).
DROP FUNCTION IF EXISTS procol_at_date;

CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT, libelle_procol TEXT) AS $$
BEGIN
  RETURN QUERY
  WITH last_action_procol AS (
    SELECT DISTINCT ON (cp.siren, cp.action_procol)
      cp.siren::VARCHAR(9), cp.date_effet, cp.action_procol, cp.stade_procol, cp.libelle_procol
    FROM sfdata.clean_procol cp
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
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns one row per entreprise with an active procédure collective on a given date, choosing the most severe proceeding in case of conflict (liquidation > redressement > sauvegarde). Completed proceedings are excluded (action_procol = "fin_procedure" or "inclusion_autre_procedure").';

---- create above / drop below ----

DROP FUNCTION IF EXISTS procol_at_date;

-- Restore previous procol_at_date (migration 022 / 029: multiple rows per siren, no libelle_procol)
CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT) AS $$
BEGIN
  RETURN QUERY
  WITH last_action_procol AS (
    SELECT DISTINCT ON (cp.siren, cp.action_procol)
      cp.siren, cp.date_effet, cp.action_procol, cp.stade_procol
    FROM sfdata.clean_procol cp
    WHERE cp.date_effet <= date_param
    ORDER BY cp.siren, cp.action_procol, cp.date_effet DESC
  )
  SELECT lap.siren, lap.date_effet, lap.action_procol, lap.stade_procol
  FROM last_action_procol lap
  WHERE lap.action_procol != 'fin_procedure' AND lap.action_procol != 'inclusion_autre_procedure';
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns entreprises that have a procédure collective in progress on a given date. A single entreprise may have several simultaneous proceedings. Completed proceedings are not counted (action_procol = "fin_procedure" or action_procol = "inclusion_autre_procedure") — closed entreprises are nevertheless displayed.';

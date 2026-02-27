-- Update procol_at_date to return a single row per siren (most severe active procol),
-- and add libelle_procol to the return type.
--
-- Create clean_procol_at_date: a materialized view of (siren, periode) pairs where
-- an active procédure collective exists, restricted to the clean_filter perimeter.
-- Covers monthly periods from 2016-01-01 to today.
-- Absence from the view means "In bonis" by definition.
DROP FUNCTION procol_at_date;

CREATE OR REPLACE FUNCTION procol_at_date(date_param date)
RETURNS TABLE(siren VARCHAR(9), date_effet DATE, action_procol TEXT, stade_procol TEXT, libelle_procol TEXT) AS $$
  WITH last_action_procol AS (
    SELECT DISTINCT ON (siren, action_procol)
      siren, date_effet, action_procol, stade_procol, libelle_procol
    FROM clean_procol
    WHERE date_effet <= date_param
    ORDER BY siren, action_procol, date_effet DESC
  )
  SELECT DISTINCT ON (siren)
    siren, date_effet, action_procol, stade_procol, libelle_procol
  FROM last_action_procol
  WHERE action_procol != 'fin_procedure'
    AND action_procol != 'inclusion_autre_procedure'
    AND NOT (stade_procol = 'plan_continuation' AND AGE(date_param, date_effet) >= INTERVAL '10 years')
  ORDER BY siren,
    CASE action_procol
      WHEN 'liquidation' THEN 1
      WHEN 'redressement' THEN 2
      WHEN 'sauvegarde' THEN 3
      ELSE 4
    END;
$$ LANGUAGE SQL;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns one row per entreprise with an active procédure collective on a given date, choosing the most severe proceeding in case of conflict (liquidation > redressement > sauvegarde). Completed proceedings are excluded (action_procol = "fin_procedure" or "inclusion_autre_procedure").';

SET LOCAL work_mem TO '256MB';

CREATE MATERIALIZED VIEW clean_procol_at_date AS
WITH periodes AS (
  SELECT generate_series(
    '2016-01-01'::date,
    date_trunc('month', CURRENT_DATE)::date,
    '1 month'::interval
  )::date AS periode
),
sirens_procol AS (
  SELECT DISTINCT siren FROM clean_procol
),
-- Call procol_at_date once per period (not once per siren×period)
active_procol_by_periode AS (
  SELECT p.periode, pad.siren, pad.date_effet, pad.action_procol, pad.stade_procol, pad.libelle_procol
  FROM periodes p
  CROSS JOIN LATERAL procol_at_date(p.periode) pad
)
SELECT
  sp.siren,
  p.periode,
  apbp.date_effet,
  apbp.action_procol,
  apbp.stade_procol,
  apbp.libelle_procol
FROM sirens_procol sp
CROSS JOIN periodes p
LEFT JOIN active_procol_by_periode apbp
  ON apbp.siren = sp.siren
  AND apbp.periode = p.periode
WHERE apbp.action_procol IS NOT NULL
WITH DATA;

CREATE INDEX ON clean_procol_at_date (siren);

---- create above / drop below ----

DROP MATERIALIZED VIEW IF EXISTS clean_procol_at_date;

-- Restore previous procol_at_date (migration 022 / 029: multiple rows per siren, no libelle_procol)
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
  WHERE action_procol != 'fin_procedure' AND action_procol != 'inclusion_autre_procedure';
$$ LANGUAGE SQL;

COMMENT ON FUNCTION procol_at_date (date) IS 'Returns entreprises that have a procédure collective in progress on a given date. A single entreprise may have several simultaneous proceedings. Completed proceedings are not counted (action_procol = "fin_procedure" or action_procol = "inclusion_autre_procedure") — closed entreprises are nevertheless displayed.';

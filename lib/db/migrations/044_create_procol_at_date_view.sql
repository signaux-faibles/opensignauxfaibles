-- Create clean_procol_at_date: a materialized view of (siren, periode) pairs where
-- an active procédure collective exists, restricted to the clean_filter perimeter.
-- Covers monthly periods from 2016-01-01 to today.
-- Absence from the view means "In bonis" by definition.
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
  SELECT DISTINCT siren FROM public.clean_procol
),
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
WITH NO DATA;

CREATE INDEX ON clean_procol_at_date (siren);

---- create above / drop below ----

DROP MATERIALIZED VIEW IF EXISTS clean_procol_at_date;

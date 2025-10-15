CREATE MATERIALIZED VIEW IF NOT EXISTS filter AS
  WITH const AS (
    -- This constant is duplicated inside the code. Beware if changing.
    SELECT
      100 AS n_months,
      10 AS min_effectif,
      ARRAY['7490', '7430', '7470', '7410', '7379', '7348', '7346', '7210', '7220', '4140', '7373', '7366', '7389', '4110', '4120', '7383', '4160'] as excluded_categories
  )

  SELECT DISTINCT LEFT(e.siret, 9) AS siren
  FROM clean_effectif e
  CROSS JOIN const
  LEFT JOIN clean_sirene_ul s ON s.siren = LEFT(e.siret, 9)
  WHERE e.effectif >= const.min_effectif
    AND e.periode >= CURRENT_DATE - const.n_months * INTERVAL '1 month'
    AND NOT (s.statut_juridique = ANY(const.excluded_categories)) OR s.statut_juridique IS NULL;

CREATE UNIQUE INDEX filter_siren_index
    ON filter(siren);


---- create above / drop below ----

DROP MATERIALIZED VIEW IF EXISTS filter;


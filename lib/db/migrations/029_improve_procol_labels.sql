CREATE OR REPLACE VIEW clean_procol AS
SELECT
  LEFT(siret, 9) as siren,
  date_effet,
  action_procol,
  stade_procol,
  CASE
    WHEN stade_procol = 'fin_procedure' THEN 'In bonis'
    WHEN action_procol = 'redressement' AND stade_procol = 'plan_continuation' THEN 'Plan de redressement'
    WHEN action_procol = 'redressement' THEN 'Redressement judiciaire'
    WHEN action_procol = 'liquidation' THEN 'Liquidation judiciaire'
    WHEN action_procol = 'sauvegarde' AND stade_procol = 'plan_continuation' THEN 'Plan de sauvegarde'
    WHEN action_procol = 'sauvegarde' THEN 'Sauvegarde'
  END AS libelle_procol
FROM stg_procol p
WHERE NOT EXISTS (SELECT siren FROM siren_blacklist b WHERE LEFT(p.siret, 9) = b.siren)
  -- On ignore la clôture administrative de la procédure
  AND stade_procol != 'solde_procedure'
GROUP BY LEFT(siret, 9), date_effet, action_procol, stade_procol;

---- create above / drop below ----

-- Rollback to previous view (migration 019)
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

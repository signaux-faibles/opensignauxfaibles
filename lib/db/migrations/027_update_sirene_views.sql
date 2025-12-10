DROP VIEW IF EXISTS clean_sirene;
DROP VIEW IF EXISTS clean_sirene_ul;

-- Update clean_sirene view to use ape instead of code_activite and nomenclature_activite
CREATE OR REPLACE VIEW clean_sirene AS
  SELECT
    s.siren,
    s.siret,
    s.siege,
    -- adresse
    TRIM(
      CONCAT_WS(' ',
        COALESCE(s.numero_voie, ''),
        COALESCE(s.indrep, ''),
        COALESCE(s.type_voie, ''),
        COALESCE(s.voie, ''),
        CASE WHEN s.complement_adresse IS NOT NULL THEN CONCAT('(', s.complement_adresse, ')') END
        ) || E'\n' ||
      CONCAT_WS(' ',
        COALESCE(s.code_postal, ''),
        COALESCE(s.commune, s.commune_etranger, ''),
        CASE WHEN s.pays_etranger IS NOT NULL THEN CONCAT(E'\n', s.pays_etranger) END
      )
    ) AS adresse,
    s.code_commune,
    s.departement,
    s.ape,
    naf.niv1 AS naf_section,
    s.date_creation,
    s.longitude,
    s.latitude,
    s.est_actif
  FROM stg_sirene s
  LEFT JOIN naf_codes naf ON s.ape = naf.niv5
  WHERE s.siren IN (SELECT siren FROM clean_filter)
    AND s.departement IS NOT NULL;

-- Update clean_sirene_ul view to add ape column
CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    s.siren,
    s.raison_sociale,
    s.statut_juridique,
    s.activite_principale,
    naf.niv1 AS naf_section,
    siege.departement,
    s.creation,
    s.est_actif
  FROM stg_sirene_ul s
  LEFT JOIN naf_codes naf ON s.activite_principale = naf.niv5
  LEFT JOIN stg_sirene siege ON s.siren = siege.siren AND siege.siege = true
  WHERE s.siren IN (SELECT siren FROM clean_filter);



---- create above / drop below ----

-- Revert to previous clean_sirene view with code_activite and nomenclature_activite
CREATE OR REPLACE VIEW clean_sirene AS
  SELECT
    siren,
    siret,
    siege,
    -- adresse
    TRIM(
      CONCAT_WS(' ',
        COALESCE(numero_voie, ''),
        COALESCE(indrep, ''),
        COALESCE(type_voie, ''),
        COALESCE(voie, ''),
        CASE WHEN complement_adresse IS NOT NULL THEN CONCAT('(', complement_adresse, ')') END
        ) || E'\n' ||
      CONCAT_WS(' ',
        COALESCE(code_postal, ''),
        COALESCE(commune, commune_etranger, ''),
        CASE WHEN pays_etranger IS NOT NULL THEN CONCAT(E'\n', pays_etranger) END
      )
    ) AS adresse,
    code_commune,
    departement,
    ape,
    date_creation,
    longitude,
    latitude,
    est_actif
  FROM stg_sirene
  WHERE siren IN (SELECT siren FROM clean_filter);

-- Revert clean_sirene_ul view
CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    siren,
    raison_sociale,
    statut_juridique,
    activite_principale,
    creation,
    est_actif
  FROM stg_sirene_ul
  WHERE siren IN (SELECT siren FROM clean_filter);

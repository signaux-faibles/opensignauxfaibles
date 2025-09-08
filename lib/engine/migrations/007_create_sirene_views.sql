CREATE OR REPLACE VIEW sfdata_clean_sirene AS
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
    code_activite,
    nomenclature_activite,
    date_creation,
    longitude,
    latitude
FROM sfdata_stg_sirene;

CREATE OR REPLACE VIEW sfdata_clean_sirene_ul AS
  SELECT
    siren,
    raison_sociale,
    statut_juridique,
    creation
  FROM sfdata_stg_sirene_ul;

---- create above / drop below ----

DROP VIEW IF EXISTS sfdata_clean_sirene;
DROP VIEW IF EXISTS sfdata_clean_sirene_ul;

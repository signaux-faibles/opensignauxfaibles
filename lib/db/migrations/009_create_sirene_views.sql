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
    latitude
  FROM stg_sirene
  WHERE siren IN (SELECT siren FROM clean_filter);

CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    siren,
    raison_sociale,
    statut_juridique,
    creation
  FROM stg_sirene_ul
  WHERE siren IN (SELECT siren FROM clean_filter);

---- create above / drop below ----

DROP VIEW IF EXISTS clean_sirene;
DROP VIEW IF EXISTS clean_sirene_ul;

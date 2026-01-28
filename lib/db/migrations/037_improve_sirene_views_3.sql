
DROP VIEW IF EXISTS clean_sirene_ul;

CREATE VIEW clean_sirene_ul AS
  SELECT
    s.siren,
    CASE WHEN s.raison_sociale IS NOT NULL AND s.raison_sociale != ''
      THEN s.raison_sociale
      ELSE
        TRIM(
          CONCAT(
            s.nom_unite_legale || ' ',
            CASE WHEN s.nom_usage_unite_legale IS NOT NULL
              THEN s.nom_usage_unite_legale || ' '
            ELSE ''
            END,
            COALESCE(prenom1_unite_legale, '') || ' ',
            COALESCE(prenom2_unite_legale, '') || ' ',
            COALESCE(prenom3_unite_legale, '') || ' ',
            COALESCE(prenom4_unite_legale, '')
          )
        ) || ' EI'
      END as raison_sociale,
    s.categorie_juridique,
    sj.libelle as libelle_categorie_juridique,
    s.activite_principale,
    naf.niv5_libelle AS libelle_activite_principale,
    naf.niv1 AS naf_section,
    naf.niv1_libelle AS libelle_naf_section,
    siege.departement,
    ldr.libelle_departement,
    ldr.region,
    ldr.libelle_region,
    s.creation,
    s.est_actif
  FROM stg_sirene_ul s
  LEFT JOIN naf_codes naf ON s.activite_principale = naf.niv5
  LEFT JOIN categories_juridiques sj ON s.categorie_juridique = sj.code
  LEFT JOIN stg_sirene siege ON s.siren = siege.siren AND siege.siege = true
  LEFT JOIN label_departement_region ldr ON siege.departement = ldr.departement;

DROP VIEW IF EXISTS clean_sirene;

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
    ldr.libelle_departement,
    ldr.region,
    ldr.libelle_region,
    s.ape,
    naf.niv1 AS naf_section,
    s.date_creation,
    s.longitude,
    s.latitude,
    s.est_actif
  FROM stg_sirene s
  LEFT JOIN naf_codes naf ON s.ape = naf.niv5
  LEFT JOIN label_departement_region ldr
    ON s.departement = ldr.departement;


---- create above / drop below ----

-- Rollback to previous views
DROP VIEW IF EXISTS clean_sirene_ul;

CREATE VIEW clean_sirene_ul AS
  SELECT
    s.siren,
    CASE WHEN s.raison_sociale IS NOT NULL AND s.raison_sociale != ''
      THEN s.raison_sociale
      ELSE
        TRIM(
          CONCAT(
            s.nom_unite_legale || ' ',
            CASE WHEN s.nom_usage_unite_legale IS NOT NULL
              THEN s.nom_usage_unite_legale || ' '
            ELSE ''
            END,
            COALESCE(prenom1_unite_legale, '') || ' ',
            COALESCE(prenom2_unite_legale, '') || ' ',
            COALESCE(prenom3_unite_legale, '') || ' ',
            COALESCE(prenom4_unite_legale, '')
          )
        ) || ' EI'
      END as raison_sociale,
    s.categorie_juridique,
    sj.libelle as libelle_categorie_juridique,
    s.activite_principale,
    naf.niv5_libelle AS libelle_activite_principale,
    naf.niv1 AS naf_section,
    naf.niv1_libelle AS libelle_naf_section,
    siege.departement,
    s.creation,
    s.est_actif
  FROM stg_sirene_ul s
  LEFT JOIN naf_codes naf ON s.activite_principale = naf.niv5
  LEFT JOIN categories_juridiques sj ON s.categorie_juridique = sj.code
  LEFT JOIN stg_sirene siege ON s.siren = siege.siren AND siege.siege = true;

DROP VIEW IF EXISTS clean_sirene;

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
  LEFT JOIN naf_codes naf ON s.ape = naf.niv5;

-- Back to previous version of the view
CREATE OR REPLACE VIEW clean_sirene_ul AS
  SELECT
    s.siren,
    CASE WHEN s.raison_sociale IS NOT NULL
      THEN s.raison_sociale
      ELSE
        TRIM(
          CONCAT(
            s.nom_unite_legale,
            CASE WHEN s.nom_usage_unite_legale IS NOT NULL
              THEN s.nom_usage_unite_legale || '/'
            ELSE ''
            END,
            COALESCE(prenom1_unite_legale, ' '),
            COALESCE(prenom2_unite_legale, ' '),
            COALESCE(prenom3_unite_legale, ' '),
            COALESCE(prenom4_unite_legale, ' ')
          )
        ) || '/'
      END as raison_sociale,
    s.statut_juridique,
    s.activite_principale,
    naf.niv1 AS naf_section,
    siege.departement,
    s.creation,
    s.est_actif
  FROM stg_sirene_ul s
  LEFT JOIN naf_codes naf ON s.activite_principale = naf.niv5
  LEFT JOIN stg_sirene siege ON s.siren = siege.siren AND siege.siege = true;

---- create above / drop below ----

-- Back to previous version of the view
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
  LEFT JOIN stg_sirene siege ON s.siren = siege.siren AND siege.siege = true;

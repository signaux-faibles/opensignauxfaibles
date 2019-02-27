#' @import dplyr
#' @importFrom lubridate %m+% %m-%
NULL

#' Connexion à la base de donnée
#'
#' \code{connect_to_database} permet de requêter des données mongoDB pour en
#' faire un dataframe ou un Spark dataframe. \code{factor_request} permet de fabriquer la requête d'aggrégation correspondante.
#'
#' @param database Nom de la base de données Mongodb.
#' @param collection Nom de la collection Mongodb.
#' @param batch Valeur du batch, format "aamm".
#' @param algo Seule la valeur par défaut ("algo2") est supportée actuellement.
#' @param siren Vecteur de numéros siren. Si omis ou égal à NULL, tous les
#'   numéros siren sont extraits.
#' @param date_inf Limite inférieure de la période de temps requêtée
#' @param date_sup Limite supérieure de la période de temps requêtée
#' @param min_effectif Limite inférieure de la taille d'établissement requêtée,
#'   en nombre d'employés.
#' @param fields Liste des champs à requêter. Si omis ou égal à NULL, tous les
#'   champs sont requêtés.
#' @param code_ape code_ape à requêter, sous la forme "7112B"
#' @param type Prend comme valeur "dataframe" pour extraire un dataframe dans R
#'   (par défaut), "spark" pour extraire un spark dataframe.
#' @param subsample Randomly samples company/period pairs. If negative, keeps all objects.
#' @param limit Number of entries to take into account. For test purposes mainly.
#'
#' @return Renvoie un dataframe ou un spark dataframe selon la valeur de
#'   \code{type}
#' @export
#'
#' @examples connect_to_database("test_signauxfaibles","testing", "1812", date_inf = "2014-01-01", date_sup = "2019-01-01", type = "spark")

connect_to_database <- function(
  database,
  collection,
  batch,
  siren = NULL,
  date_inf = NULL,
  date_sup = NULL,
  min_effectif = 10,
  fields = NULL,
  code_ape = NULL,
  type = "dataframe",
  subsample = NULL,
  .limit = NULL
  ) {

  requete <- factor_request(
    batch,
    siren,
    date_inf,
    date_sup,
    min_effectif,
    fields,
    code_ape,
    subsample
    )

  #cat("\n")
  #cat(requete)
  #cat("\n")
  #browser()



  TYPES <- c("dataframe", "spark")

  assertthat::assert_that(type %in% TYPES,
    msg = "connect_to_database:
    specification du type d'import/export non valide")

    if (type == TYPES[1]) {

      cat("Connexion à la collection mongodb ", collection, " ...")

      dbconnection <- mongolite::mongo(
        collection = collection,
        db = database,
        verbose = TRUE,
        url = "mongodb://localhost:27017")
      cat(" Fini.", "\n")

      # Import dataframe
      cat("Import ...", "\n")

      donnees <- dbconnection$aggregate(requete)


      cat(" Fini.", "\n")

      if (dim(donnees)[1] == 0) {
        return(dplyr::data_frame(siret = character(0), periode = character(0)))
      }
      assertthat::assert_that(
        all(c("periode", "siret") %in% names(donnees)))

      assertthat::assert_that(
        anyDuplicated(donnees %>% select(siret, periode)) == 0,
        msg = "La base importee contient des lignes identiques"
        )

      table_wholesample <- donnees %>%
        arrange(periode) %>%
        mutate_if(lubridate::is.POSIXct, as.Date) %>%
        tibbletime::as_tbl_time(periode)

      n_eta <- table_wholesample$siret %>%
        n_distinct()
      n_ent <- table_wholesample$siret %>%
        stringr::str_sub(1, 9) %>%
        n_distinct()
      cat("Import de", n_eta, "etablissements issus de",
        n_ent, "entreprises", "\n")

      cat(" Fini.", "\n")

    } else if (type == TYPES[2]) {
      browser()
      cat("Ouverture de la connection spark", "\n")

      #FIX ME: neater way to deal with already open connection
      #spark_disconnect_all()
      sc <- sparklyr::spark_connection_find()[[1]]

      cat("Connection à MongoDB et import des donnees dans spark", "\n")
      load <- sparklyr::invoke(sparklyr::spark_session(sc), "read") %>%
        sparklyr::invoke("format", "com.mongodb.spark.sql.DefaultSource") %>%
        sparklyr::invoke("option", "pipeline", requete) %>%
        sparklyr::invoke("load")

      table_wholesample <- sparklyr::sdf_register(load)

      if (ncol(table_wholesample) == 0) {
        return(table_wholesample)
      }


      type_frame <- sparklyr::sdf_schema(table_wholesample) %>%
        lapply(function(x) do.call(tibble::data_frame, x)) %>%
        bind_rows()

      assertthat::assert_that(
        all(type_frame$type != "NullType"),
        msg = paste("NullTypes found in imported spark frame: ",
          type_frame$name[type_frame$type == "NullType"], collapse = " :: ")
        )

    }

    # Champs manquants
    champs_manquants <- fields[!fields %in% tbl_vars(table_wholesample)]
    if (length(champs_manquants) >= 1) {
      cat("Champs manquants: ")
      cat(champs_manquants, "\n")

      cat("Remplacements par NA", "\n")
      if (type == "dataframe") {
        NA_variable <- NA_character_
      } else if (type == "spark") {
        NA_variable <- "NA_character_"
      }
      remplacement <- paste0(NA_variable, character(length(champs_manquants))) %>%
        setNames(champs_manquants)
      table_wholesample <- table_wholesample %>%
        mutate_(.dots = remplacement)


    }


    return(table_wholesample)


}

###################################


#' @rdname connect_to_database
factor_request <- function(
  batch,
  siren,
  date_inf,
  date_sup,
  min_effectif,
  fields,
  code_ape,
  subsample,
  .limit = NULL
  ) {
  # Util functions

  make_query  <- function(keyword, x) {
    if (any(x != "")) return(paste0('{"$', keyword, '":{',
        paste0(x[x != ""], collapse = ', '), '}}'))
  else return("")
  }

  if (!is.null(.limit))
    limit_req  <- paste0('{"$limit":', .limit, "}")
  else
    limit_req <- ""
  ## Construction de la requête ##
  match_id  <- paste0('"info.batch":"', batch, '"')

  # Filtrage siren

  if (is.null(siren)) {
    match_siren <- ""
  } else {
    match_siren  <- c()
    for (i in seq_along(siren)) {
      match_siren  <- c(
        match_siren,
        paste0('{"info.siren":"', siren[i], '"}')
        )
    }

    match_siren <- paste0('"$or":[', paste(match_siren, collapse = ","), "]")
  }




  # Unwind du tableau
  unwind_req <- '{"$unwind":{"path": "$value"}}'

  # Sample
  if (is.null(subsample)) {
    sample_req <- ""
  } else {
    sample_req <- paste0('{"$sample": {"size": ', subsample, '}}')
  }

  # Filtrage code APE

  if (is.null(code_ape)) {
    match_APE <- ""
  } else {
    niveau_code_ape <- nchar(code_ape)
    if (any(niveau_code_ape >= 2)) {
      ape_or_naf <- "code_ape"
    } else {
      ape_or_naf <- "code_naf"
    }
    match_APE  <- c()
    for (i in seq_along(code_ape)) {
      match_APE  <- c(
        match_APE,
        paste0('{"value.', ape_or_naf, '":
          {"$regex":"^', code_ape[i], '", "$options":"i"}
    }')
    )
  }
  match_APE <- paste0('"$or":[', paste(match_APE, collapse = ","), "]")
  # FIX ME: Requete nettement sous-optimale
}

# Filtrage effectif
if (is.null(min_effectif)) {
  match_eff <- ""
} else {
  match_eff <- paste0(
    '"value.effectif":{"$gte":',
    min_effectif, "}")
}

# Filtrage date
if (is.null(date_inf)) {
  match_date_1  <- ""
} else {
  match_date_1  <- paste0('"info.periode":{
    "$gte": {"$date":"', date_inf, 'T00:00:00Z"}}')
}

# Filtrage date
if (is.null(date_sup)) {
  match_date_2  <- ""
} else {
  match_date_2  <- paste0('"info.periode":{"$lt": {"$date":"', date_sup, 'T00:00:00Z"}}')

}

match_req <- make_query(
  "match",
  c(
    match_id,
    match_siren,
    match_date_1,
    match_date_2,
    match_APE,
    match_eff
    ))

# Construction de la projection

if (is.null(fields)) {
  projection_req  <- ""
} else {
  assertthat::validate_that( #does not return an error if not verified
    all(c("periode", "siret") %in% fields),
    msg = "Beware: siret and periode are not included in the request")
  projection_req  <- paste0('"', fields, '":1')
  projection_req  <- paste(projection_req, collapse = ",")
  projection_req  <- paste0('{"$project":{', projection_req, "}}")
}

# Remplacement de la racine

new_root <- "{\"$replaceRoot\" : {\"newRoot\": \"$value\"}}"
#new_root = ""

reqs <- c(
  #match_req_1,
  #unwind_req,
  limit_req,
  match_req,
  sample_req,
  new_root,
  projection_req
  )

requete  <- paste(
  reqs[reqs != ""],
  collapse = ", ")
requete <- paste0(
  "[",
  requete,
  "]")

return(requete)

}

#' Wrapper pour se connecter à H2O et spark avec la bonne configuration
#'
#' Créé une instance ou rejoint l'instance existante le cas échéant.
#' N'utiliser qu'avec des R dataframes, sinon rsparkling s'occupe de créer une instance.
#'
#' @return
#' @export
#'
#' @examples
connect_to_h2o <- function() {
  h2o::h2o.init(
    ip = "localhost",
    port = 4444,
    min_mem_size = "5G",
    log_dir = rprojroot::find_rstudio_root_file("logs"))
}


#' @rdname connect_to_h2o
connect_to_spark <- function(database = NULL, collection = NULL) {
  config <- sparklyr::spark_config()
  config$sparklyr.defaultPackages <-
    c("org.mongodb.spark:mongo-spark-connector_2.11:2.4.0")
  if (!is.null(database) && !is.null(collection)) {
    config$spark.mongodb.input.uri <-
      paste0("mongodb://127.0.0.1:27017/", database, ".", collection)
    config$spark.mongodb.output.uri <-
      paste0("mongodb://127.0.0.1:27017/", database, ".", collection)
  }

  sc <- sparklyr::spark_connect(master = "local[*]", config = config)

  if ("spark.mongodb.input.uri" %in% sc$config &&
    config$spark.mongodb.input.uri != sc$config$spark.mongodb.input.uri) {
    sparklyr::spark_disconnect_all()
    sc <- sparklyr::spark_connect(master = "local[*]", config = config)
  }
  #FIX ME: add config if spark has been initiated without this config file.
  #Spark_connect does not correct this automatically.

  return(sc)

}


#' Exporte une requête de base mongodb as a csv file
#'
#' Attention, fonction non maintenue car plus utilisée. Faisais maladroitement appel à un fichier sh. Les données transitent dorénavant par Spark.
#'
#' @param database Nom de la base de donnée mongodb
#' @param batch Numéro de batch, au format "AAMM"
#' @param fields Liste de champs à extraire.
#' @param min_effectif Limite minimale du filtrage sur l'effectif. NE FONCTIONNE PAS. POUR L'INSTANT EFFECTIF_MIN = 10.
#'
#' @return
#'
#'
#' @examples
export_to_csv <- function(database, algo, batch, fields, min_effectif) {
  path_1 <- rprojroot::find_rstudio_root_file(
    "..", "dbmongo", "export", "export.sh"
    )

  path_2 <- rprojroot::find_rstudio_root_file(
    "..", "dbmongo", "export", "export_fields.txt"
    )

  ## Write fields to file path_1
  write(fields, path_2, append = FALSE, sep = "\n")

  ##

  # FIX ME: ignore complètement min_effectif !! (par défaut, min_effectif = 10)
  system2("bash", args = c(path_1, database, algo, batch, min_effectif))
}


#' Get a list of field names
#'
#' Utilitary function to quickly get a field name list. Put training = TRUE to
#' get only fields used for training. Put an input to 0 to remove the related
#' fields, to 1 to keep basic fields, and to 2 to include all fields.
#'
#' @param training
#' @param siren
#' @param urssaf
#' @param delai
#' @param effectif
#' @param diane
#' @param bdf
#' @param apart
#' @param procol
#' @param interim
#' @param target_encode Only available for training
#' @param info Exercices Diane & Bdf
#'
#' @return  Vector of field names
#' @export
#'
#' @examples
get_fields <- function(
  training,
  siren = 2,
  urssaf = 2,
  delai = 2,
  effectif = 2,
  note_preface = 2,
  diane = 2,
  bdf = 2,
  apart = 2,
  procol = 1,
  interim = 0,
  target_encode = 2,
  info = 0
  ) {
  # TODO change this into data-frame !
  fields  <- c()
  if (siren >= 1 && !training) {
    fields  <- c(fields,
      "siret",
      "siren",
      "periode",
      "code_ape",
      "code_ape_niveau2",
      "code_ape_niveau3",
      "code_naf"
      )
  }
  if (siren >= 1) {
    fields <- c(
      fields,
      "age"
      )
  }

  if (urssaf >= 1) {
    fields  <- c(fields,
      "montant_part_patronale",
      "montant_part_ouvriere",
      "montant_echeancier",
      "delai",
      "duree_delai",
      "ratio_dette",
      "ratio_dette_moy12m",
      "cotisation_moy12m"
      )
  }
  if (urssaf >= 2) {
    fields  <- c(fields,
      "montant_part_patronale_past_1",
      "montant_part_ouvriere_past_1",
      "montant_part_patronale_past_2",
      "montant_part_ouvriere_past_2",
      "montant_part_patronale_past_3",
      "montant_part_ouvriere_past_3",
      "montant_part_patronale_past_6",
      "montant_part_ouvriere_past_6",
      "montant_part_patronale_past_12",
      "montant_part_ouvriere_past_12"
      )
  }

  if (apart >= 1) {
    fields <- c(
      fields,
      "apart_heures_consommees",
      "apart_heures_autorisees"
      )
  }

  if (effectif >= 1) {
    fields <- c(fields,
      "effectif",
      "effectif_entreprise"
      )
  }
  if (effectif >= 2) {
    fields <- c(fields,
      "effectif_past_6",
      "effectif_past_12",
      "effectif_past_18",
      "effectif_past_24"
      )
  }

  if (note_preface >= 1) {
    fields <- c(fields,
      "note_preface",
      "note_preface_past_1",
      "note_preface_past_2"
    )
  }

  if (diane >= 1) {


    fields <- c(fields,
      "dette_fiscale_et_sociale",
      "effectif_consolide",
      "frais_de_RetD",
      "conces_brev_et_droits_sim",
      "nombre_etab_secondaire",
      "nombre_filiale",
      "taille_compo_groupe",
      "concours_bancaire_courant",
      "equilibre_financier",
      "independance_financiere",
      "endettement",
      "autonomie_financiere",
      "degre_immo_corporelle",
      "financement_actif_circulant",
      "liquidite_generale",
      "liquidite_reduite",
      "rotation_stocks",
      "credit_client",
      "credit_fournisseur",
      "taux_interet_financier",
      "taux_interet_sur_ca",
      "endettement_global",
      "taux_endettement",
      "capacite_remboursement",
      "capacite_autofinancement",
      "couverture_ca_fdr",
      "couverture_ca_besoin_fdr",
      "poids_bfr_exploitation",
      "exportation",
      "efficacite_economique",
      "productivite_potentiel_production",
      "productivite_capital_financier",
      "productivite_capital_investi",
      "taux_d_investissement_productif",
      "rentabilite_economique",
      "performance",
      "rendement_brut_fonds_propres",
      "rentabilite_nette",
      "rendement_capitaux_propres",
      "rendement_ressources_durables",
      "taux_marge_commerciale",
      "taux_valeur_ajoutee",
      "part_salaries",
      "part_etat",
      "part_preteur",
      "part_autofinancement",
      "ca",
      "ca_exportation",
      "achat_marchandises",
      "achat_matieres_premieres",
      "production",
      "marge_commerciale",
      "consommation",
      "autres_achats_charges_externes",
      "valeur_ajoutee",
      "charge_personnel",
      "impots_taxes",
      "subventions_d_exploitation",
      "excedent_brut_d_exploitation",
      "autres_produits_charges_reprises",
      "dotation_amortissement",
      "resultat_expl",
      "operations_commun",
      "produits_financiers",
      "charges_financieres",
      "interets",
      "resultat_avant_impot",
      "produit_exceptionnel",
      "charge_exceptionnelle",
      "participation_salaries",
      "impot_benefice",
      "benefice_ou_perte"
      )
  }

  if (diane >= 2) {
    fields  <- c(
      fields,
      # DIANE PAST 1
      "effectif_consolide_past_1",
      "dette_fiscale_et_sociale_past_1",
      "frais_de_RetD_past_1",
      "conces_brev_et_droits_sim_past_1",
      "nombre_etab_secondaire_past_1",
      "nombre_filiale_past_1",
      "taille_compo_groupe_past_1",
      "concours_bancaire_courant_past_1",
      "equilibre_financier_past_1",
      "independance_financiere_past_1",
      "endettement_past_1",
      "autonomie_financiere_past_1",
      "degre_immo_corporelle_past_1",
      "financement_actif_circulant_past_1",
      "liquidite_generale_past_1",
      "liquidite_reduite_past_1",
      "rotation_stocks_past_1",
      "credit_client_past_1",
      "credit_fournisseur_past_1",
      "taux_interet_financier_past_1",
      "taux_interet_sur_ca_past_1",
      "endettement_global_past_1",
      "taux_endettement_past_1",
      "capacite_remboursement_past_1",
      "capacite_autofinancement_past_1",
      "couverture_ca_fdr_past_1",
      "couverture_ca_besoin_fdr_past_1",
      "poids_bfr_exploitation_past_1",
      "exportation_past_1",
      "efficacite_economique_past_1",
      "productivite_potentiel_production_past_1",
      "productivite_capital_financier_past_1",
      "productivite_capital_investi_past_1",
      "taux_d_investissement_productif_past_1",
      "rentabilite_economique_past_1",
      "performance_past_1",
      "rendement_brut_fonds_propres_past_1",
      "rentabilite_nette_past_1",
      "rendement_capitaux_propres_past_1",
      "rendement_ressources_durables_past_1",
      "taux_marge_commerciale_past_1",
      "taux_valeur_ajoutee_past_1",
      "part_salaries_past_1",
      "part_etat_past_1",
      "part_preteur_past_1",
      "part_autofinancement_past_1",
      "ca_past_1",
      "ca_exportation_past_1",
      "achat_marchandises_past_1",
      "achat_matieres_premieres_past_1",
      "production_past_1",
      "marge_commerciale_past_1",
      "consommation_past_1",
      "autres_achats_charges_externes_past_1",
      "valeur_ajoutee_past_1",
      "charge_personnel_past_1",
      "impots_taxes_past_1",
      "subventions_d_exploitation_past_1",
      "excedent_brut_d_exploitation_past_1",
      "autres_produits_charges_reprises_past_1",
      "dotation_amortissement_past_1",
      "resultat_expl_past_1",
      "operations_commun_past_1",
      "produits_financiers_past_1",
      "charges_financieres_past_1",
      "interets_past_1",
      "resultat_avant_impot_past_1",
      "produit_exceptionnel_past_1",
      "charge_exceptionnelle_past_1",
      "participation_salaries_past_1",
      "impot_benefice_past_1",
      "benefice_ou_perte_past_1",
      # DIANE PAST 2
      "effectif_consolide_past_2",
      "dette_fiscale_et_sociale_past_2",
      "frais_de_RetD_past_2",
      "conces_brev_et_droits_sim_past_2",
      "nombre_etab_secondaire_past_2",
      "nombre_filiale_past_2",
      "taille_compo_groupe_past_2",
      "concours_bancaire_courant_past_2",
      "equilibre_financier_past_2",
      "independance_financiere_past_2",
      "endettement_past_2",
      "autonomie_financiere_past_2",
      "degre_immo_corporelle_past_2",
      "financement_actif_circulant_past_2",
      "liquidite_generale_past_2",
      "liquidite_reduite_past_2",
      "rotation_stocks_past_2",
      "credit_client_past_2",
      "credit_fournisseur_past_2",
      "taux_interet_financier_past_2",
      "taux_interet_sur_ca_past_2",
      "endettement_global_past_2",
      "taux_endettement_past_2",
      "capacite_remboursement_past_2",
      "capacite_autofinancement_past_2",
      "couverture_ca_fdr_past_2",
      "couverture_ca_besoin_fdr_past_2",
      "poids_bfr_exploitation_past_2",
      "exportation_past_2",
      "efficacite_economique_past_2",
      "productivite_potentiel_production_past_2",
      "productivite_capital_financier_past_2",
      "productivite_capital_investi_past_2",
      "taux_d_investissement_productif_past_2",
      "rentabilite_economique_past_2",
      "performance_past_2",
      "rendement_brut_fonds_propres_past_2",
      "rentabilite_nette_past_2",
      "rendement_capitaux_propres_past_2",
      "rendement_ressources_durables_past_2",
      "taux_marge_commerciale_past_2",
      "taux_valeur_ajoutee_past_2",
      "part_salaries_past_2",
      "part_etat_past_2",
      "part_preteur_past_2",
      "part_autofinancement_past_2",
      "ca_past_2",
      "ca_exportation_past_2",
      "achat_marchandises_past_2",
      "achat_matieres_premieres_past_2",
      "production_past_2",
      "marge_commerciale_past_2",
      "consommation_past_2",
      "autres_achats_charges_externes_past_2",
      "valeur_ajoutee_past_2",
      "charge_personnel_past_2",
      "impots_taxes_past_2",
      "subventions_d_exploitation_past_2",
      "excedent_brut_d_exploitation_past_2",
      "autres_produits_charges_reprises_past_2",
      "dotation_amortissement_past_2",
      "resultat_expl_past_2",
      "operations_commun_past_2",
      "produits_financiers_past_2",
      "charges_financieres_past_2",
      "interets_past_2",
      "resultat_avant_impot_past_2",
      "produit_exceptionnel_past_2",
      "charge_exceptionnelle_past_2",
      "participation_salaries_past_2",
      "impot_benefice_past_2",
      "benefice_ou_perte_past_2"
      )
  }
  if (bdf >= 1) {
    fields <- c(
      fields,
      "taux_marge",
      "delai_fournisseur",
      "poids_frng",
      "financier_court_terme",
      "frais_financier",
      "dette_fiscale"
      )
  }

  if (bdf >= 2) {
    fields  <- c(
      fields,
      #BDF PAST 1
      "taux_marge_past_1",
      "delai_fournisseur_past_1",
      "poids_frng_past_1",
      "financier_court_terme_past_1",
      "frais_financier_past_1",
      "dette_fiscale_past_1",
      #BDF PAST 2
      "taux_marge_past_2",
      "delai_fournisseur_past_2",
      "poids_frng_past_2",
      "financier_court_terme_past_2",
      "frais_financier_past_2",
      "dette_fiscale_past_2"
      )
  }

  if (procol >= 1) {
    fields  <- c(fields,
      "etat_proc_collective"
      )
  }
  if (procol >= 1  && !training) {
    fields  <- c(fields,
      # ALTARES
      "outcome",
      "time_til_outcome"
      )
  } else if (procol >= 2 && !training) {
    fields <- c(fields,
      "tag_outcome",
      "tag_failure"
      )
  }

  if (interim >= 1) {
    fields <- c(
      fields,
      "interim_proportion",
      "interim_ratio_past_6",
      "interim_ratio_past_12",
      "interim_ratio_past_18",
      "interim_ratio_past_24"
      )
  }

  if (target_encode >= 1 && training) {
    fields <- c(
      fields,
      "TargetEncode_code_ape_niveau2",
      "TargetEncode_code_ape_niveau3"
      )
  }
  if (info >= 1 && !training) {
    fields <- c(
      fields,
      "region",
      "departement",
      "arrete_bilan_diane",
      "exercice_bdf",
      "exercice_diane"
      )
  }
  return(fields)
}


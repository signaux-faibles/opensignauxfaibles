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
#' @param limit Randomly samples company/period pairs. If negative, keeps all objects.
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
  limit = NULL
) {

  requete <- factor_request(
    batch,
    siren,
    date_inf,
    date_sup,
    min_effectif,
    fields,
    code_ape,
    limit
  )

  #cat("\n")
  #cat(requete)
  #cat("\n")



  TYPES <- c("dataframe", "spark")

  assertthat::assert_that(type %in% TYPES,
                          msg = "connect_to_database:
    specification du type d'import/export non valide")

  if (type == TYPES[1]) {

    cat("Connexion à la collection mongodb ...")

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
  limit
) {
  # Util functions

  make_query  <- function(keyword, x) {
    if (any(x != "")) return(paste0('{"$', keyword, '":{',
                                    paste0(x[x != ""], collapse = ', '), '}}'))
    else return("")
  }

  ## Construction de la requête ##
  match_id  <- paste0('"_id.batch":"', batch, '"')

  # Filtrage siren

  if (is.null(siren)) {
    match_siren <- ""
  } else {
    match_siren  <- c()
    for (i in seq_along(siren)) {
      match_siren  <- c(
        match_siren,
        paste0('{"_id.siren":"', siren[i], '"}')
      )
    }

    match_siren <- paste0('"$or":[', paste(match_siren, collapse = ","), "]")
  }



  match_req_1 <- make_query("match", c(match_id, match_siren))

  # Unwind du tableau
  unwind_req <- '{"$unwind":{"path": "$value"}}'

  # Sample
  if (is.null(limit)) {
    sample_req <- ""
  } else {
    sample_req <- paste0('{"$sample": {"size": ', limit, '}}')
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
    match_date_1  <- paste0('"value.periode":{
      "$gte": {"$date":"', date_inf, 'T00:00:00Z"}}')
  }

  # Filtrage date
  if (is.null(date_sup)) {
    match_date_2  <- ""
  } else {
    match_date_2  <- paste0('"value.periode":{"$lt": {"$date":"', date_sup, 'T00:00:00Z"}}')

  }

  match_req_2 <- make_query("match",
                            c(match_APE, match_eff, match_date_1, match_date_2))

  # Construction de la projection

  if (is.null(fields)) {
    projection_req  <- ""
  } else {
    assertthat::assert_that(
      all(c("periode", "siret") %in% fields))
    projection_req  <- paste0('"', fields, '":1')
    projection_req  <- paste(projection_req, collapse = ",")
    projection_req  <- paste0('{"$project":{', projection_req, "}}")
  }

  # Remplacement de la racine

  new_root <- "{\"$replaceRoot\" : {\"newRoot\": \"$value\"}}"
  #new_root = ""

  reqs <- c(
    match_req_1,
    unwind_req,
    sample_req,
    match_req_2,
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

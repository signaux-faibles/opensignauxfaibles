#' Prepare_for_export
#'
#' If number of periods exceeds two, last period is dropped.
#'
#' @param donnees
#' @param collection
#' @param export_fields
#' @param database
#' @param last_batch
#'
#' @return
#' @export
#'
#' @examples
format_for_export <- function(
  donnees,
  collection,
  export_fields,
  database,
  last_batch,
  known_sirens_filenames = c("sirets_connus_pdl.csv", "sirets_connus_bfc.csv"),
  verbose = TRUE) {

  require(logger)
  if (verbose) {
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  prediction <- donnees %>%
    select(siret, periode, score)

  all_periods <- prediction %>%
    arrange(periode) %>%
    .$periode %>%
    unique()

  if (length(all_periods) < 2) {
    log_warning(
      "Less than two periods do not allow to compute score variations,
      or to monitor company appearing since last period. You can add more
      periods with the rollback_months parameter from load_new_data function"
    )
  }

  # Compute additional indicators about progression
  pred_data <- prediction %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(
      last_score = dplyr::lag(score)#,
      # last_periode = dplyr::lag(periode),
      # next_periode = dplyr::lead(periode),
      ) %>%
    ungroup() %>%
    # select(-c(last_periode, next_periode)) %>%
    mutate(score_diff = score - last_score) %>%
    select(-c(last_score))

  if (length(all_periods) >= 2){
    pred_data <- pred_data %>%
      filter(periode > min(all_periods), periode <= max(all_periods))
  }

  first_period <- min(all_periods)
  last_period <- max(all_periods)
  log_info("Préparation à l'export ... ")
  log_info("Dernière période connue: {last_period}")

  donnees <- connect_to_database(
    database,
    collection,
    last_batch,
    date_inf = first_period,
    date_sup = last_period %m+% months(1),
    min_effectif = 10,
    fields = export_fields[!export_fields %in% c(
      "connu", "score_diff", "score"
      )]
    )

  donnees <- pred_data %>%
    mutate(siret = as.character(siret)) %>%
    left_join(donnees %>% mutate(siret = as.character(siret)),
      by =
        c("siret", "periode")
      ) %>%
    dplyr::mutate(CCSF = date_ccsf) %>%
    dplyr::arrange(dplyr::desc(score), siret, dplyr::desc(periode))

  # Report des dernières infos financieres connues

  donnees <- donnees %>%
    mark_known_sirets(names = known_sirens_filenames) %>%
    select(export_fields)

  all_names <- names(donnees)
  log_info("Les variables suivantes sont absentes du dataframe:
    {export_fields[!(export_fields %in% all_names)]}")
    export_fields <- export_fields[export_fields %in% all_names]


    # if (is.emp)
    to_export <- donnees %>%
      dplyr::select(one_of(export_fields))
    return(to_export)
}



#' Export des scores dans une collection mongodb
#'
#' Exporte les scores vers une collection mongodb à partir des données formattées par la fonction
#' \code{\link{format_for_export}}.
#'
#'
#' @param formatted_data `data.frame()` \cr Données avec les champs "siret",
#'   "periode", "score" et "score_diff". C'est le cas des données formatées par
#'   \code{\link{format_for_export}}.
#' @param algo `character(1)` \cr Nom de l'algo qui figurera dans les objets
#'   exportés
#' @param batch `character(1)` \cr Nom du batch qui figurera dans les objets
#'   exportés
#' @param f_scores `character(2)` \cr Vecteur de scores F1 et F2. Doit être
#'   nommé avec comme noms "F1" et "F2"
#' @param database `character(1)` \cr Nom de la base de données vers laquelle
#'   param exporter.
#' @param collection `character(1)' \cr Nom de la collection vers laquelle
#'   exporter.
#'
#' @return Retourne TRUE. \cr Les objets insérés dans la base de données ont les
#'   champs:
#'   * "_id" (ObjectId générés),
#'   * "alert", qui peut prendre les valeurs _Alerte seuil F1_, _Alerte seuil F2_ et _Pas d'alerte_,
#'   * "algo" et "batch" tel qu'entrés en paramètres,
#'   * "siret", "periode", "score" et "score_diff" tel qu'extraits de la table \code{formatted_data},
#'   * "timestamp" qui donne la date et l'heure.
#'
#' @export
#'
#' @examples
export_scores_to_mongodb <- function(
  formatted_data,
  algo,
  batch,
  f_scores,
  database,
  collection
  ){

  exported_columns <- c("siret", "periode", "score", "score_diff")
  assertthat::assert_that(
    all(exported_columns %in% names(formatted_data)),
    msg = paste(
      paste0(exported_columns, collapse = ", "),
      "are compulsary column names to export the scores to mongodb"
      )
    )
  assertthat::assert_that(is.character(batch) && length(batch) == 1,
    msg = "Batch shoud be a length 1 character vector"
    )
  assertthat::assert_that(is.character(database) && length(database) == 1,
    msg = "Database shoud be a length 1 character vector"
    )
  assertthat::assert_that(is.character(collection) && length(collection) == 1,
    msg = "Collection shoud be a length 1 character vector"
    )

  assertthat::assert_that(length(f_scores) == 2)
  assertthat::assert_that(all(c("F1", "F2") %in% names(f_scores)))

  dbconnection <- mongolite::mongo(
    collection = collection,
    db = database,
    verbose = FALSE,
    # TODO: à lire du fichier de config
    url = "mongodb://localhost:27017"
    )

  data_to_export <- formatted_data %>%
    dplyr::select(dplyr::one_of(exported_columns)) %>%
    dplyr::mutate(
      algo = algo,
      batch = batch,
      timestamp = Sys.time()
      )
  data_to_export <- data_to_export %>%
    mutate(alert = alert_levels(score, f_scores["F1"], f_scores["F2"]))

  dbconnection$insert(data_to_export)

  return(TRUE)
}

#' Export des scores dans un tableau csv
#'
#' Exporte les scores vers un tableau csv à partir des données formattées par la fonction
#' \code{\link{format_for_export}}.
#'
#' @param formatted_data `data.frame()` \cr Données à exporter. Typiquement,
#'   le data.frame retourné par la fonction \code{\link{format_for_export}}.
#' @param algo `character(1)` \cr Nom de l'algo qui figurera dans les objets
#'   exportés
#' @param batch `character(1)` \cr Nom du batch qui figurera dans les objets
#'   exportés
#' @param relative_path `character(1)` \cr Chemin relatif du dossier dans lequel effectuer
#'   l'export, par rapport à la racine du paquet R.
#'
#'
#' @return Retourne TRUE. \cr Exporte un fichier csv, dans le dossier spécifié
#'   en entrée, nommé automatiquement "aaaa_mm_jj_export_{algo}_{batch}" avec
#'   éventuellement un suffixe "_vX" pour ne pas écraser de fichier existant.
#'
#' @export
#'
#' @examples
export_scores_to_csv  <- function(
  formatted_data,
  algo,
  batch,
  relative_path
  ){

  assertthat::assert_that(is.character(batch) && length(batch) == 1,
    msg = "Batch shoud be a length 1 character vector"
    )

  data_to_export <- formatted_data

  # create unique filename
  full_path <- name_file(
    relative_path,
    file_detail = paste("export", algo, batch, sep = "_"),
    file_extension = "csv",
    full_path = TRUE
    )

  write.table(data_to_export,
    row.names = F,
    dec = ",",
    sep = ";",
    file = full_path,
    quote = T,
    append = F
    )

  return(TRUE)
}

export_company_card  <- function(formatted_data){

  # TODO

}

#' Exports a dataframe to csv, or export scores to mongodb
#'
#' @param donnees Data to export.
#' @param batch Batch name.
#' @param algo Algorithm name.
#' @param database Database name. Database is created if non-existent. Only
#' used for mongodb export
#' @param collection Collection name.  Collection is created if non-existent. Only
#' used for mongodb export.
#' @param destination Either "csv", "mongodb" or "json" (case insensitive). See details.
#' @param relative_path Relative path used for csv export.
#' @param F_scores Named vector with F-scores, and names F1 and F2. Allows to
#' add the alert status.
#'
#' @section Details of csv export:
#' If destination is "csv", then a full csv export is operated.
#'
#' @section Details of mongodb export:
#' If destination is "mongodb", then the scores are exported to mongodb
#' \code(database) and \code(collection).  Input data \code(donnees) needs to
#' have columns "siret", "periode" and "score".
#'
#' Inserts object with fields "siret", "periode", "score", "batch" and
#' "timestamp" which gives the timestamp of insertion.
#'
#' @return
#' @export
#'
#' @examples
# export_scores <- function(
#   donnees,
#   batch,
#   algo = "algo",
#   database = "test_signauxfaibles",
#   collection = "Scores",
#   destination = "csv",
#   relative_path = file.path("..", "output"),
#   F_scores = NULL) {

#   assertthat::assert_that(tolower(destination) %in% c("csv", "mongodb", "json"),
#     msg = "Wrong export destination argument.
#     Possible destinations are 'csv', 'mongodb' or 'json'"
#     )

#   assertthat::assert_that(is.character(batch) && length(batch) == 1,
#     msg = "Batch shoud be a length 1 character vector"
#     )
#   assertthat::assert_that(is.character(database) && length(database) == 1,
#     msg = "Database shoud be a length 1 character vector"
#     )
#   assertthat::assert_that(is.character(collection) && length(collection) == 1,
#     msg = "Collection shoud be a length 1 character vector"
#     )

#   assertthat::assert_that(grepl("^[[:alnum:]_]+$", algo),
#     msg = "Please use alphanumeric and underscore characters for algorithm name"
#     )

#   if (tolower(destination) == "csv") {
#     fullpath <- name_file(
#       relative_path,
#       file_detail = paste0("detection", batch),
#       file_extension = "csv",
#       full_path = TRUE
#       )

#     write.table(donnees,
#       row.names = F,
#       dec = ",",
#       sep = ";",
#       file = fullpath,
#       quote = T,
#       append = F
#       )
#   } else if (tolower(destination) == "mongodb") {
#     dbconnection <- mongolite::mongo(
#       collection = collection,
#       db = database,
#       verbose = TRUE,
#       url = "mongodb://localhost:27017"
#       )

#     compulsory_columns <- c("siret", "periode", "score", "score_diff")

#     assertthat::assert_that(all(compulsory_columns %in% names(donnees)),
#       msg = paste(
#         paste0(compulsory_columns, collapse = ", "),
#         "are compulsary column names to export the scores to mongodb"
#         )
#       )

#     donnees_export <- donnees %>%
#       dplyr::select(siret, periode, score, score_diff) %>%
#       dplyr::mutate(
#         batch = batch,
#         timestamp = Sys.time(),
#         algo = algo
#         ) %>%
#       dplyr::rename(score = score)

#     if (!is.null(F_scores)) {
#       assertthat::assert_that(length(F_scores) == 2)
#       assertthat::assert_that(all(c("F1", "F2") %in% names(F_scores)))

#       donnees_export <- donnees_export %>%
#         mutate(alert = alert_levels(score, F_scores["F1"], F_scores["F2"]))
#     }

#     dbconnection$insert(donnees_export)
#   } else if (tolower(destination) == "json") {
#     error("Export to json not implemented yet !")
#   }
#   return()
# }

#' Marks sirets from files as "known"
#'
#' @param df
#' @param names File names
#'
#' @return
#' @export
#'
#' @examples
mark_known_sirets <- function(df, names) {
  sirens <- c()
  for (name in names) {
    sirets <- readLines(
      rprojroot::find_rstudio_root_file("..", "data-raw", name)
    )
    sirens <- c(sirens, substr(sirets, 1, 9))
  }
  df <- df %>%
    dplyr::mutate(connu = as.numeric(substr(siret, 1, 9) %in% sirens))

  return(df)
}

#' Get the scores for specific sirets and periods
#'
#' @param database String. Mongodb database name to query.
#' @param collection String. Mongodb collection name to query.
#' @param algo. Algorithm from which to take the score.
#' @param method String, either "first" or "last". See details.
#' @param siret Vector of sirets to query.
#'
#' @section Details:
#' If method is "first", then the last score available for the first batch
#' scoring this period.
#'
#' @return Returns a dataframe with columns "siret", "periode", "score"
#' (failure probability), "batch" (batch of the model used), "algo" (name of
#' the algorithm)
#'
#' @export
#'
#' @examples
get_scores <- function(
  database = "test_signauxfaibles",
  collection = "Scores",
  algo = "algo",
  method = "first",
  batchs = NULL,
  sirets = NULL
  ){

  assertthat::assert_that(is.character(database) && length(database) == 1,
    msg = "Database name should be a length 1 string"
  )
  assertthat::assert_that(is.character(collection) && length(collection) == 1,
    msg = "Collection name should be a length 1 string"
  )
  assertthat::assert_that(tolower(method) %in% c("first", "last"))

  assertthat::assert_that(grepl("^[[:alnum:]_]+$", algo),
    msg = "Please use alphanumeric and underscore characters for algorithm name"
  )

  assertthat::assert_that(all(
    grepl("^[0-9]+$", sirets)
  ),
  msg = "Sirets should only be numeric characters"
  )


  result <- data.frame(
    siret = character(),
    periode = as.Date(character()),
    score = double(),
    batch = character()
  )

  admin <- mongolite::mongo(db = "admin")
  db_names <- admin$run('{"listDatabases":1}')$databases$name
  assertthat::assert_that(database %in% db_names,
    msg = "Database not found"
  )

  db <- mongolite::mongo(db = database, collection = collection)
  col_names <- db$run('{"listCollections":1}')$cursor$firstBatch$name #nolint
  assertthat::assert_that(collection %in% col_names,
    msg = "Collection not found"
  )

  if (tolower(method) == "first") {
    batch_sort <- 1
  } else if (tolower(method) == "last") {
    batch_sort <- -1
  }



  aggregation <- paste0(
    '[
    {
      "$match" : {',
      if (!is.null(sirets)){
        paste0('"siret" : {
          "$in" : [',
    paste0(paste0('"', sirets, '"'), collapse = ", "), '
            ]
        },')} else {''},
      if (!is.null(batchs)){
        paste0('"batch" : {
          "$in" : [',
    paste0(paste0('"', batchs, '"'), collapse = ", "), '
            ]
        },')} else {''},
        '"algo" : "', algo, '"
      }
    },
    {
      "$sort" : {
        "siret" : 1.0,
        "batch" : ', batch_sort, ',
        "periode" : 1.0,
        "timestamp" : -1.0
      }
    },
    {
      "$group" : {
        "_id" : {
          "siret" : "$siret",
          "periode" : "$periode"
        },
        "score" : {
          "$first" : "$score"
        },
        "batch" : {
          "$first" : "$batch"
        },
        "algo" : {
          "$first" : "$algo"
        },
        "timestamp" : {
          "$first" : "$timestamp"
        }
      }
    },
    {
      "$addFields" : {
        "siret" : "$_id.siret",
        "periode" : "$_id.periode"
      }
    },
    {
      "$project" : {
        "_id" : 0.0
      }
    }
    ]'
  )

  query <- db$aggregate(
    aggregation
  )

  if (nrow(query) >= 1) result <- query

  return(result)
}

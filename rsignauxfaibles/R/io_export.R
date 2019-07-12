#' Prepare_for_export
#'
#' Si le nombre de périodes disponible dans les données est supérieur à deux,
#' alors la dernière période disponible est absente du data.frame de retour.
#'
#' @param data_to_export `data.frame(1)` \cr Données à exporter. Doit avoir
#'   les champs "siret", "periode", et "score".
#' @param export_fields `character()` \cr Liste des champs à charger pour
#'   enrichir les scores. Ces champs sont chargés directement depuis la base
#'   avec les informations de connection.
#' @inheritParams mongodb_connection
#' @param last_batch `character(1)` \cr Nom du batch à prendre en compte pour
#'   les données.
#' @param known_sirens_full_paths `character()` \cr Chemins absolus des
#'   fichiers contenant des listes de sirens connus.
#' @param verbose `logical()` \cr Faut-il imprimer des informations sur le
#'   chargement des données supplémentaires ?
#'
#' @return `data.frame()`\cr
#'   Données formatées avec les champs "siret", "periode", "score",
#'   "score_diff", "connu" ainsi que les champs spécifiés dans
#'   `export_fields`.
#' @export
format_for_export <- function(
  data_to_export,
  export_fields,
  database,
  collection,
  last_batch,
  # TODO: no default needed here
  # ec79ac15-b8b3-426c-b804-f441851ac58b
  known_sirens_full_paths = c(
    rprojroot::find_package_root_file(
      "..",
      "data-raw",
      "sirets_connus_pdl.csv"
    ),
    rprojroot::find_package_root_file(
      "..",
      "data-raw",
      "sirets_connus_bfc.csv"
    )
  ),
  verbose = TRUE) {

  require(logger)
  if (verbose) {
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  prediction <- data_to_export %>%
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
  log_info("Preparation a l'export ... ")
  log_info("Derniere periode connue: {last_period}")

  donnees <- connect_to_database(
    database = database,
    collection = collection,
    batch = last_batch,
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

  # Report des dernieres infos financieres connues

  # TODO couper nom / chemin du dossier
  donnees <- donnees %>%
    mark_known_sirets(names = basename(known_sirens_full_paths),
      absolute_path = unique(dirname(known_sirens_full_paths))) %>%
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
export_scores_to_mongodb <- function(
  formatted_data,
  algo,
  batch,
  f_scores,
  database,
  collection
  ){

  browser()
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
    # TODO: a lire du fichier de config
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
#' @return `TRUE`. \cr Exporte un fichier csv, dans le dossier spécifié
#'   en entrée, nommé automatiquement "aaaa_mm_jj_export_{algo}_{batch}" avec
#'   éventuellement un suffixe "_vX" pour ne pas écraser de fichier existant.
#'
#' @export
#'
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

#' Marks sirets from files as "known"
#'
#' Récupère des sirens d'une liste de fichier et ajoute une colonne à un
#' data.frame pour spécifier dans une colonne "connu" si le siren est présent dans l'un de ces
#' fichiers.
#'
#' @param df `data.frame()` \cr Table avec au moins une colonne "siret"
#' @param names `character()` \cr Nom des fichiers
#' @param absolute_path `character()` \cr Chemin d'accès absolu
#'
#' @return `data.frame()` \cr Table donnée en entrée à laquelle a été apposé
#' une colonne "connu", qui vaut 1 si le siren correspondant a été trouvé dans
#' les fichiers spécifiés, 0 sinon.
#'
#' @export
#'
mark_known_sirets <- function(
  df,
  names,
  absolute_path = rprojroot::find_package_root_file("..", "data-raw")
) {
  sirens <- c()
  for (name in names) {
    sirets <- readLines(
      file.path(absolute_path, name)
    )
    sirens <- c(sirens, substr(sirets, 1, 9))
  }
  df <- df %>%
    dplyr::mutate(connu = as.numeric(substr(siret, 1, 9) %in% sirens))

  return(df)
}

#' Récupération des scores
#'
#' @inheritParams mongodb_connection
#' @param algo`character(1)`\cr Algorithm from which to take the score. Defaults to "algo".
#' @param method `character(1)`\cr Ou bien "first" ou bien "last". See details.
#' @param siret `character()` \cr Vecteur de sirets à requêter
#'
#' @section Details:
#' Si la méthode est "first", alors le dernier score disponible pour le
#' **premier** batch disponible pour cette période est récupéré. Sinon, le
#' dernier score disponible pour le **dernier** batch disponible est récupéré.
#' "first" représente plus fidèlement les scores qui ont été obtenues par le
#' passé, là où "last" représente les scores qui auraient été obtenues avec
#' le dernier algorithme en date (même si ce n'est pas toujours le cas, ça
#' dépend des périodes qui ont été  exportées avec cet algorithme).
#'
#' @return Returns a dataframe with columns "siret", "periode", "score"
#' (failure probability), "batch" (batch of the model used), "algo" (name of
#' the algorithm)
#'
#' @export
#'
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


#' Export d'une fiche visite
#'
#' @param sirets Vecteur de sirets des établissements dont il faut exporter la
#' fiche
#' @inheritParams mongodb_connection
#' @param batch `character(1)` Choix du batch pour lequel se fait l'export
#' @param with_urssaf `logical(1)`\cr Si TRUE, les informations du montant des
#'  dettes URSSAF sont incluses à la fiche.
#' @param absolute_path `character(1)` \cr Chemin d'accès absolu au dossier
#' d'export.
#'
#' @return `NULL`
#'   Exporte la ou les fiches visites sous le nom
#'   "Fiche_visite_{raison_sociale}.pdf"
#' @export
#'
export_fiche_visite <- function(
  sirets,
  database = "test_signauxfaibles",
  collection = "Features",
  batch,
  with_urssaf = FALSE,
  absolute_path = rprojroot::find_package_root_file("..", "output", "Fiches", batch)){

  for (i in seq_along(sirets)) {
    elementary_info <- connect_to_database(
      database = database,
      collection = collection,
      batch = batch,
      siren = substr(sirets[i], 1, 9),
      fields = c("siret", "raison_sociale", "periode", "code_ape_niveau3"),
      min_effectif = 0
    ) %>%
      filter(siret == sirets[i]) %>%
      head(n = 1)

    raison_sociale <- elementary_info %>%
      .$raison_sociale

    require(rprojroot)
    rmarkdown::render(rprojroot::find_package_root_file("R", "post_fiche_visite.Rmd"),
      params = list(
        siret = sirets[i],
        batch = batch,
        database = database,
        collection = collection,
        with_urssaf = with_urssaf
      ),
      output_file =  file.path(absolute_path, paste0("Fiche_visite_", stringr::str_replace_all(raison_sociale, "[^[:alnum:]]", "_"), ".pdf")),
      clean = TRUE
    )
  }
}

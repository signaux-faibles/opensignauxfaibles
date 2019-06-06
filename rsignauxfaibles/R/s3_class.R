sf_task <- function(
  verbose,
  database,
  collection,
  experiment_aim,
  experiment_description
  ){
  res <- list(database = database, collection = collection)
  class(res) <- "sf_task"
  attr(res, "verbose") <- verbose
  attr(res, "to_log") <- list(
    experiment_aim = experiment_aim,
    experiment_description = experiment_description)
  return(res)
}

# Generic functions
load_hist_data <- function(task, ...){
  UseMethod("load_hist_data", task)
}
load_new_data <- function(task, ...){
  UseMethod("load_new_data", task)
}
hold_out <- function(task, ...){
  UseMethod("hold_out", task)
}
prepare <- function(task, ...){
  UseMethod("prepare", task)
}
train <- function(task, ...){
  UseMethod("train", task)
}
load <- function(task, ...){
  UseMethod("load", task)
}
load.default <- function(task, ...){
  base::load(task, ...)
}
save <- function(task, ...){
  UseMethod("save", task)
}
save.default <- function(task, ...){
  base::save(task, ...)
}
export <- function(task, ...){
  UseMethod("export", task)
}
evaluate <- function(task, ...){
  UseMethod("evaluate", task)
}
log <- function(task, ...){
  UseMethod("log", task)
}
explain <- function(task, ...){
  UseMethod("explain", task)
}


#' Chargement de données historiques
#'
#' Charge les données historiques de signaux faibles et les stocke dans un
#' champ "hist_data" de l'objet \code{task}
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param batch `character(1)` \cr Batch auquel doit être importées les
#'   données. Les modifications opérées par les batchs ultérieurs sont
#'   ignorées.
#' @param database `character(1)` \cr Nom de la base de données vers laquelle
#'   param exporter. Par défaut, celle stockée dans \code{task}.
#' @param collection `character(1)` \cr Nom de la collection vers laquelle
#'   exporter. Par défaut, celle stockée dans \code{task}.
#' @param subsample `integer(1)` \cr Nombre d'objets (c'est-à-dire de couples
#'   siret x periode) à échantillonner.
#' @param fields `character()` \cr Noms des champs à requêter dans la base de
#'   données. Doit contenir "siret" et "periode". Si égal à \code{NULL}, alors
#'   charge tous les champs disponibles.
#' @param min_effectif `integer(1)` \cr Limite basse du filtrage de l'effectif
#'   (la limite est incluse)
#' @param siren `character()` \cr Liste de sirens à exporter. Si égale à
#'   \code{NULL}, charge tous les sirens disponibles.
#' @param code_ape `character() \cr Liste de codes APE à exporter. Si égale à
#'   \code{NULL}, charge tous les codes disponibles.
#'
#' @return `[sf_task]` \cr
#'   L'objet \code{task} donné en entrée auquel le champs "hist_data" a été
#'   ajouté (ou écrasé), contenant un  data.frame() avec les colonnes incluses dans le paramètre d'entrée
#'  \code{fields}, et pour chaque ligne un couple unique siret x periode.
#'
#' @export
#'
#' @examples
load_hist_data.sf_task <- function(
  task,
  batch,
  database = task[["database"]],
  collection = task[["collection"]],
  subsample = 200000L,
  fields = get_fields(training = FALSE),
  date_inf = as.Date("2015-01-01"),
  date_sup = as.Date("2017-01-01"),
  min_effectif = 10L,
  siren = NULL,
  code_ape = NULL
  ){

  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  log_info("Loading historical data")

  hist_data <- connect_to_database(
    database,
    collection,
    batch,
    siren = NULL,
    date_inf = date_inf,
    date_sup = date_sup,
    min_effectif = min_effectif,
    fields = fields,
    code_ape = NULL,
    type = "dataframe",
    subsample = subsample,
    verbose = attr(task, "verbose")
    )

  if (nrow(hist_data) > 1) {
    log_info("Data has been loaded successfully")
  } else {
    log_warn("No data has been loaded, no data could be found")
  }
  task[["hist_data"]] <- hist_data
  return(task)
}

#' Chargement de nouvelles données
#'
#' @param task
#'
#'  @return
#'  @export
#'
#'  @examples
load_new_data.sf_task <- function(
  task,
  last_batch,
  periods,
  database = task[["database"]],
  collection = task[["collection"]],
  fields = get_fields(training = FALSE),
  min_effectif = 10,
  rollback_months = 1
  ){

  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  log_info("Loading data from last batch")
  task[["new_data"]] <- get_last_batch(
    database = database,
    collection = collection,
    last_batch = last_batch,
    periods = periods,
    fields = fields,
    min_effectif = min_effectif,
    rollback_months =  rollback_months
    )

  return(task)
}

hold_out.sf_task <- function(
  task,
  frac_train = 0.6,
  frac_val = 0.2
  ){
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }



  log_info("Historical data is splitted into a train, validation and test
    frame")

    assertthat::assert_that("hist_data" %in% names(task),
      msg = "Please load historical data before holding out test data")

    if (frac_train == 1) {
      frac_val <- 0
      task[["train_data"]] <- task[["hist_data"]]
    } else {
      res <- split_snapshot_rdm_month(
        task[["hist_data"]],
        frac_train,
        frac_val
        )
      task[["train_data"]] <- task[["hist_data"]] %>%
        semi_join(res[["train"]], by = c("siret", "periode"))
      task[["validation_data"]] <- task[["hist_data"]] %>%
        semi_join(res[["validation"]], by = c("siret", "periode"))
      task[["test_data"]] <- task[["hist_data"]] %>%
        semi_join(res[["test"]], by = c("siret", "periode"))
    }


    attr(task, "to_log")[["resampling_strategy"]]  <-  "holdout"
    attr(task, "to_log")[["train_val_test_shares"]] <- c(
          frac_train,
          frac_val,
          1 - frac_train - frac_val
          )
    attr(task, "to_log")[["test_frame"]] <- task[["validation_data"]]
    attr(task, "to_log")[["train_frame"]] <- task[["train_data"]]

    return(task)
}

prepare.sf_task <- function(
  task,
  data_names = c(
    "train_data",
    "validation_data",
    "test_data",
    "new_data"
    )
  ){
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  aux_function <- function(task, name){
    if (name == "train_data"){
      test_or_train <- "train"
    } else {
      test_or_train <- "test"
      assertthat::assert_that("preparation_map" %in% names(task),
        msg = 'No preparation map has been found. Have you loaded a task, or
        prepared the "train_data" first ?')
    }

    log_info("Preparing data for training and predicting")
    if (name %in% names(task)){

      out <-  prepare_frame(
        data_to_prepare = task[[name]],
        save_or_load_map = FALSE,
        te_map = task[["preparation_map"]],
        test_or_train = test_or_train,
        outcome = "outcome"
        )

      task[[paste0("prepared_", name)]]  <- out[["data"]]
      if (name == "train_data"){
        task[["preparation_map"]]  <- out[["te_map"]]
      }

      return(task)

    } else {
      log_warn("There is no {name} in current task")
      return(NULL)
    }
  }

  for (name in data_names){
    task  <- aux_function(task, name)
  }

  attr(
    task,
    "to_log"
  )[["preprocessing_strategy"]] <- "H2OFrame, Target encoding with H2O"

  return(task)
}

train.sf_task <- function(
  task,
  fields = get_fields(training = TRUE)
  ){
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  log_info("Model is being trained.")

  task[["features"]] <- fields

  model <- train_light_gradient_boosting(
    h2o_train_data = task[["prepared_train_data"]],
    x_fields_model = fields,
    save_results = FALSE
    )

  log_info("Model trained_successfully")
  task[["model"]] <- model

  attr(task, "to_log")[["model_name"]] <- "light gradient boosting"
  attr(task, "to_log")[["model"]] <- "model"
  attr(task, "to_log")[["model_features"]] <- fields
  attr(task, "to_log")[["model_parameters"]] <- list(
        "learn_rate <- 0.1",
        "max_depth <- 4",
        "ntrees <- 60",
        "seed <- 123"
        )
  attr(task, "to_log")[["model_target"]] <- "18 mois, defaut et défaillance"

  return(task)
}

load.sf_task <- function(task){
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }
  task[["model"]] <- load_h2o_object("lgb", "model", last = TRUE)
  task[["preparation_map"]] <- load_h2o_object("te_map", "temap", last = TRUE)

  return(task)
}

save.sf_task <- function(task, ...) {
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  assertthat::assert_that(all(c("model", "preparation_map") %in% names(task)),
    msg = "Task should have a model and a preparation_map to be saved")
  save_h2o_object(task[["model"]], "lgb")
  save_h2o_object(task[["preparation_map"]], "preparation_map")

  return(task)
}

#' Predict from a trained model
#'
#' Predict on some data.
#' Data should be prepared.
#'
#' @param object
#'
#' @return
#' @export
#'
#' @examples
predict.sf_task <- function(
  object,
  data_names = c("new_data", "train_data", "validation_data", "test_data")
){

  task  <- object
  require(logger)
  if (attr(task, "verbose")){
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  assertthat::assert_that(all(c("model") %in% names(task)),
    msg = "Task should have a model to predict on new
    data")

  predict_on_given_data <- function(data_name, task){

    prepared_data_name <- paste0("prepared_", data_name)
    if (!prepared_data_name %in% names(task)){
      log_warn("{data_name} is missing or has not been prepared yet")
      return(task)
    }

    log_info("Model is being applied on {prepared_data_name}")

    prediction <- predict_model(
      model = task[["model"]],
      new_data = task[[prepared_data_name]]
      )
    task[[data_name]] <- task[[data_name]] %>%
      left_join(prediction, by = c("siret", "periode"))
    return(task)
  }

  for (name in data_names){
    task  <- predict_on_given_data(name, task)
  }

  log_info("Prediction successfully done.")
  return(task)
}

export.sf_task <- function(task, ...){
  export_fields <- c(
    "siret",
    "periode",
    "raison_sociale",
    "departement",
    "region",
    "score",
    "score_diff",
    "connu",
    "date_ccsf",
    "etat_proc_collective",
    "date_proc_collective",
    "interessante_urssaf",
    # "default_urssaf",
    "effectif",
    "libelle_naf",
    "libelle_ape5",
    "code_ape",
    "montant_part_ouvriere",
    "montant_part_patronale",
    "ca",
    "ca_past_1",
    "benefice_ou_perte",
    "benefice_ou_perte_past_1",
    "resultat_expl",
    "resultat_expl_past_1",
    "poids_frng",
    "taux_marge",
    "frais_financier",
    "financier_court_terme",
    "delai_fournisseur",
    "dette_fiscale",
    "apart_heures_consommees",
    "apart_heures_autorisees",
    # "cotisation_moy12m",
    "compte_urssaf",
    "montant_majorations",
    "exercice_bdf",
    "exercice_diane",
    "delai"
    )

  f_scores <- c(F1 = 0.31, F2 = 0.13) # TODO TODO

  if (!is.null(export_type) && export_type != "none") {
    assertthat::assert_that(all(export_type %in% c("csv", "mongodb")))

    log_info("Adding additional fields for export")
    res <- task[["new_data"]] %>%
      format_for_export(
        export_fields = export_fields,
        database = database,
        collection = collection,
        last_batch = last_batch
        )

    log_info("Data is exported to {paste(export_type, collapse = ' and ')}")
    purrr::walk(
      .x = export_type,
      .f = function(x, ...){
        if (x == "csv") {
          export_scores_to_csv(
            ...,
            relative_path = "../output/"
            )
        } else if (x == "mongodb") {
          export_scores_to_mongodb(
            ...,
            f_scores = f_scores,
            database = database,
            collection = collection
            )
      }},
      formatted_data = res,
      batch = last_batch,
      algo = "algo"
      )
  }
  log_info("Data exported with success to
    {paste(export_type, collapse = ' and ')}")
    return(task)
}

evaluate.sf_task <- function(
  task,
  data_name = c("validataion_data")
){

  browser()

  assertthat::assert_that(
    length(data_name) == 1,
    msg = "Evaluation can only be made on a single data.frame at once"
  )

  require(MLsegmentr)
  eval <- Assesser$new(
    task[[data_name]] %>%
      filter(periode == max(periode))
    )

  eval$set_predictions("score")
  eval$set_targets("outcome")

  eval$evaluation_funs <- eval_precision_recall()

  perf <- eval$assess_model()

  attr(task, "to_log")[["model_performance"]] <- perf %>%
        select(evaluation_name, evaluation) %>%
        filter(evaluation_name != "prcurve")

  return(task)
  ## Log model performance.
}


log.sf_task <- function(
  task,
  database = task[["database"]],
  collection = "ml_logs",
  ...
  ){

  require("MLlogr")
  logger <- MLLogger$new(database, collection)
  do.call(logger$set, args = attr(task, "to_log"))
  browser()
  logger$log(...)

  return(task)
}

explain.sf_task <- function(task, ...){
  # Later
}


print.sf_task <- function(x, ...){
  cat("-- FIELDS --\n")
  aux_fun <- function(name, x){
    if (!is.character(x) || length(x) > 1){
      cat(paste0("  * ", name, " (", paste0(class(x), collapse = ", "), ")\n"))
    } else {
      cat(paste0("  * ", name, " : ", x, "\n"))
    }
  }
  purrr::walk2(names(x), x, aux_fun)

  cat("-- INFO --\n")
  purrr::walk2(
    names(attr(x, "to_log")),
    attr(x, "to_log"),
    aux_fun
  )
  return()
}

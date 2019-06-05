sf_task <- function(
  verbose,
  database,
  collection
){
  res <- list(database = database, collection = collection)
  class(res) <- "sf_task"
  attr(res, "verbose") <- verbose
  require(MLlogr)
  attr(res, "MLLogger") <- MLLogger$new(database, "ml_logs")
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


# Definitions
load_hist_data.sf_task <- function(
  task,
  last_batch,
  database = task[["database"]],
  collection = task[["collection"]],
  subsample = 200000,
  fields = get_fields(training = FALSE),
  date_inf = as.Date("2015-01-01"),
  date_sup = as.Date("2017-01-01"),
  min_effectif = 10,
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
    last_batch,
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

    attr(task, "MLLogger")$set(
      resampling_strategy = "holdout",
      train_val_test_shares = c(
        frac_train,
        frac_val,
        1 - frac_train - frac_val
      ),
      test_frame = task[["test_data"]]
      )

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
        test_or_train = test_or_train
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

  attr(task, "MLLogger")$set(
    preprocessing_strategy = "H2OFrame, Target encoding with H2O"
    )
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

  attr(task, "MLLogger")$set(
    model_name = "light gradient boosting",
    model = "model",
    model_features = fields,
    model_parameters = list(
      "learn_rate = 0.1",
      "max_depth = 4",
      "ntrees = 60",
      "seed = 123"
    ),
    model_target = "18 mois, defaut et défaillance"
    )
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
  save_h2o_object(taks[["te_map"]], "te_map")

  return(task)
}

#' Predict from a trained model
#'
#' Predict on some data.
#' Data should be prepared.
#'
#' @param
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
      log_warning("{data_name} is missing or has not been prepared yet")
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



    #task[["prediction"]] <- pred_data
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
    "delai",
    "apparait",
    "disparait"
    )

  F_scores <- c(F1 = 0.31, F2 = 0.13)


  if (!is.null(export_type) && export_type != "none") {
    assertthat::assert_that(all(export_type %in% c("csv", "mongodb")))

    log_info("Adding additional fields for export")
    res <- pred_data %>%
      prepare_for_export(
        export_fields = export_fields,
        database = database,
        collection = collection,
        last_batch = last_batch
        )
    log_info("Data is exported to {paste(export_type, collapse = ' and ')}")
    purrr::map(
      .x = export_type, function(x, ...) export(destination = x, ...),
      donnees = res,
      batch = last_batch,
      F_scores = F_scores
      )
  }
  return(task)
}

evaluate.sf_task <- function(task){

  require(MLsegmentr)
  eval <- Assesser$new(
    task[["new_data"]] %>%
      as_tibble() %>%
      filter(periode == max(periode)) %>%
      mutate(outcome = sample(
        c(TRUE, FALSE),
        length(outcome),
        replace = TRUE
      )) %>%
      mutate(prediction = as.numeric(as.vector(task[["prediction"]]$prob)))
      )

  eval$set_predictions(
    MLsegmentr::add_id(
      data.frame(prediction = as.numeric(as.vector(task[["prediction"]]$prob)))
    )
  )
  eval$set_targets("outcome")

  eval$evaluation_funs <- eval_precision_recall()

  perf <- eval$assess_model()

  attr(task, "MLLogger")$set(
    model_performance = perf %>%
      select(evaluation_name, evaluation) %>%
      filter(evaluation_name != "prcurve")
    )

  return(task)
  ## Log model performance.
}


log.sf_task <- function(
  task,
  experiment_aim,
  experiment_description,
  database = "test_signaux_faibles",
  collection = "ml_logs"
  ){
  require("MLlogr")

  attr(task, "MLLogger")$set(
    experiment_aim = experiment_aim,
    experiment_description = experiment_description
    )
  attr(task, "MLLogger")$log()
}

explain.sf_task <- function(task, ...){
  # Later
}

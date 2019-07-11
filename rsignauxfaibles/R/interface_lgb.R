#' Quickly run the model -- temporary function
#'
#' Charge les données, les prépare, entraîne l'algorithme, charge les
#' dernières données, prédit sur les dernières données et exporte les
#' résultats en un appel de fonction.
#'
#' @param sample_size `integer()` \cr Taille de l'échantillon sur lequel
#'   l'algorithme est entraîné.
#' @param  batch `character()` \cr Nom du batch pour lequel le modèle est
#' entraîné et les dernières données sont sélectionnées.
#'
#'  @return `TRUE`
#'  @export
full_light_gradient_boosting <- function(
  sample_size,
  batch
){
  h2o::h2o.no_progress()
  devtools::load_all()
  # Script parameters
  database <- "test_signauxfaibles"
  collection <- "Features"
  actual_period <- as.Date("2019-03-01")
  last_batch <- "1905"
  min_effectif <- 10
  retrain_model <- TRUE
  type <- "dataframe"
  export_type <- c("csv", "mongodb")
  connect_to_h2o()
  task <- sf_task(
    verbose = TRUE,
    database,
    collection,
    experiment_name = "Reference evaluation canvas",
    experiment_description = "This experiment fixes the conditions for the
    evaluation of a new algorithm compared to former versions"
    )

  task <- load_hist_data(
    task = task,
    batch = batch,
    subsample = sample_size
    )
  task <- hold_out(task)
  task <- prepare(task, data_names = c("train_data"))
  task  <- train(task)
  task <- prepare(task, data_names = c("validation_data"))
  task <- predict(task, data_names = c("validation_data"))
  task <- evaluate(task, data_name = "validation_data")
  task <- log(task)

  return(TRUE)
}

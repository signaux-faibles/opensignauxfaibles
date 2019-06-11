#' Quickly run the model -- temporary function
#'
#' @param param
#'
#'  @return
#'  @export
#'
#'  @examples
full_light_gradient_boosting <- function(sample_size){
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
    experiment_aim = "Reference evaluation canvas",
    experiment_description = "This experiment fixes the conditions for the
    evaluation of a new algorithm compared to former versions"
    )

  task <- load_hist_data(
    task = task,
    batch = "1905",
    subsample = sample_size
    )
  task <- hold_out(task)
  task <- prepare(task, data_names = c("train_data"))
  task  <- train(task)
  task <- prepare(task, data_names = c("validation_data"))
  task <- predict(task, data_names = c("validation_data"))
  task <- evaluate(task, data_name = "validation_data")
  task <- log(task)
}

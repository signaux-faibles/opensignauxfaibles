#' Train a light gradient boosting model
#'
#' @param train_data données d'entraînement, sous la forme d'un H2OFrame
#' @param validation_data données d'évaluation, sous la forme d'un H2OFrame
#' @param x_fields_model Champs sur lequel entraîner le modèle.
#' @param save_results si TRUE, sauvegarde le modèle sur le disque.
#'
#' @return Liste("model", "temap")
#' @export
#'
#' @examples
train_light_gradient_boosting <- function(
  h2o_train_data,
  h2o_validation_data = NULL,
  x_fields_model,
  save_results = TRUE
){

  #
  # Train the model
  #


  y <- "outcome"

  model <- h2o.xgboost(
    model_id = "Model_train",
    x = x_fields_model,
    y = y,
    training_frame = h2o_train_data,
    validation_frame = h2o_validation_data,
    tree_method = "hist",
    grow_policy = "lossguide",
    learn_rate = 0.1,
    max_depth = 4,
    ntrees = 60,
    seed = 123
    )

  if (save_results){
    save_h2o_object(model, "lgb") #lgb like Light Gradient Boosting
  }
  return(model)
}

#' Préparation des données pour l'entraînement
#'
#' Inclut la conversion en H2O, et le target encoding.
#' La map pour le target encoding peut être sauvegardé, afin de permettre de
#' l'appliquer à de nouvelles données.
#'
#' @param train_data
#' @param validation_data
#' @param save_results
#'
#' @return
#' @export
#'
#' @examples
prepare_frame_light_gradient_boosting <- function(
  train_data,
  validation_data = NULL,
  save_results = TRUE
) {
  browser()
  with_validation <- !is.null(validation_data)

  h2o_train_data <- convert_to_h2o(train_data)
  if (with_validation){
    h2o_validation_data <- convert_to_h2o(validation_data)
  }

  #
  # Target Encoding de differents groupes sectoriels
  #

  te_map <- h2o.target_encode_create(
    h2o_train_data,
    x = list(c("code_naf"),
             c("code_ape_niveau2"),
             c("code_ape_niveau3")),
    y = "outcome")

  h2o_train_data <- h2o_target_encode(
    te_map,
    h2o_train_data,
    "train")

  if (!is.null(validation_data)){
    h2o_validation_data <- h2o_target_encode(
      te_map,
      h2o_validation_data,
      "test")
  }


  if (save_results){
    save_h2o_object(te_map, "te_map")
  }

  res <- list(
    train_data = h2o_train_data,
    te_map = te_map)
  if (!is.null(validation_data)){
    res[["validation_data"]] <- h2o_validation_data
  }

  return(res)
}

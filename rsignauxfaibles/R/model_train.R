#' Train a light gradient boosting model
#'
#' Trains a light gradient boosting model on training data.
#'
#' @param h2o_train_data `H2OFrame` \cr données d'entraînement, sous la forme d'un H2OFrame
#' @param h2o_validation_data `H2OFrame` \cr Données d'évaluation, sous la forme d'un H2OFrame
#' @param x_fields_model Champs sur lequel entraîner le modèle.
#' @param save_results si TRUE, sauvegarde le modèle sur le disque.
#' @param outcome `character(1)` \cr Nom de la variable qui sert de cible
#'   d'apprentissage
#' @param learn_rate `numeric(1)` \cr
#' @param max_depth `integer(1)` \cr  Specify the maximum tree depth. Higher
#'   values will make the model more complex and can lead to overfitting.
#'   Setting this value to 0 specifies no limit.
#' @param ntrees `integer(1)` \cr Specify the number of trees to build.
#' @param min_child_weight `integer(1)` \cr Specify the minimum number of
#' observations for a leaf.
#' @param seed `integer(1)` \cr Graine aléatoire pour que les opérations
#'   aléatoires soient reproductibles.
#' @param save_results `logical(1)` \cr Faut-il sauvegarder le modèle pour un
#'   chargement ultérieur ?
#'
#' @return `H2OBinomialModel` \cr
#' @export
#'
train_light_gradient_boosting <- function(
  h2o_train_data,
  h2o_validation_data = NULL,
  x_fields_model,
  outcome = "outcome",
  learn_rate = 0.1,
  max_depth = 4,
  ntrees = 60,
  min_child_weight = 1,
  seed = 123,
  save_results = TRUE) {

  #
  # Train the model
  #

  require(uuid)
  model <- h2o::h2o.xgboost(
    model_id = paste0("Train_model_", uuid::UUIDgenerate()),
    x = x_fields_model,
    y = outcome,
    training_frame = h2o_train_data,
    validation_frame = h2o_validation_data,
    tree_method = "hist",
    grow_policy = "lossguide",
    learn_rate = learn_rate,
    max_depth = max_depth,
    ntrees = ntrees,
    min_child_weight = min_child_weight,
    seed = seed
  )

  if (save_results) {
    save_h2o_object(model, "lgb") # lgb like Light Gradient Boosting
  }
  return(model)
}

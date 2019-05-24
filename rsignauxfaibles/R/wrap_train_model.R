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
                                          outcome = "outcome",
                                          save_results = TRUE) {

  #
  # Train the model
  #



  # FIX ME: RISQUE D'ECRASER UN MODELE EXISTANT AVEC CE MODEL_ID

  model <- h2o::h2o.xgboost(
    model_id = paste0("Train_model", floor(runif(1) * 1000)),
    x = x_fields_model,
    y = outcome,
    training_frame = h2o_train_data,
    validation_frame = h2o_validation_data,
    tree_method = "hist",
    grow_policy = "lossguide",
    learn_rate = 0.1,
    max_depth = 4,
    ntrees = 60,
    seed = 123
  )

  if (save_results) {
    save_h2o_object(model, "lgb") # lgb like Light Gradient Boosting
  }
  return(model)
}


AUCPR <- function(y_pred, y_true){

  PR <-  pr.curve(scores.class0 = y_pred,
                  weights.class0 =  as.numeric(y_true),curve = TRUE)

  return(PR$auc.integral)
}


# define custom callback class

#AUCPR_keras_callback <- function(){
#  AUCPR_keras <- R6::R6Class("AUCPR_keras_callback",
#                             inherit = KerasCallback,
#
#                             public = list(
#
#                               AUCPR = NULL,
#
#                               on_epoch_end = function(epoch, logs = list()) {
#                                 X_val = self.validation_data[0]
#                                 Y_val = self.validation_data[1]
#                                 Y_predict = model.predict(X_val)
#                                 self$AUCPR <- c(self$AUCPR,
#                                                 AUCPR(
#                                                   y_true = Y_val,
#                                                   y_pred = Y_predict)
#                                 )
#                                 cat(self$AUCPR)
#                               }
#                             ))
#
#  return(AUCPR$new())
#}

#' Fscore_from_prediction
#'
#' @param prediction
#' @param target
#' @param alpha
#'
#' @return
#' @export
#'
#' @examples
Fscore_from_prediction <- function(prediction, target, alpha){
  pr_object <- PRROC::pr.curve(
    scores.class0 = prediction,
    weights.class0 =  target,
    curve = TRUE)
  res <- as.data.frame(pr_object$curve) %>%
    dplyr::mutate(Fscore = Fscore_from_prcurve(pr_object$curve, alpha)) %>%
    dplyr::filter(Fscore == max(Fscore, na.rm = TRUE)) %>%
    .[1, 3]
  return(res)
}

#' Trace la courbe précision rappel pour le modèle H2O considéré
#'
#' @param model Modèle H2O
#' @param my_data H2OFrame sur lequel la courbe PRROC est tracée (validation, test etc.)
#' @param new_fig Si TRUE, ouvre une nouvelle fenêtre graphique. Sinon, trace dans la dernière fenêtre active.
#'
#' @return
#' @export
#'
#' @examples
plotPR <- function(model, my_data, new_fig = TRUE, model_objective = "outcome"){
  if (new_fig)  plot(1, type = 'n', xlab = "recall", ylab = "precision", xlim = c(0,1), ylim = c(0,1))
  perf <- h2o::h2o.performance(model, newdata = my_data)
  pred <- h2o::h2o.predict(model, my_data)

  true_res <- as.vector(h2o::as.numeric(my_data[model_objective]))



  precision <- h2o::h2o.precision(perf)
  recall <- h2o::h2o.recall(perf)
  lines(recall$tpr, precision$precision, col = rgb(runif(5),runif(5),runif(5)) )

  F2 <- h2o::h2o.F2(perf)
  precision_F2 <- precision$precision[which.max(F2$f2)]
  recall_F2 <- recall$tpr[which.max(F2$f2)]
  points(recall_F2, precision_F2, col = 'red', pch = 4)
  text(recall_F2, precision_F2,'F2', pos = 4, col = 'red')

  F1 <- h2o::h2o.F1(perf)
  precision_F1 <- precision$precision[which.max(F1$f1)]
  recall_F1 <- recall$tpr[which.max(F1$f1)]
  points(recall_F1, precision_F1, col = 'green', pch = 4)
  text(recall_F1, precision_F1,'F1', pos = 4, col = 'green')


  cat(str(PRROC::pr.curve(scores.class0 = as.vector(pred[[3]]), weights.class0 = true_res)))
}


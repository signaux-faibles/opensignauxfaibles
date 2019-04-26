AUCPR <- function(y_pred, y_true){

  PR <-  pr.curve(
    scores.class0 = y_pred,
    weights.class0 = as.numeric(y_true),
    curve = TRUE
  )

  return(PR$auc.integral)
}


#' Fonction utilitaire pour calculer le Fscore_threshold
#'
#' @param predicted
#' @param outcome
#' @param alpha Paramètre du Fscore_threshold (alpha = 1 pour le plus courant, le
#' F1-score)
#'
#' @return
#' @export
#'
#' @examples
Fscore_threshold <- function(prediction, outcome, alpha){
  pr_object <- PRROC::pr.curve(
    scores.class0 = prediction,
    weights.class0 = outcome,
    #as.numeric(outcome),
    curve = TRUE
  )

   res <- as.data.frame(pr_object$curve) %>%
    dplyr::mutate(
      Fscore_threshold = MLsegmentr::Fscore_from_prcurve(pr_object$curve, alpha)
    ) %>%
    dplyr::filter(Fscore_threshold == max(Fscore_threshold, na.rm = TRUE)) %>%
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
plotpr <- function(
  model,
  my_data,
  new_fig = TRUE,
  model_objective = "outcome"
){

  if (new_fig)  plot(
    1,
    type = "n",
    xlab = "recall",
    ylab = "precision",
    xlim = c(0, 1),
    ylim = c(0, 1)
  )
  perf <- h2o::h2o.performance(model, newdata = my_data)
  pred <- h2o::h2o.predict(model, my_data)

  true_res <- as.vector(h2o::as.numeric(my_data[model_objective]))



  precision <- h2o::h2o.precision(perf)
  recall <- h2o::h2o.recall(perf)
  lines(recall$tpr, precision$precision, col = rgb(
    runif(5),
    runif(5),
    runif(5)
  ))

  F2 <- h2o::h2o.F2(perf)
  precision_F2 <- precision$precision[which.max(F2$f2)]
  recall_F2 <- recall$tpr[which.max(F2$f2)]
  points(recall_F2, precision_F2, col = "red", pch = 4)
  text(recall_F2, precision_F2, "F2", pos = 4, col = "red")

  F1 <- h2o::h2o.F1(perf)
  precision_F1 <- precision$precision[which.max(F1$f1)]
  recall_F1 <- recall$tpr[which.max(F1$f1)]
  points(recall_F1, precision_F1, col = "green", pch = 4)
  text(recall_F1, precision_F1, "F1", pos = 4, col = "green")


  cat(str(PRROC::pr.curve(scores.class0 = as.vector(
    pred[[3]]
  ), weights.class0 = true_res)))
}



#' Fonction d'évaluation de survenue de nouvelles dettes urssaf.
#'
#' @param eval_frame
#' @param additional_data In additional_data we need to have: debit_18_mois,
#' interet URSSAF, siret
#' @return
#' @export
#'
#' @examples
custom_eval_urssaf <- function(eval_frame, additional_data) {

  eval_frame <- dplyr::left_join(eval_frame, additional_data, by = ".id")

  assertthat::assert_that(all(
      c("siret", "periode") %in% names(additional_data)
      ), msg = "Données manquantes")


  F1_thres <- Fscore_threshold(
    prediction = eval_frame$prediction,
    outcome = eval_frame$outcome,
    alpha = 1
    )
  F2_thres <- Fscore_threshold(
    prediction = eval_frame$prediction,
    outcome = eval_frame$outcome,
    alpha = 2
    )
  aux_rates <- function(eval) {

    assertthat::assert_that(F2_thres <= F1_thres)
    first_detection  <- eval %>%
      dplyr::group_by(siret) %>%
      dplyr::arrange(siret, periode) %>%
      dplyr::mutate(
        any_above_F2 = any(prediction >= F2_thres),
        rank = rank(x = periode, ties.method = "first"),
        kept_rank = ifelse(!any_above_F2, 1,
          which.max(prediction >= F2_thres))
        ) %>%
      dplyr::filter(rank == kept_rank) %>%
      dplyr::select(-c(any_above_F2, kept_rank, rank)) %>%
      dplyr::mutate(etat = .bincode(
        # cut( # Does not handle non unique breaks !
          x = prediction,
          breaks = c(-1e-4, F2_thres, F1_thres, 1 + 1e-4),
          )
        ) %>%
      ungroup()

      res <- first_detection %>%
        dplyr::group_by(etat) %>%
        dplyr::summarise(
          conf = list(
            dplyr::as_data_frame(
            Hmisc::binconf(
              x = sum(as.numeric(outcome == 1)),
              n = dplyr::n(),
              alpha = 0.05
              )
            )
            ),
          count = dplyr::n()
          ) %>%
      tidyr::unnest() %>%
      dplyr::rename(
        moyenne = PointEst,
        lconf = Lower,
        hconf = Upper
        )
      return(res)
  }
  aux_eval_fun <- MLsegmentr::wrap_simple_eval(aux_rates)
  res  <- aux_eval_fun(eval_frame) %>%
    mutate(etat = factor(etat,
        levels = 1:3,
        labels = c("Pas d'alerte", "Alerte seuil F2", "Alerte seuil F1")
        ))
  return(res)
}

custom_plot_urssaf  <- function(evaluation){
  require(ggplot2)
  ggplot2::ggplot(
   evaluation,
   ggplot2::aes(
     x = etat,
     y = moyenne,
     ymin = lconf,
     ymax = hconf,
     fill = outcome_type,
     group = outcome_type,
     label = count
   )) +
     ggplot2::geom_bar(stat = "identity", position = "dodge") +
     ggplot2::geom_errorbar(position = "dodge") +
     ggplot2::facet_wrap(ggplot2::vars(segment, model)) +
     ggplot2::geom_text(y = -0.015)
}

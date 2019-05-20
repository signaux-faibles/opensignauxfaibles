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


## Depreciée Utiliser MLsegmentr avec evaluation_precision_rappel à la
## place.
##' Trace la courbe précision rappel pour le modèle H2O considéré
##'
##' @param model Modèle H2O
##' @param my_data H2OFrame sur lequel la courbe PRROC est tracée (validation, test etc.)
##' @param new_fig Si TRUE, ouvre une nouvelle fenêtre graphique. Sinon, trace dans la dernière fenêtre active.
##'
##' @return
##' @export
##'
##' @examples
#plotpr <- function(
#  model,
#  my_data,
#  new_fig = TRUE,
#  model_objective = "outcome"
#){
#
#  if (new_fig)  plot(
#    1,
#    type = "n",
#    xlab = "recall",
#    ylab = "precision",
#    xlim = c(0, 1),
#    ylim = c(0, 1)
#  )
#  perf <- h2o::h2o.performance(model, newdata = my_data)
#  pred <- h2o::h2o.predict(model, my_data)
#
#  true_res <- as.vector(h2o::as.numeric(my_data[model_objective]))
#
#
#
#  precision <- h2o::h2o.precision(perf)
#  recall <- h2o::h2o.recall(perf)
#  lines(recall$tpr, precision$precision, col = rgb(
#    runif(5),
#    runif(5),
#    runif(5)
#  ))
#
#  F2 <- h2o::h2o.F2(perf)
#  precision_F2 <- precision$precision[which.max(F2$f2)]
#  recall_F2 <- recall$tpr[which.max(F2$f2)]
#  points(recall_F2, precision_F2, col = "red", pch = 4)
#  text(recall_F2, precision_F2, "F2", pos = 4, col = "red")
#
#  F1 <- h2o::h2o.F1(perf)
#  precision_F1 <- precision$precision[which.max(F1$f1)]
#  recall_F1 <- recall$tpr[which.max(F1$f1)]
#  points(recall_F1, precision_F1, col = "green", pch = 4)
#  text(recall_F1, precision_F1, "F1", pos = 4, col = "green")
#
#
#  cat(str(PRROC::pr.curve(scores.class0 = as.vector(
#    pred[[3]]
#  ), weights.class0 = true_res)))
#}



#' Fonction d'évaluation de survenue de nouvelles dettes urssaf.
#'
#' Cette fonction est prévue pour construire un objet EvaluationFunctions
#' afin d'être utilisé dans un objet Assesser du package MLsegmentr. À
#' utiliser avec les options par défaut wrap_eval_fun = TRUE et unnest = TRUE.
#'
#' @param eval_frame
#' @param additional_data In additional_data we need to have: debit_18_mois,
#' interet URSSAF, siret, periode
#' @param F1_thres, F2_thres Seuils qui maximisent les scores F1 et F2
#' respectivement, calculés sur les données d'entraînement.
#' (Possibilité d'utiliser la fonction Fscore_threshold)
#'
#' @return data.frame
#' @export
#'
#' @examples
aux_custom_eval_urssaf <- function(eval_frame, additional_data, F1_thres, F2_thres) {

 ### TODO change
  F1_thres = 0.31
  F2_thres = 0.13

  # TODO simplifier: F1 et F2 sont en fait une segmentation ! Donc à faire
  # sur les données a priori. Cette fonction ne calcule au final que la
  # précision avec intervalle de confiance + count.

  assertthat::assert_that(F2_thres <= F1_thres)

  # eval_frame <- dplyr::left_join(eval_frame, additional_data, by = ".id")

  # assertthat::assert_that(all(
  #     c("siret", "periode") %in% names(additional_data)
  #     ), msg = "Données manquantes")

  # À calculer a priori sur les données d'entraînement;
  #F1_thres <- Fscore_threshold(
  #  prediction = eval_frame$prediction,
  #  outcome = eval_frame$outcome,
  #  alpha = 1
  #  )
  #F2_thres <- Fscore_threshold(
  #  prediction = eval_frame$prediction,
  #  outcome = eval_frame$outcome,
  #  alpha = 2
  #  )

  # Calcul de la première détection au seuils F2 et F1
    first_detection  <- eval_frame %>%
      dplyr::group_by(siret) %>%
      dplyr::arrange(siret, periode) %>%
      dplyr::filter(any(prediction >= F2_thres)) %>%
      dplyr::mutate(
        rank = rank(x = periode, ties.method = "first")
        ) %>%
      dplyr::filter(
        rank == which.max(prediction >= F2_thres) |
        rank == which.max(prediction >= F1_thres)) %>%
      select(-rank) %>%
      dplyr::ungroup()


    # On prend une période au hasard pour les entreprises non détectées.
    not_detected <- eval_frame %>%
      dplyr::group_by(siret) %>%
      dplyr::arrange(siret, periode) %>%
      dplyr::filter(all(prediction < F2_thres)) %>%
      dplyr::mutate(
        rank = rank(x = periode, ties.method = "first"),
        random_rank = ceiling(runif(1, 0, max(rank)))
        ) %>%
      dplyr::filter(rank == random_rank) %>%
      select(-rank, -random_rank) %>%
      dplyr::ungroup()

    # On fusionne ces deux df et on indique le statut (Non détectée, détectée
    # seuil F2, détecté seuil F1)
    all_data <- dplyr::bind_rows(first_detection, not_detected) %>%
      dplyr::mutate(statut = .bincode(
          x = prediction,
          breaks = c(-1e-4, F2_thres, F1_thres, 1 + 1e-4),
          ))  %>%
      dplyr::mutate(statut = factor(statut,
          levels = 1:3,
          labels = c("Pas d'alerte", "Alerte seuil F2", "Alerte seuil F1")
          ))


      # Taux de défaillance moyen avec intervalles de confiance.
      res <- all_data %>%
        dplyr::group_by(statut) %>%
        dplyr::summarise(
          conf = list(
            tibble::as_tibble(
            Hmisc::binconf(
              x = sum(as.numeric(target == 1)),
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

#' Fonction graphique pour la précision
#'
#' Cette fonction est prévue pour construire un objet EvaluationFunctions
#' afin d'être utilisé dans un objet Assesser du package MLsegmentr. À
#' utiliser avec les options par défaut wrap_eval_fun = TRUE et unnest = TRUE.
#'
#' Elle s'utilise en complément de aux_custom_eval_urssaf, donc elle prend le
#' résultat en variable d'entrée.
#'
#' @param eval_frame, tel que donné par la fonction aux_custom_eval_urssaf
#'  avec paramètres wrap_eval_fun et unnest = TRUE.
#'
#' @return NULL
#' @export
#'
#' @examples
aux_custom_plot_urssaf  <- function(evaluation){
  if (require(ggplot2)) {
  ggplot2::ggplot(
   evaluation,
   ggplot2::aes(
     x = statut,
     y = moyenne,
     ymin = lconf,
     ymax = hconf,
     fill = target_type,
     group = target_type,
     label = count
   )) +
     ggplot2::geom_bar(stat = "identity", position = "dodge") +
     ggplot2::geom_errorbar(position = "dodge") +
     ggplot2::facet_wrap(ggplot2::vars(segment, model)) +
     ggplot2::geom_text(y = -0.015)
  } else {
    warning("Cette fonction nécessite le package ggplot2; veuillez l'installer pour l'utiliser")
  }
}

#' Fonctions d'évaluation de la précision (spéciale URSSAF)
#'
#' Cette fonction construit un objet EvaluationFunction avec les fonctions
#' auxiliaires ci-dessus.
#'
#' Etapes pour utiliser cette objet:
#' 1- Entraînement du modèle. Déterminer les seuils F1 et F2 à partir de
#' l'échantillon d'entraînement.
#' 2- faire un objet assesseur, avec cet objet EvaluationFunction,
#' 3 -définir  comme segment la colonne interessant urssaf (6 segments)
#' 4- Définir comme target l'anticipation d'un débit à 18 mois
#'
#'
#' @return Objet R6 EvaluationFunction
#' @export
#'
#' @examples
eval_urssaf  <- function(){
  require(MLsegmentr)
  evo <- EvaluationFunction$new()
  evo$eval_fun <- aux_custom_eval_urssaf
  evo$wrap_eval_fun  <- TRUE
  evo$unnest <- TRUE
  evo$plot_fun  <- aux_custom_eval_urssaf

  return(evo)

}

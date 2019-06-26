#' Fonction d'évaluation de survenue de nouvelles dettes urssaf.
#'
#' Cette fonction est prévue pour construire un objet de type
#' `MLsegmentr::eval_function`, afin d'être utilisé dans un objet
#' `MLsegmentr::Assesser`. À utiliser avec les options par défaut
#' wrap_eval_fun = TRUE et unnest = TRUE.
#'
#' @param eval_frame `data.frame()` \cr Un data.frame à évaluer comme produit par
#'   l'objet `MLsegmentr::Assesser`. Se référer à la documentation
#'   TODO: référence doc MLsegmentr 5fb65cac-b478-4413-a359-294e647103e3
#' @param F1_thres Seuil qui maximise le score F1, calculés sur
#' les données d'entraînement.
#' @param F2_thres Seuil qui maximise le score F2.
#' @return `data.frame` \cr avec en plus des colonnes classiques "model",
#' "target_type", "segment"; les colonnes "statut" (niveau d'alerte), "count"
#' (nombre d'entreprises); "moyenne" (taux de défaillance moyen), "lconf" et
#' "hconf" (limites basse et haute respectivement de l'intervalle de confiance
#' issue d'une loi binomiale, à 95%)
#' @export
aux_custom_eval_urssaf <- function(eval_frame, F1_thres, F2_thres) {

  # TODO: calculer les F-scores !
  # fa85ece1-26cb-4a88-b12b-e846ef16a60f
  # TODO simplifier: F1 et F2 sont en fait une segmentation ! Donc a faire
  # sur les donnees a priori. Cette fonction ne calcule au final que la
  # precision avec intervalle de confiance + count.
  # La fusion des segments doit se faire avant l'application de la fonction
  # d'evaluation !
  F1_thres <- 0.31
  F2_thres <- 0.13

  assertthat::assert_that(F2_thres <= F1_thres)

  # Calcul de la premiere detection au seuils F2 et F1
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


  # On prend une periode au hasard pour les entreprises non detectees.
  not_detected <- eval_frame %>%
    dplyr::group_by(siret) %>%
    dplyr::arrange(siret, periode) %>%
    dplyr::filter(all(prediction < F2_thres)) %>%
    dplyr::mutate(
      rank = rank(x = periode, ties.method = "first"),
      random_rank = ceiling(stats::runif(1, 0, max(rank)))
      ) %>%
    dplyr::filter(rank == random_rank) %>%
    dplyr::select(-rank, -random_rank) %>%
    dplyr::ungroup()

  # On fusionne ces deux df et on indique le statut (Non detectee, detectee
  # seuil F2, detecte seuil F1)
  all_data <- dplyr::bind_rows(first_detection, not_detected) %>%
    dplyr::mutate(statut = .bincode(
        x = prediction,
        breaks = c(-1e-4, F2_thres, F1_thres, 1 + 1e-4),
        ))  %>%
  dplyr::mutate(statut = factor(statut,
      levels = 1:3,
      labels = c("Pas d'alerte", "Alerte seuil F2", "Alerte seuil F1")
      ))


  # Taux de defaillance moyen avec intervalles de confiance.
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
      moyenne = PointEst, #nolint
      lconf = Lower, #nolint
      hconf = Upper #nolint
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
#' @param eval_frame `data.frame()` \cr Table évaluée tel que donné par la fonction
#' `[aux_custom_eval_urssaf]` avec paramètres wrap_eval_fun et unnest = TRUE.
#'
#' @return NULL
#' @export
aux_custom_plot_urssaf  <- function(eval_frame){
  if (require(ggplot2)) {
    p <- ggplot2::ggplot(
      eval_frame,
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

  plot(p)
  } else {
    warning(
      "Cette fonction necessite le package ggplot2; veuillez l'installer pour
      l'utiliser"
      )
  }
}

#' Fonctions d'évaluation de la précision (spéciale URSSAF)
#'
#' Cette fonction construit un objet eval_function avec les fonctions
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
#' @return `MLsegmentr::eval_function`
#' @export
#'
custom_eval_urssaf  <- function(){
  require(MLsegmentr)
  evo <- MLsegmentr::eval_function(
    eval_fun = aux_custom_eval_urssaf,
    wrap_eval_fun = TRUE,
    unnest =  TRUE,
    plot_fun  = aux_custom_plot_urssaf,
    compulsory_fields = c("siret", "periode")
    )

  return(evo)
}

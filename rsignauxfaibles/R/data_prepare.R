#' Préparation des données pour l'entraînement ou la prédiction
#'
#' Inclut la conversion en H2O, et le target encoding.
#' La map pour le target encoding peut être sauvegardé, afin de permettre de
#' l'appliquer à de nouvelles données.
#' Le target encoding se fait avec du bruit pour éviter le surapprentissage.
#'
#' @param data_to_prepare `data.frame()` \cr Données à préparer
#' @param test_or_train `"test" | "train"` \cr "train" if a new target
#' encoding map (`te_map`) should be computed; "test" otherwise
#' @param te_map `list(H2OFrame)`\cr Carte de calcul du target encoding. Utile
#'   uniquement si `test_or_train` = "test", sinon elle est recalculée.
#' @param save_or_load_map `logical(1)` \cr Faut-il charger la carte de target
#' encoding (si `test_or_train`="test"), ou l'enregistrer (si `test_or_train`
#' = "train"). Chargement et sauvegarde se font avec le nom "te_map",
#' l'extension "te_map" et le chemin d'accès par défaut, cf
#' `[save_h2o_object]`.
#' @param outcome `character(1)`\cr Nom de la cible d'apprentissage sur
#' laquelle effectuer le target encoding.
#'
#' @return `list(2)` \cr
#'   Une liste avec deux champs:
#'   * "data": `H2OFrame`; les données préparées
#'   * "te_map": `list(H2OFrame)`; la carte de target encoding
#' @export
prepare_frame <- function(
                          data_to_prepare,
                          test_or_train,
                          te_map = NULL,
                          save_or_load_map,
                          outcome) {

  #if ( (!is.null(te_map) && test_or_train == "train") ||
  #  (!is.null(te_map) && save_or_load_map == TRUE)) {
  #  error('te_map should not be specified if it is computed (test_or_train =
  #    "train") or if map is loaded (save_or_load_map = TRUE)')
  #}

  h2o_data <- convert_to_h2o(data_to_prepare)

  #
  # Target Encoding de differents groupes sectoriels
  #

  if (test_or_train == "train") {
    te_map <- h2o::h2o.target_encode_create(
      h2o_data,
      x = list(
        c("code_naf"),
        c("code_ape_niveau2"),
        c("code_ape_niveau3")
      ),
      y = outcome
    )

    if (save_or_load_map) {
      save_h2o_object(te_map, "te_map")
    }
  } else if (test_or_train == "test" && save_or_load_map) {
    te_map <- load_h2o_object("te_map", "temap", last = TRUE)
  }

  h2o_data <- h2o_target_encode(
    te_map,
    h2o_data,
    test_or_train
  )

  res <- list(
    data = h2o_data,
    te_map = te_map
  )
  return(res)
}

#' Récupère les dernières données disponibles pour un batch
#'
#' Cette fonction permet d'accéder rapidement aux dernières données
#' disponibles.
#'
#' @param last_batch `character(1)` \cr Batch auquel doit être importées les
#'   données. Les modifications opérées par les batchs ultérieurs sont
#'   ignorées.
#' @param periods `[Date()]` \cr Périodes d'intérêt, auquels charger les
#'   données. Des périodes supplémentairs peuvent être chargées selon la
#'   valeur de rollback_months.
#' @inheritParams mongodb_connection
#' @param fields `character()` \cr Noms des champs à requêter dans la base de
#'   données. Doit contenir "siret" et "periode". Si égal à \code{NULL}, alors
#'   charge tous les champs disponibles.
#' @param min_effectif `integer(1)` \cr Limite basse du filtrage de l'effectif
#'   (la limite est incluse)
#' @param rollback_months `integer(1)`\cr Nombre de mois précédant le premier mois de
#'   `periods` à charger. Permet d'effectuer des calculs de différences ou de
#'   moyennes glissantes pour les périodes d'intérêt.
#'
#' @return `data.frame()` \cr
#' Données avec les colonnes décrites dans `fields`, pour les périodes
#' définies par `periods` et `rollback_months`
#' @export
get_last_batch <- function(
  last_batch,
  periods,
  database,
  collection,
  fields,
  min_effectif,
  rollback_months) {

  current_data <- connect_to_database(
    database,
    collection,
    last_batch,
    date_inf = min(periods) %m-% months(rollback_months),
    date_sup = max(periods) %m+% months(1),
    min_effectif = min_effectif,
    fields = fields
  )

  if ("periode" %in% fields && max(current_data$periode) != max(periods)) {
    log_warn("Data is missing at actual period !")
  }
  return(current_data)
}



h2o_target_encode <- function(
                              te_map,
                              h2o_frame,
                              train_or_test) {
  assertthat::assert_that(train_or_test %in% c("train", "test"))

  if (train_or_test == "train") {
    holdout_type <- "LeaveOneOut"
    blended_avg <- TRUE
    noise_level <- 0.02
  } else if (train_or_test == "test") {
    holdout_type <- "None"
    blended_avg <- FALSE
    noise_level <- 0
  }

  res <- h2o::h2o.target_encode_apply(
    h2o_frame,
    x = as.list(names(te_map)),
    y = "outcome",
    target_encode_map = te_map,
    holdout_type = holdout_type,
    blended_avg = blended_avg,
    # fold_column = "fold_column", should be not necessary
    noise_level = noise_level,
    seed = 1234
  )
  return(res)
}

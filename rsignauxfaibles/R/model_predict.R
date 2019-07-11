#' Prédiction du model
#'
#' Produit une prédiction du modèle sur de nouvelles données. Renvoie le
#' résultat converti sous forme de `data.frame()`
#'
#' @param model `H2OModel()` \cr
#'   Un modèle H2O.
#' @param new_data `H2OFrame()` \cr  Données au bon format pour alimenter le
#' modèle.
#'
#' @return `data.frame()` \cr
#'   Un data.frame avec les colonnes "predicted_outcome" (`TRUE` ou `FALSE`),
#'   "score" (entre 0 et 1), "siret" et "periode"
#' @export
predict_model <- function(model, new_data) {
  prediction <- h2o::h2o.cbind(
    new_data[, c("siret", "periode")],
    h2o::h2o.predict(model, new_data) %>%
      .[, c(1, 3)] %>%
      setNames(list("predicted_outcome", "score"))
  ) %>%
    tibble::as_tibble()

  prediction <- prediction %>%
    mutate(
      # H2O bug ??
      periode = as.Date(structure(periode / 1000,
        class = c("POSIXct", "POSIXt")
      )),
      siret = as.character(siret)
    )
  return(prediction)
}

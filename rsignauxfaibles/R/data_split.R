#' Découper l'échantillon en échantillon d'entraînement, de test et de validation
#'
#' Scinde les données en échantillon d'entraînement, de validation et de
#' test, selon les proportions souhaitées. S'assure que deux établissements de la même entreprise ne soient pas
#' à la fois dans deux échantillons différents pour éviter la fuite
#' d'information d'un échantillon vers l'autre.
#'
#'  La fraction de l'échantillon de test est calculée par
#'  1 - frac_train - frac_val. (frac_train + frac_val) doit donc être inférieur
#'  à 1.
#'
#' @param data_sample `dataframe()`\cr
#' Données à scinder avec des champs "periode" et "siret".
#' @param frac_train `numeric(1)` \cr Fraction des données utilisées pour
#'   l'entraînement. Doit être entre 0 et 1.
#' @param frac_val `numeric(1)` \cr Fraction des données utilisées pour la
#'   validation. Doit être entre 0 et 1.
#' @param seed `integer(1)` \cr Seed pour assurer la reproductibilité des
#'   opérations aléatoires.
#'
#' @return `list(data.frames(3))` \cr
#' Liste de trois data.frames: \code{train}, \code{validation} and \code{test}
#' avec les couples (siret x periode) des trois échantillons calculés.
#' @export
split_snapshot_rdm_month <- function(
  data_sample,
  frac_train,
  frac_val,
  seed = 1234
  ) {

  assertthat::assert_that(
    frac_train > 0,
    frac_val > 0,
    frac_train + frac_val <= 1,
    msg = "Fractions must be positive and not exceed 1"
  )

  frac_test <- 1 - (frac_train + frac_val)

  data_sample <- data_sample %>%
    select(siret, periode) %>%
    arrange(siret, periode) %>%
    mutate(siren = substr(siret, 1, 9))

  sirens <- data_sample %>%
    select(siren) %>%
    distinct()

  set.seed(seed)
  sirens <- sirens %>%
    mutate(ss = sample(1:3,
        size = nrow(sirens), replace = TRUE,
        prob = c(frac_train, frac_val, frac_test)
        )) %>%
    mutate(ss = factor(ss, levels = c(1, 2, 3)))

  data_sample <- data_sample %>%
    left_join(y = sirens, by = "siren") %>%
    select(siret, periode, ss)

  result <- setNames(
    split(data_sample %>% select(-ss), data_sample %>% select(ss)),
    c("train", "validation", "test")
  )

  return(result)
}

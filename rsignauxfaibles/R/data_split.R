#' Découper l'échantillon en échantillon d'entraînement, de test et de validation
#'
#' Scinde les données en échantillon d'entraînement, de validation et de
#' test, selon les proportions souhaitées. S'assure que deux établissements de la même entreprise ne soient pas
#' a la fois dans deux échantillons différents pour éviter la fuite
#' d'information d'un échantillon vers l'autre.
#'
#'  La fraction de l'échantillon de test est calculée par
#'  1 - frac_train - frac_val. (frac_train + frac_val) doit donc être inférieur
#'  a 1.
#'
#' @param data_sample `dataframe()`\cr
#' Données à scinder avec des champs "periode" et "siret".
#' @param fracs `numeric()` \cr Fraction des différents sous-échantillons.
#' Toutes les fractions doivent être dans ]0,1] et leur somme doit être égale
#' à 1.
#' @param names `character()` \cr Vecteur de noms des sous-échantillons, qui doit être de même
#' longueur que les fractions, ou de longueur 1. Si de longueur 1, alors un
#' préfixe "_X", avec X la numérotation des échantillons, est automatiquement ajouté.
#' @param seed `integer(1)` \cr Seed pour assurer la reproductibilité des
#'   opérations aléatoires.
#'
#' @return `list(data.frames(3))` \cr
#' Liste de trois data.frames: \code{train}, \code{validation} and \code{test}
#' avec les couples (siret x periode) des trois échantillons calculés.
#' @export
split_snapshot_rdm_month <- function(
  data_sample,
  fracs,
  names,
  seed = 1234
  ) {

  assertthat::assert_that(
    sum(fracs) == 1,
    msg = "Sum of fractions should be equal to 1"
    )
  assertthat::assert_that(
    all(fracs > 0),
    msg = "All fractions should be strictly positivie"
    )
  assertthat::assert_that(
    length(names) == length(fracs) ||
      length(names) == 1,
    msg = "Les noms specifies doivent etre de meme longueur que les fractions,
    ou de longueur 1")


  data_sample <- data_sample %>%
    select(siret, periode) %>%
    arrange(siret, periode) %>%
    mutate(siren = substr(siret, 1, 9))

  sirens <- data_sample %>%
    select(siren) %>%
    distinct()

  set.seed(seed)

  random_vec <- stats::runif(n = nrow(sirens))
  random_cats <- .bincode(random_vec, breaks = c(0, cumsum(fracs)), TRUE, TRUE)
  sirens <- sirens %>%
    mutate(ss = random_cats) %>%
    mutate(ss = factor(ss, levels = 1:length(fracs)))

  data_sample <- data_sample %>%
    left_join(y = sirens, by = "siren") %>%
    select(siret, periode, ss)

  if (length(names) == 1 && length(fracs) > 1){
    names <- paste0(names, "_", 1:length(fracs))
  }
  result <- setNames(
    base::split(data_sample %>% select(-ss), data_sample %>% select(ss)),
    names
  )

  return(result)
}

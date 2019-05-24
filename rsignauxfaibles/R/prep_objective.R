objective_default_or_failure <- function(data, n_months, threshold, lookback) {
  proc <- c("plan_redressement", "liquidation", "plan_sauvegarde")
  data <- data %>%
    mutate(dette_cumulee_aux = ifelse((!is.na(cotisation_moy12m) & cotisation_moy12m > 1e-10),
      (montant_part_patronale + montant_part_ouvriere) / cotisation_moy12m, 0
    )) %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(
      default_urssaf = check_n_successive_defaults(dette_cumulee_aux, n_months, threshold),
      failure_aux = etat_proc_collective %in% proc,
      default_any = any(default_urssaf | failure_aux),
      failure = with_lookback(failure_aux, lookback),
      default = with_lookback(default_urssaf | failure_aux, lookback)
    ) %>%
    ungroup() %>%
    select(-failure_aux, dette_cumulee_aux)

  assertthat::assert_that(!any(is.na(data$default)))
  assertthat::assert_that(!any(is.na(data$failure)))
  return(data)
}

check_n_successive_defaults <- function(data, n_months, threshold) {
  exceedance <- data >= threshold
  max_consecutive <- sequence(rle(exceedance)$lengths) * exceedance
  return(max_consecutive >= n_months)
}

with_lookback <- function(data, lookback) {
  output <- sapply(1:length(data),
    FUN = function(x)
      any(data[x:min(x + lookback, length(data))])
  )
  return(output)
}

set_objective <- function(data, objective) {
  data[["outcome"]] <- fct_recode(
    as.factor(data[[objective]]),
    default = "TRUE",
    non_default = "FALSE"
  ) %>%
    fct_relevel(c("default", "non_default"))

  cat("L'entraînent se fait désormais sur l'objectif suivant:", objective)
  assertthat::assert_that(!any(is.na(data$outcome)))
  return(data)
}

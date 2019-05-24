elapsed_months <- function(end_date, start_date) {
  ed <- as.POSIXlt(end_date)
  sd <- as.POSIXlt(start_date)
  12 * (ed$year - sd$year) + (ed$mon - sd$mon)
}


compare_vects <- function(data, field1, field2) {
  smy <- data %>%
    group_by(siret) %>%
    summarize(a = first(!!sym(field1)), b = first(!!sym(field2)))
  return(list(only_first = smy$a & !smy$b, only_second = smy$b & !smy$a))
}


count_etab_entr <- function(df) {
  nb_etab <- n_distinct(df %>% select(siret))
  nb_entr <- n_distinct(df %>% select(siren))

  cat(nb_etab, " établissements dans ", nb_entr, " entreprises", "
")
}


#' Replace NAs in a data.frame
#'
#' @param data A data frame
#' @param replace_missing A list where names are column names and values
#' replace values for NAs
#' @param fail_if_column_missing What to do if a replace_missing name is not a column
#' name ? Nothing if drop_silently is true, warning otherwise
#'
#' @return
#' @export
#'
#' @examples
replace_na <- function(
  data,
  replace_missing,
  fail_on_missing_col = TRUE
) {

  if (any(!names(replace_missing) %in% colnames(data)) &&
    fail_on_missing_col &&
    require(logger)){
    stop(
      "{names(replace_missing)[!names(replace_missing) %in% colnames(data)]}
      is missing from the dataframe"
    )
  }

  replace_missing <- replace_missing[names(replace_missing) %in% colnames(data)]

  purrr::walk2(replace_missing, names(replace_missing),
    function(na_value, name) {
      data[is.na(data[, name]), name] <<- na_value
    }
    )
  return(data)
}

average_12m <- function(vec) {
  sapply(
    1:length(vec),
    FUN = function(x) {
      res <- tail(vec[1:x], 12)
      if (sum(!is.na(res)) > 3) {
        return(mean(res, na.rm = TRUE))
      } else {
        return(NA)
      }
    }
  )
}

#' Gives alert levels from prediction and F-scores
#'
#' Lower thresholds are strict (a prediction falling on the threshold is
#' binned to the lower alert level)
#'
#' @param prediction Vector, list of predictions between 0 and 1.
#' @param F1 `Double(1)`. F1_score
#' @param F2 `Double(1)`. F2_score
#'
#' @return A factor vector with alert levels.
#' @export
#'
#' @examples
alert_levels <- function(prediction, F1, F2) {
  assertthat::assert_that(F2 <= F1,
    msg = "F2 score cannot be less than F1 score. Could you have entered the
    scores in the wrong order ?"
  )
  alert <- .bincode(
    x = prediction,
    breaks = c(-1e-4, F2, F1, 1 + 1e-4),
  ) %>%
    factor(
      levels = 1:3,
      labels = c("Pas d'alerte", "Alerte seuil F2", "Alerte seuil F1")
    )
}

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
    msg = "F2 score cannot be less than F1 score. Could you have inverted the
    F scores ?"
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


name_file <- function(
  relative_path,
  file_detail,
  file_extension = "",
  full_path = FALSE) {

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)

  assertthat::assert_that(dir.exists(full_dir_path),
    msg = "Directory not found. Check relative path"
  )

  file_list <- list.files(full_dir_path)


  n_different <- grepl(
    paste0(
      "^", Sys.Date(), "_v[0-9]*_",
      file_detail, "\\.", file_extension, "$"
    ),
    file_list
  ) %>%
    sum()


  file_name <- paste0(
    Sys.Date(),
    "_v",
    n_different + 1,
    "_",
    file_detail,
    ".",
    file_extension
  )

  if (full_path) {
    full_file_path <- file.path(full_dir_path, file_name)

    return(full_file_path)
  } else {
    return(file_name)
  }
}


#' Convert a spark or R dataframe to an H2O Frame
#'
#' It deletes the input frame from memory and fixes H2O types for learning and predicting.
#'
#' @param table_to_convert Spark Dataframe or R data.frame. Deletes the frame from memory.
#'
#' @return The converted H2OFrame
#' @export
#'
#' @examples
convert_to_h2o <- function(table_to_convert) {
  if (table_to_convert %>% inherits("tbl_spark")) {
    sc <- sparklyr::spark_connection_find()[[1]]
    h2o_table <- rsparkling::as_h2o_frame(sc, table_to_convert)
    # tbl_uncache(sc, spark_table_name(table_to_convert)) # Gain some memory space
  } else if (table_to_convert %>% inherits("data.frame")) {
    h2o_table <- h2o::as.h2o(table_to_convert)
    rm(table_to_convert) # Gain some memory space
  }
  h2o_table <- set_h2o_types(h2o_table)
  return(h2o_table)
}


#' Conversion des colonnes "string" en "enum" dans un H2OFrame
#'
#' @param h2o_table Table h2o
#'
#' @return Table h2o avec les colonnes de type "string" remplacées par des colonnes de type "enum"
#' @export
#'
#' @examples
set_h2o_types <- function(h2o_table) {
  fields <- names(h2o_table)[h2o::h2o.getTypes(h2o_table) == "string"]

  if ("outcome" %in% names(h2o_table)) fields <- c(fields, "outcome")

  aux_type_factor <- function(colname) {
    h2o_table[colname] <<- h2o::h2o.asfactor(h2o_table[colname])
  }

  plyr::l_ply(fields, aux_type_factor)
  utils::capture.output(h2o_table, file = "/dev/null") # ? cf JIRA H2O BUG PUBDEV-6221
  return(h2o_table)
}

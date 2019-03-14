#' Title
#'
#' @param donnees
#' @param collection
#' @param export_fields
#' @param database
#' @param last_batch
#'
#' @return
#' @export
#'
#' @examples
prepare_for_export <- function(
  donnees,
  collection,
  export_fields,
  database,
  last_batch,
  known_sirens_filenames = c("sirets_connus_pdl.csv", "sirets_connus_bfc.csv")
){

  last_period  <- max(donnees$periode, na.rm = TRUE)
  cat("Préparation à l'export ... \n")
  cat(paste0("Dernière période connue: ",
             last_period, "\n"))

  full_data <- connect_to_database(
    database,
    collection,
    last_batch,
    date_inf = last_period,
    date_sup = last_period %m+% months(1),
    min_effectif = 10,
    fields = export_fields[!export_fields %in% c("connu", "diff", "prob", "apparait", "disparait")]
  )

  donnees <- donnees %>%
    mutate(siret = as.character(siret)) %>%
    left_join(full_data %>% mutate(siret = as.character(siret)), by = c("siret", "periode")) %>%
    dplyr::mutate(CCSF = date_ccsf) %>%
    dplyr::arrange(dplyr::desc(prob))

  # Report des dernières infos financieres connues

  donnees <- donnees %>%
    mark_known_sirets(names = known_sirens_filenames) %>%
    select(export_fields)

  all_names <- names(donnees)
  cat("Les variables suivantes sont absentes du dataframe:", "\n")
  cat(export_fields[!(export_fields %in% all_names)])
  export_fields <- export_fields[export_fields %in% all_names]


  #if (is.emp)
  to_export <- donnees %>%
    dplyr::select(one_of(export_fields))
}



#' Exports a dataframe to csv
#'
#' @param donnees
#' @param batch
#' @param destination
#' @param relative_path
#'
#' @return
#' @export
#'
#' @examples
export <- function(
  donnees,
  batch,
  destination = "csv",
  relative_path =  file.path("..", "output")) {


  assertthat::assert_that(tolower(destination) %in% c("csv", "mongodb", "json"),
    msg = "Wrong export destination argument.
    Possible destinations are 'csv', 'mongodb' or 'json'")

  if (tolower(destination) == "csv") {


    fullpath <- name_file(
      relative_path,
      file_detail = paste0("detection", batch),
      file_extension = "csv",
      full_name = TRUE
      )

    write.table(donnees,
      row.names = F,
      dec = ",",
      sep = ";",
      file = fullpath,
      quote = T,
      append = F
      )



  } else if (tolower(destination) == "mongodb") {

    error("Export to mongodb not implemented yet !")

  } else if (tolower(destination) == "json"){

    error("Export to json not implemented yet !")

  }
}

mark_known_sirets <- function(df, names){
  sirens <- c()
  for (name in names) {
    sirets <- readLines(rprojroot::find_rstudio_root_file('..','data-raw',name))
    sirens <- c(sirens, substr(sirets,1,9))
  }
  df <- df %>%
    dplyr::mutate(connu = as.numeric(substr(siret,1,9) %in% sirens))

  return(df)
}


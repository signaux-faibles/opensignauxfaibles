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
convert_to_h2o <- function(table_to_convert){
  if (table_to_convert %>% inherits("tbl_spark")) {

    sc <- sparklyr::spark_connection_find()[[1]]
    h2o_table <- rsparkling::as_h2o_frame(sc, table_to_convert)
    #tbl_uncache(sc, spark_table_name(table_to_convert)) # Gain some memory space

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
set_h2o_types <- function(h2o_table){

  fields <- names(h2o_table)[h2o::h2o.getTypes(h2o_table) == "string"]

  if ("outcome" %in% names(h2o_table)) fields <- c(fields, "outcome")

  aux_type_factor <- function(colname){
      h2o_table[colname] <<- h2o::h2o.asfactor(h2o_table[colname])
  }

  plyr::l_ply(fields, aux_type_factor)
  utils::capture.output(h2o_table, file = '/dev/null') #? cf JIRA H2O BUG PUBDEV-6221
  return(h2o_table)
}






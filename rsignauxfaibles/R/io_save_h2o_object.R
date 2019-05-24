#' Enregistrer un objet H2O
#'
#' @param object Objet à sauvegarder. Ou bien un "H2OBinomialModel", ou bien une liste de "H2OFrame"
#' @param object_name Nom générique à donner au fichier, auxquels seront automatiquement accolés une date et une version.
#' @param relative_path Chemin relatif par rapport à la racine du projet R
#'
#' @return
#' @export
#'
#' @examples
save_h2o_object <- function(
                            object,
                            object_name,
                            relative_path = file.path("..", "output", "model")) {
  if (object %>% inherits("H2OBinomialModel")) {
    save_function <- save_H2OModel
    extension <- "model"
  } else if (object %>% inherits("list")) {
    save_function <- save_H2OFrame_list
    extension <- "temap"
  }

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)

  filename <- name_file(
    relative_path,
    file_detail = object_name,
    file_extension = extension
  )

  save_function(object, full_dir_path, filename)

  return(TRUE)
}


#' Enregistrer un objet de classe "H2OBinomialModel"
#'
#' @param object Objet de classe "H2OBinomialModel"
#' @param path Chemin d'accès absolu.
#' @param filename Nom complet du fichier à sauvegarder.
#'
#' @return
#'
#' @examples
save_H2OModel <- function(object, path, filename) {
  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function saves a H2OBinomialModel, not a", class(object))
  )

  h2o.saveModel(object, path)

  # Rename model object, which has the model_id as name at this stage
  file.rename(file.path(path, object@model_id), file.path(path, filename))
}

#' Enregistre une liste de H2OFrames
#'
#' @param object Liste de H2OFrames
#' @param path Chemin absolu
#' @param filename Nom complet du fichier à enregistrer
#'
#' @return
#'
#' @examples
save_H2OFrame_list <- function(object, path, filename) {
  assertthat::assert_that(class(object) == "list",
    msg = paste("This function saves a list, not a", class(object))
  )

  # Convert H2OFrames to dataframes
  res <- lapply(object, as.data.frame)

  saveRDS(res, file.path(path, filename))
}

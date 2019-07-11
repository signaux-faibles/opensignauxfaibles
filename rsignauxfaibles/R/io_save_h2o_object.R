#' Enregistrer un objet H2O
#'
#' Enregistre un objet de type "H2OBinomialModel" ou "H2OFrame".
#' Attention, les modèles enregistrés ne sont pas compatibles d'une version
#' H2O sur l'autre.
#'
#' @param object Objet à sauvegarder. Ou bien un "H2OBinomialModel", ou bien une liste de "H2OFrame"
#' @param object_name `character(1)` \cr Nom générique à donner au fichier, auxquels seront
#' automatiquement accolés une date et une version. Cf `[name_files]`.
#' @param relative_path `character(1)` \cr Chemin relatif par rapport à la racine du projet R
#'
#' @return `character(1)`\cr
#'   Chemin complet du fichier créé.
#' @export
save_h2o_object <- function(
  object,
  object_name,
  relative_path = file.path("..", "output", "model")
  ){

  if (object %>% inherits("H2OBinomialModel")) {
    save_function <- save_H2OModel
    extension <- "model"
  } else if (object %>% inherits("list")) {
    save_function <- save_H2OFrame_list
    extension <- "temap"
  }

  require(rprojroot)
  absolute_path <- rprojroot::find_rstudio_root_file(relative_path)

  filename <- name_file(
    relative_path,
    file_detail = object_name,
    file_extension = extension
  )

  save_function(object, absolute_path, filename)

  return(file.path(absolute_path, filename))
}

#' Enregistre une liste de H2OFrames
#'
#' Enregistre une liste de `H2OFrame`.
#'
#' @param object `list(H2OFrame)` \cr Liste de H2OFrames
#' @param absolute_path `character()`\cr Chemin absolu (sans nom de fichier)
#' @param filename `character()`\cr Nom complet du fichier à enregistrer
#'
#' @return `TRUE`
save_H2OFrame_list <- function(object, absolute_path, filename) {
  assertthat::assert_that(class(object) == "list",
    msg = paste("This function saves a list, not a", class(object))
  )

  # Convert H2OFrames to dataframes
  res <- lapply(object, as.data.frame)

  saveRDS(res, file.path(absolute_path, filename))
  return(TRUE)
}

#' Enregistrer un objet de classe "H2OBinomialModel"
#'
#' @param object `H2OBinomialModel`\cr Modèle H2O à enregistrer
#' @inheritParams save_H2OFrame_list
#' @return `TRUE`
save_H2OModel <- function(object, absolute_path, filename) {
  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function saves a H2OBinomialModel, not a", class(object))
  )

  h2o::h2o.saveModel(object, absolute_path)

  # Rename model object, which has the model_id as name at this stage
  file.rename(
    file.path(absolute_path, object@model_id),
    file.path(absolute_path, filename)
  )
  return(TRUE)
}

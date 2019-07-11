#' Charge un fichier enregistré par save_h2o_object
#'
#' Cette fonction charge un fichier enregistré avec \code{save_h2o_object}, ou bien avec la racine \code{name}, l'extension du fichier, et le paramètre \code{last = TRUE}(qui charge le dernier fichier disponible), ou bien avec le nom exact avec le paramètre \code{file_name} et \code{last = FALSE}
#'
#' @param object_name `character(1)`\cr Nom de l'objet, spécifié de manière
#' identique qu'au moment de l'enregistrement. Cf `[save_h2o_object]`
#' @param extension `character(1)`\cr
#'   L'extension du fichier à charger, ou bien "model" (pour un
#'   `H2OBinomialModel`) ou bien "temap" (pour une `list(H2OFrame)`).
#' @param relative_path Chemin d'accès relatif par rapport à la racine du projet.
#' @param last `logical()` \cr Si TRUE, charge le dernier fichier sauvegardé
#'  correspondant à `object_name`, `extension` et `relative_path`. Sinon, il
#'  faut spécifier le nom complet: `file_name`.
#' @param file_name `character(1)`\cr Nom complet et exact du fichier à charger. N'est pris en compte uniquement si last est FALSE.
#'
#' @return `H2OBinomialModel() | `list(H2OFrame)` \cr
#'   depending on the loaded model.
#' @export
load_h2o_object <- function(
                            object_name,
                            extension = NULL,
                            relative_path = file.path("..", "output", "model"),
                            last = TRUE,
                            file_name = "") {
  if (substr(relative_path, nchar(relative_path), nchar(relative_path))
  %in% c("/", "\\")) {
    relative_path <- substr(relative_path, 1, nchar(relative_path) - 1)
  }

  if (file_name != "") {
    fnsplit <- strsplit(file_name, ".", fixed = TRUE)
    extension <- tail(fnsplit[[1]], n = 1)
  }

  assertthat::assert_that(extension %in% c("model", "temap"),
    msg = 'Unsupported extension. Supported extensions are "model" and "temap"'
  )

  if (extension == "model") {
    load_function <- load_H2OModel
  } else if (extension == "temap") {
    load_function <- load_H2OFrame_list
  }

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)
  assertthat::assert_that(dir.exists(full_dir_path),
    msg = "Directory not found. Check relative path"
  )

  if (last) {
    file_candidates <- list.files(full_dir_path) %>%
      grep(pattern = paste0(object_name, ".", extension), value = TRUE)

    assertthat::assert_that(length(file_candidates) > 0,
      msg = "No such file, please check object_name and extension"
    )

    file_name <- file_candidates %>%
      sort(decreasing = TRUE) %>%
      .[1]
  }

  full_path <- file.path(full_dir_path, file_name)

  assertthat::assert_that(file.exists(full_path),
    msg = "No such file, please check file_name"
  )

  res <- load_function(full_dir_path, file_name)

  return(res)
}

#' Chargement d'une liste de H2OFrames
#'
#' Charge une liste de H2OFrames. Utilisé pour le target encoding.
#' @inheritParams save_H2OFrame_list
#' @return `H2OBinomialModel`
load_H2OFrame_list <- function(absolute_path, filename) {
  object <- readRDS(file.path(absolute_path, filename))

  assertthat::assert_that(class(object) == "list",
    msg = paste("This function loads a list, not a", class(object))
  )

  # Convert dataframes to H2OFrames
  res <- lapply(object, h2o::as.h2o)

  return(res)
}

#' Chargement d'un fichier de classe "H2OBinomialModel"
#'
#' Charge un fichier de classe "H2OBinomialModel".
#'
#' @param absolute_path `character(1)`\cr Chemin d'accès
#' @param filename Nom exact du fichier à charger.
#'
#' @return `list(H2OFrame)`
load_H2OModel <- function(absolute_path, filename) {
  object <- h2o::h2o.loadModel(file.path(absolute_path, filename))

  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function loads a H2OBinomialModel, not a", class(object))
  )

  return(object)
}

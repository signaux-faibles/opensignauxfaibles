#' Charge un fichier enregistré par save_h2o_object
#'
#' Cette fonction charge un fichier enregistré avec \code{save_h2o_object}, ou bien avec la racine \code{name}, l'extension du fichier, et le paramètre \code{last = TRUE}(qui charge le dernier fichier disponible), ou bien avec le nom exact avec le paramètre \code{file_name} et \code{last = FALSE}
#'
#' @param name Nom de l'objet
#' @param extension L'extension du fichier à charger, ou bien "model" ou bien "temap"
#' @param relative_path Chemin d'accès relatif par rapport à la racine du projet.
#' @param last Si TRUE, charge le dernier fichier sauvegardé.
#' @param file_name Nom exact du fichier à charger. N'est pris en compte uniquement si last est FALSE.
#'
#' @return
#' @export
#'
#' @examples
load_h2o_object <- function(
                            name,
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
      grep(pattern = paste0(name, ".", extension), value = TRUE)

    assertthat::assert_that(length(file_candidates) > 0,
      msg = "No such file, please check name and extension"
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
#' @param path Chemin d'accès relatif
#' @param filename Nom complet du fichier
#'
#' @return
#'
#' @examples
load_H2OFrame_list <- function(path, filename) {
  object <- readRDS(file.path(path, filename))

  assertthat::assert_that(class(object) == "list",
    msg = paste("This function loads a list, not a", class(object))
  )

  # Convert dataframes to H2OFrames
  res <- lapply(object, as.h2o)

  return(res)
}

#' Chargement d'un fichier de classe "H2OBinomialModel"
#'
#' Charge un fichier de classe "H2OBinomialModel".
#'
#' @param path Chemin d'accès
#' @param filename Nom exact du fichier à charger.
#'
#' @return
#' @examples
load_H2OModel <- function(path, filename) {
  object <- h2o.loadModel(file.path(path, filename))

  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function loads a H2OBinomialModel, not a", class(object))
  )

  return(object)
}


read_h2oframe_from_csv <- function() {
  path <- rprojroot::find_rstudio_root_file(
    "..", "output", "features", "features.csv"
  )

  # Read csv file
  return(h2o::h2o.importFile(path))
}

# Generic functions

#' Chargement de données historiques
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#'
#' @return
#' @export
#'
#' @examples
load_hist_data <- function(task, ...){
  UseMethod("load_hist_data", task)
}

#' Loads new data
#' @export
load_new_data <- function(task, ...){
  UseMethod("load_new_data", task)
}

#' Splits data
#' @export
split_data <- function(task, ...){
  UseMethod("split_data", task)
}

#' Prepare data
#' @export
prepare <- function(task, ...){
  UseMethod("prepare", task)
}

#' Optimize hyperparameters
#' @export
optimize_hyperparameters <- function(task, ...){
  UseMethod("optimize_hyperparameters", task)
}

#' Trains model on data
#' @export
train <- function(task, ...){
  UseMethod("train", task)
}

#' Loads task
#' @export
load <- function(task, ...){
  UseMethod("load", task)
}

load.default <- function(task, ...){
  base::load(task, ...)
}

#' Saves task
#' @export
save <- function(task, ...){
  UseMethod("save", task)
}

save.default <- function(task, ...){
  base::save(task, ...)
}

#' exports data
#' @export
export <- function(task, ...){
  UseMethod("export", task)
}

#' evaluates a model
#' @export
evaluate <- function(task, ...){
  UseMethod("evaluate", task)
}

#' Log a machine learning experiment
#' @export
log <- function(task, ...){
  UseMethod("log", task)
}
#' Explains model results
#' @export
explain <- function(task, ...){
  UseMethod("explain", task)
}

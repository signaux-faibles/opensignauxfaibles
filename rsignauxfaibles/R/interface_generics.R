# Generic functions

#' Chargement de données historiques
#'
#' @param task `[sf_task]` \cr Objet s3 de type sf_task
#' @param ... Additional parameters depending on class of task
#' @export
load_hist_data <- function(task, ...){
  UseMethod("load_hist_data", task)
}

#' Loads new data
#' @inheritParams load_hist_data
#' @export
load_new_data <- function(task, ...){
  UseMethod("load_new_data", task)
}

#' Splits data
#' @inheritParams load_hist_data
#' @export
split_data <- function(task, ...){
  UseMethod("split_data", task)
}

#' Prepare data
#' @inheritParams load_hist_data
#' @export
prepare <- function(task, ...){
  UseMethod("prepare", task)
}

#' Optimize hyperparameters
#' @inheritParams load_hist_data
#' @export
optimize_hyperparameters <- function(task, ...){
  UseMethod("optimize_hyperparameters", task)
}

#' Trains model on data
#' @inheritParams load_hist_data
#' @export
train <- function(task, ...){
  UseMethod("train", task)
}

#' Loads task
#' @inheritParams load_hist_data
#' @export
load <- function(task, ...){
  UseMethod("load", task)
}

load.default <- function(task, ...){
  base::load(task, ...)
}

#' Saves task
#' @inheritParams load_hist_data
#' @export
save <- function(task, ...){
  UseMethod("save", task)
}

save.default <- function(task, ...){
  base::save(task, ...)
}

#' exports data
#' @inheritParams load_hist_data
#' @export
export <- function(task, ...){
  UseMethod("export", task)
}

#' evaluates a model
#' @inheritParams load_hist_data
#' @export
evaluate <- function(task, ...){
  UseMethod("evaluate", task)
}

#' Log a machine learning experiment
#' @inheritParams load_hist_data
#' @export
log <- function(task, ...){
  UseMethod("log", task)
}

#' Explains model results
#' @inheritParams load_hist_data
#' @export
explain <- function(task, ...){
  UseMethod("explain", task)
}

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
load_new_data <- function(task, ...){
  UseMethod("load_new_data", task)
}
hold_out <- function(task, ...){
  UseMethod("hold_out", task)
}
prepare <- function(task, ...){
  UseMethod("prepare", task)
}
optimize_hyperparameters <- function(task, ...){
  UseMethod("optimize_hyperparameters", task)
}
train <- function(task, ...){
  UseMethod("train", task)
}
load <- function(task, ...){
  UseMethod("load", task)
}
load.default <- function(task, ...){
  base::load(task, ...)
}
save <- function(task, ...){
  UseMethod("save", task)
}
save.default <- function(task, ...){
  base::save(task, ...)
}
export <- function(task, ...){
  UseMethod("export", task)
}
evaluate <- function(task, ...){
  UseMethod("evaluate", task)
}
log <- function(task, ...){
  UseMethod("log", task)
}
explain <- function(task, ...){
  UseMethod("explain", task)
}

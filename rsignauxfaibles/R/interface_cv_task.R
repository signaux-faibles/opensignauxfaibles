#' Prepare cross-validation
#'
#' Splits the data into `n_folds` for cross validation
#'
#' @inheritParams split_data.sf_task
#' @param n_folds `integer(1)` \cr Number of folds
#' @param frac_test `numeric(1)` \cr Fraction of data which should be kept
#'   aside
#'
#' @return `c(cv_task, sf_task)` \cr
#'  A cross validated task with field "cross_validation", which is a list of
#'  tasks, one for each fold.
#' @export
split_n_folds <- function(
  task,
  n_folds,
  frac_test,
  remove_strong_signals
  # TODO: not working
  # 55015842-b766-4c4d-8d0e-51b814803479
){

  assertthat::assert_that(frac_test >= 0 && frac_test < 1)
  frac_cv <- (1 - frac_test) / n_folds
  cv_names  <- paste0("cv_", 1:n_folds)
  task <- split_data(
    task,
    fracs = c(rep(frac_cv, n_folds), frac_test),
    names = c(cv_names, "test"),
    remove_strong_signals = remove_strong_signals
  )
  cv_names <- paste0(cv_names, "_data")
  cv_chunks <- task[cv_names]
  task[cv_names] <- NULL

  create_cv_task <- function(cv_number, cv_chunks){
    cv_task  <- sf_task(
        verbose = FALSE,
        database = task[["database"]],
        collection = task[["collection"]],
        experiment_name = paste0(
          #TODO change in special field
          attr(task, "to_log")[["experiment_aim"]],
          "_cv",
          cv_number
        ),
      experiment_description = attr(task, "to_log")[["experiment_description"]]
      )
    cv_task[["validation_data"]] <- cv_chunks[[cv_number]]
    cv_task[["train_data"]] <- dplyr::bind_rows(cv_chunks[-cv_number])

    attr(cv_task, "to_log")[["resampling_strategy"]] <- paste0(
      n_folds,
      "-folds cross validation"
    )

    return(cv_task)
  }

  task[["cross_validation"]] <- purrr::map(
    1:n_folds,
    create_cv_task,
    cv_chunks = cv_chunks
  )

  if (!"cv_task" %in% class(task)){
    class(task) <- c("cv_task", class(task))
  }

  return(task)
}

#' Prepares a cross-validated task
#'
#' Prepares all subtasks (one for each fold) from cv_task.
#'
#' @inheritParams prepare.sf_task
#'
#'  @return The task given as input where each fold-task has been prepared.
#'  @export
prepare.cv_task <- function(
  task,
  data_names = c("train_data",
    "validation_data",
    "test_data",
    "new_data"
    )
  ){

  task[["cross_validation"]] <- purrr::map(
    .x = task[["cross_validation"]],
    .f = prepare.sf_task,
    data_names = data_names
  )
  return(task)
}

#' Train a model on a cross-validated task
#'
#' Train models for each subtask, i.e. for each cross-validation fold.
#'
#' @inheritParams train.sf_task
#'
#'  @return The input task, where each cv-subtask has a trained model
#'  @export
train.cv_task  <- function(
  task,
  fields = get_fields(training = TRUE),
  parameters = NULL,
  seed = 123
  ){

  task[["cross_validation"]]  <- purrr::map(
    task[["cross_validation"]],
    train.sf_task,
    fields = fields,
    parameters = parameters,
    seed = seed
    )

  return(task)
}

#' Predict model on new data for a cross-validated task
#'
#' Predict on new data for each cross-validated fold.
#'
#' @inheritParams predict.sf_task
#'
#'  @return A task where each cross-validated fold has predictions
#'  @export
predict.cv_task <- function(
  task,
  data_names = c(
    "new_data",
    "train_data",
    "validation_data",
    "test_data"
    )
  ){

  task[["cross_validation"]] <- purrr::map(
    task[["cross_validation"]],
    predict.sf_task,
    data_names = data_names
    )
  return(task)
}


#' Evaluate a cross-validated task
#'
#'  @inheritParams evaluate.sf_task
#'
#'  @return task with evaluation on first level. (no individual
#'  cross-validated fold evaluation)
#'  @export
evaluate.cv_task <- function(
  task,
  eval_function = MLsegmentr::eval_precision_recall(),
  data_name = c("validation_data"),
  plot = TRUE,
  prediction_names = "score",
  target_names = "outcome",
  segment_names = NULL
  ){

  combined_data <- purrr::map(task[["cross_validation"]], "validation_data")
  #TODO no need for tibble_time !!
  # Task f338d4ff-6516-480a-bf85-7c5c1da4fbd7
  combined_data <- purrr::map(combined_data, tibble::as_tibble)  %>%
    dplyr::bind_rows(.id = "model") %>%
    mutate(model = paste0("cv_", model))
  model_names  <- unique(combined_data$model)
  combined_data  <- combined_data %>%
    tidyr::spread(key = model, value = "score")

  browser()
  modified_task <- task
  modified_task[["validation_data"]] <- combined_data
  modified_task <- evaluate.sf_task(
    modified_task,
    eval_function = eval_function,
    data_name = data_name,
    plot = plot,
    prediction_names = model_names,
    target_names = target_names,
    segment_names = NULL
  )
  # TODO: why isn't there the column model in output ?
  # Task 282fe34e-45df-4c5d-9e57-0ec0ef2e8a77

  return (modified_task)
}

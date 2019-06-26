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
  n_folds = 4,
  frac_test = 0.2,
  tracker = task[["tracker"]]
){
  require(purrr)

  assertthat::assert_that(frac_test >= 0 && frac_test < 1)
  frac_cv <- (1 - frac_test) / n_folds
  cv_names  <- paste0("cv_", 1:n_folds)
  task <- split_data(
    task,
    fracs = c(rep(frac_cv, n_folds), frac_test),
    names = c(cv_names, "test")
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
          task[["tracker"]]$values[["experiment_name"]],
          "_cv",
          cv_number
        ),
      experiment_description =
        task[["tracker"]]$values[["experiment_description"]]
      )
    cv_task[["validation_data"]] <- cv_chunks[[cv_number]]
    cv_task[["train_data"]] <- dplyr::bind_rows(cv_chunks[-cv_number])

    if (!is.null(tracker)){
      suppressWarnings(
        tracker$set(
          resampling_strategy = paste0(
            n_folds,
            "-folds cross validation"
            )
          )
        )
    }

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
#' @return The task given as input where each fold-task has been prepared.
#' @export
prepare.cv_task <- function(
  task,
  data_names = c("train_data",
    "validation_data",
    "test_data",
    "new_data"
    )
  ){

  require(purrr)

  task[["cross_validation"]] <- purrr::map(
    .x = task[["cross_validation"]],
    .f = prepare.sf_task,
    data_names = data_names,
    tracker = task[["tracker"]]
  )
  return(task)
}

#' Train a model on a cross-validated task
#'
#' Train models for each subtask, i.e. for each cross-validation fold.
#'
#' @inheritParams train.sf_task
#' @return The input task, where each cv-subtask has a trained model
#' @export
train.cv_task  <- function(
  task,
  fields = get_fields(training = TRUE),
  parameters = NULL,
  seed = 123
  ){

  require(purrr)
  if (is.null(parameters) && "model_parameters" %in% names(task)){
    parameters <- task[["model_parameters"]]
  }

  task[["features"]] <- fields
  # Sort out what should be stored in parent task and in cross validation
  # folds
  task[["cross_validation"]]  <- purrr::map(
    task[["cross_validation"]],
    train.sf_task,
    fields = fields,
    parameters = parameters,
    seed = seed,
    tracker = task[["tracker"]]
    )

  return(task)
}

#' Predict model on new data for a cross-validated task
#'
#' Predict on new data for each cross-validated fold.
#'
#' @inheritParams predict.sf_task
#' @return A task where each cross-validated fold has predictions
#' @export
predict.cv_task <- function(
  object,
  data_names = c(
    "new_data",
    "train_data",
    "validation_data",
    "test_data"
    )
  ){

  require(purrr)
  task <- object
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
  segment_names = NULL,
  remove_strong_signals = TRUE
  ){

  require(purrr)
  combined_data <- purrr::map(task[["cross_validation"]], "validation_data")
  combined_data <- combined_data  %>%
    dplyr::bind_rows(.id = "model") %>%
    mutate(model = paste0("cv_", model))
  model_names  <- unique(combined_data$model)
  combined_data  <- combined_data %>%
    tidyr::spread(key = model, value = "score")

  modified_task <- task
  modified_task[["validation_data"]] <- combined_data
  modified_task <- evaluate.sf_task(
    modified_task,
    eval_function = eval_function,
    data_name = data_name,
    plot = plot,
    prediction_names = model_names,
    target_names = target_names,
    segment_names = NULL,
    remove_strong_signals = remove_strong_signals,
    tracker = task[["tracker"]]
  )

  # TODO: why isn't there the column model in output ?
  # Task 282fe34e-45df-4c5d-9e57-0ec0ef2e8a77
  #
  # TODO: next section to be replaced with build-in summary functions of
  # MLsegmentr
  # see c5823da4-48ae-4c33-a345-2a5efdc808e7

  summary_eval <- modified_task[["model_performance"]] %>%
    filter(evaluation_name == "aucpr")  %>%
    group_by(target_type, segment, evaluation_name) %>%
    tidyr::unnest() %>%
    summarize(evaluation = mean(evaluation), model = factor("cv_averaged")) %>%
    #to have same col order
    select(names(modified_task[["model_performance"]])) %>%
    ungroup()

  summary_eval$evaluation <- list(summary_eval$evaluation)
  summary_eval  <- summary_eval %>%
    mutate_if(is.factor, as.character)
  modified_task[["model_performance"]] <-
    modified_task[["model_performance"]] %>%
    mutate_if(is.factor, as.character)

  modified_task[["model_performance"]] <- bind_rows(
    summary_eval,
    modified_task[["model_performance"]]
    )
  return (modified_task)
}

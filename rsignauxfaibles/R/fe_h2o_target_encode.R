h2o_target_encode <- function(
                              te_map,
                              h2o_frame,
                              train_or_test) {
  assertthat::assert_that(train_or_test %in% c("train", "test"))

  if (train_or_test == "train") {
    holdout_type <- "LeaveOneOut"
    blended_avg <- TRUE
    noise_level <- 0.02
  } else if (train_or_test == "test") {
    holdout_type <- "None"
    blended_avg <- FALSE
    noise_level <- 0
  }

  res <- h2o.target_encode_apply(
    h2o_frame,
    x = as.list(names(te_map)),
    y = "outcome",
    target_encode_map = te_map,
    holdout_type = holdout_type,
    blended_avg = blended_avg,
    # fold_column = "fold_column", should be not necessary
    noise_level = noise_level,
    seed = 1234
  )



  return(res)
}

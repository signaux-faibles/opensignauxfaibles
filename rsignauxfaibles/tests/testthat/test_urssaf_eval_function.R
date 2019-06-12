context("Test urssaf specific eval_function")

eval_frame <- MLsegmentr::get_test_eval_frame()
n <- nrow(eval_frame)

set.seed(1)
additional_data <- data.frame(
  .id = 1:n,
  siret = sample(1:(n / 2), n, replace = TRUE),
  periode = round(runif(n) + 1)
  )

eval_frame  <- eval_frame %>%
  left_join(additional_data, by = c(".id"))

testthat::test_that("custom_eval_urssaf works correctly", {
  eval_obj <- eval_urssaf()
  expect_error(
    actual_out <- MLsegmentr::assess_eval_frame(eval_obj, eval_frame, plot = FALSE),
    NA
    )
  expected_moyenne  <- c(0.07692, 0.33333, 0.95833, 0, 0, 0.81818,
    0, 0.5, 0.58823, 0, 0, 0.60714)
  expected_lconf  <- c(0.00394, 0.01709, 0.79758, 2.20966e-17, 0,
    0.61483, 0, 0.02564, 0.42221, 0, 0, 0.42409)
  expected_hconf  <- c(0.33313, 0.79234, 0.99786, 0.20388, 0.94870, 0.92693,
    0.65761, 0.97435, 0.73634, 0.56149, 0.32440, 0.76434)

  expect_equal(actual_out$moyenne, expected_moyenne, tolerance = 1e-2 )
  expect_equal(actual_out$lconf, expected_lconf, tolerance = 1e-2 )
  expect_equal(actual_out$hconf, expected_hconf, tolerance = 1e-2 )
  }
  )

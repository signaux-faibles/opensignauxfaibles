context("Test urssaf specific eval_function")

eval_frame <- MLsegmentr::get_test_eval_frame()
n <- nrow(eval_frame)

set.seed(1)
additional_data <- data.frame(
  .id = 1:n,
  siret = sample(1:(n / 2), n, replace = TRUE),
  periode = round(runif(n) + 1)
 )

testthat::test_that("custom_eval_urssaf works correctly", {
  actual_out <- custom_eval_urssaf(eval_frame, additional_data)
  actual_moyenne  <- c(0.125, 0.913, 0, 0.958, 0.133, 0.5, 0.909, 0, 0, 0.958)
  actual_lconf  <- c(0.034, 0.732, 1.131e-17, 0.797, 0.037,
    0.025, 0.721, 0, 0, 0.797)
  actual_hconf  <- c(0.360, 0.975, 0.184, 0.997, 0.378, 0.974,
    0.974, 0.215, 0.561, 0.997)
  expect_equal(actual_out$moyenne, actual_moyenne, tolerance = 1e-2 )
  expect_equal(actual_out$lconf, actual_lconf, tolerance = 1e-2 )
  expect_equal(actual_out$hconf, actual_hconf, tolerance = 1e-2 )

  expect_error(custom_plot_urssaf(actual_out) , NA)
}
)

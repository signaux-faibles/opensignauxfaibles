context("test utilitary functions")

test_that("alert_levels works as expected", {
  fake_predictions <- c(0, 0.1, 0.2, 0.21, 0.3, 0.4, 0.5, 1)
  fake_F2 <- 0.2
  fake_F1 <- 0.4
  actual <- alert_levels(fake_predictions, fake_F1, fake_F2)
  expected <- structure(c(1L, 1L, 1L, 2L, 2L, 2L, 3L, 3L), .Label = c(
    "Pas d'alerte",
    "Alerte seuil F2", "Alerte seuil F1"
  ), class = "factor")
  expect_equal(expected, actual)
})

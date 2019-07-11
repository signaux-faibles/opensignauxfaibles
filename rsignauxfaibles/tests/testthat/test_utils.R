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


test_that("replace_na works as expected", {

  a <- data.frame(foo = c(1, NA, 3), bar = c(NA, NA, 3))
  actual <- replace_na(a, list(foo = 2, bar = 1))
  expected <- data.frame(foo = c(1, 2, 3), bar = c(1, 1, 3))
  expect_equal(actual, expected)
  expect_error(replace_na(a, list(toto = 3)))
  expect_error(
    actual <- replace_na(
      a,
      list(toto = 3, foo = 2, bar = 1),
      fail_if_column_missing = FALSE
    ),
    NA)
  expect_equal(actual, expected)
})

context("Test s3 sf_tasks")


test_that("sf_task initialization works as expected", {
  my_test_task  <<- expect_error(
    sf_task(
      verbose = TRUE,
      database = "unittest_signauxfaibles",
      collection  = "Features_for_tests",
      experiment_name = "test",
      experiment_description = "test description"),
    NA
    )

  expect_true(
    all(
      names(my_test_task) %in% c("database", "collection")
      )
    )
  expect_true( attr(my_test_task, "verbose"))
  expect_equal( attr(my_test_task, "to_log"),
    list(
      experiment_aim = "test",
      experiment_description = "test description"
      )
    )
})


test_that("load_hist_data works as expected", {
  expect_error(
    my_test_task <<- load_hist_data.sf_task(
      my_test_task,
      batch = "1901_interim"
    ),
    NA
    )

  expect_true("hist_data" %in% names(my_test_task))
  expect_equal(nrow(my_test_task[["hist_data"]]), 440)
})

test_that("split_data works as expected", {
  expect_error(split_data(my_test_task, fracs <- c (0.6, 0.1, 0.1)))
  expect_error(
    my_test_task  <<- split_data(my_test_task),
    NA
  )
  expect_true(
    all(
      c("train_data", "test_data", "validation_data") %in% names(my_test_task)
      )
    )

  expect_error(split_data(sf_task(
        verbose = FALSE,
        experiment_name = "empty_task",
        experiment_description = ""
        )))
})


test_that("Data preparation works as expected", {
  connect_to_h2o()
  expect_error(
    my_test_task <<- prepare(my_test_task, c(
      "train_data",
      "validation_data",
      "new_data"
    )),
    NA
    )
  expect_error(
    second_time <<- prepare(my_test_task, "validation_data"),
    NA
    )
  expect_true(all(
      c("prepared_train_data", "prepared_validation_data") %in%
        names(my_test_task)
      ))
})

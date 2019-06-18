context("Test s3 sf_tasks")

test_database <- "unittest_signauxfaibles"
test_that("sf_task initialization works as expected", {
  expect_error(
    my_test_task <<- sf_task(
      verbose = TRUE,
      database = test_database,
      collection  = "Features_for_tests",
      experiment_name = "test",
      experiment_description = "test description"),
    NA
    )

  expect_true(
    all(
      names(my_test_task) %in% c("database", "collection", "tracker")
      )
    )
  expect_true( attr(my_test_task, "verbose"))
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

test_that("Model training works as expected", {
  expect_error(
    my_test_task <<- train(my_test_task),
    NA
    )
  wrong_task <- my_test_task
  wrong_task[["prepared_train_data"]] <- NULL
  expect_error(train(wrong_task))
  ## TODO check setting parameters before works
  ## setting parameters in function calls overwrites
  # 8cdb79fe-caf5-4e9c-aee0-c64c8f2f519d
})


test_that("Predicting works", {
  expect_error(
    my_test_task <<- predict(my_test_task, data_names = "validation_data"),
    NA
    )
  expect_true("score" %in% names(my_test_task[["validation_data"]]))
  ## Predicting twice still works
  my_test_task <- predict(my_test_task, data_names = "validation_data")
  expect_true("score" %in% names(my_test_task[["validation_data"]]))

  ## If data is not prepared: log_warn, no error
  expect_error(predict(my_test_task, data_names = "test_data"),
    NA)
})


test_that("evaluation works", {

   expect_error(
     my_test_task <<- evaluate(
       my_test_task,
       plot = FALSE,
       remove_strong_signals = FALSE
     ),
      NA
      )
   expect_error(
     my_test_task <<- evaluate(my_test_task, plot = FALSE),
      NA
      )
  #TODO better test remove_strong_signals
  # 35191df2-7b94-4cc7-8646-5a9f97921d1d
})


test_that("log works", {
  expect_error(
    log(
      my_test_task,
      database = test_database,
      collection = "ml_logs_for_tests"
    ),
    NA
    )
  # TODO clean mllogs for tests database
  # ef922acc-3f0c-4ff8-8c51-cbd9a53eec02
})

test_that("load_new_data works", {
  expect_error(
    my_test_task <<- load_new_data(
      my_test_task,
      as.Date("2019-01-01"),
      "1901_interim"
    ),
    NA
    )
  expect_equal(nrow(my_test_task[["new_data"]]), 31)
})

test_that("Export works", {

})

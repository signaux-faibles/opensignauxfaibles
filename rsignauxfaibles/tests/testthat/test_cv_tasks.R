context("Test cross validation tasks")

test_database <- "unittest_signauxfaibles"

connect_to_h2o()

test_that("quick test that everything works", {
  my_test_task <- sf_task(
    verbose = TRUE,
    database = test_database,
    collection  = "Features_for_tests",
    experiment_name = "test",
    experiment_description = "test description")


  my_test_task <- load_hist_data(
    task = my_test_task,
    batch = "1901_interim"
    )



  my_test_task <- split_n_folds(
    my_test_task,
    n_folds = 2,
    frac_test = 0.2
  )

  # Afin qu'il y ait des issues positives
  set.seed(12321)
  my_test_task[["cross_validation"]]  <- lapply(
    my_test_task[["cross_validation"]],
    function(x){
      x[["validation_data"]]$outcome  <- round(
        runif(nrow(x[["validation_data"]]))
      )
      x[["train_data"]]$outcome  <- round(
        runif(nrow(x[["train_data"]]))
      )
      return(x)
    }
  )

  my_test_task  <- prepare(
    my_test_task,
    data_names = c("train_data", "validation_data")
  )
  my_test_task  <- train(my_test_task)
  my_test_task <- predict(my_test_task, data_names = "validation_data")
  my_test_task <- evaluate(my_test_task, data_name = "validation_data")

  expect_equal(
    my_test_task[["model_performance"]]$evaluation[1][[1]],
    0.465518,
    tolerance = 1e-4
    )

  my_test_task  <- optimize_hyperparameters(
    my_test_task,
    n_init = 5,
    n_iter = 5,
    fields = c("montant_part_patronale", "ca")
    )
  my_test_task <- log(
    my_test_task,
    database = "unittest_signauxfaibles",
    collection = "ml_logs_for_tests"
  )
})

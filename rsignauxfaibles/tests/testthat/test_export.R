context("Test that exports and import of scores work as expected")


# test_dataframe <- data.frame(
#   siret = c("1", "1", "2"),
#   periode = c("2018", "2019", "2019"),
#   score = c(0.4, 0.6, 0.2),
#   score_diff = c(0, 0.2, -0.2)
# )


test_database <- "unittest_signauxfaibles"
test_collection <- "Scores_for_tests"

dbconnection <- mongolite::mongo(
  collection = test_collection,
  db = test_database,
  verbose = TRUE,
  url = "mongodb://localhost:27017"
)

f_scores = c(F2 = 0.2, F1 = 0.4)
# test_that("export_scores_to_mongodb works as expected", {

#     expect_error(
#       export_scores_to_mongodb(
#         formatted_data = test_dataframe,
#         algo = "test_algo",
#         batch = "test_batch",
#         f_scores = c(F2 = 0.2, F1 = 0.4),
#         database = test_database,
#         collection = test_collection
#         ),
#       NA
#       )

# })

# Sample object for the tests
siret_df <- data.frame(
  siret = as.character(rep(1:10, each = 3)),
  periode = seq.Date(
    from = as.Date("2019-01-01"),
    by = "month",
    length.out = 3
  ),
  score_diff = 1
)


set.seed(1)
scores_1 <- runif(n = nrow(siret_df))
scores_2 <- runif(n = nrow(siret_df))
scores_3 <- runif(n = nrow(siret_df))


# Empty test collection

# ------------------------------------------------------------------------------
# -------- Export to mongodb ---------------------------------------------------
# ------------------------------------------------------------------------------
test_that("Test inputs for export_scores function to mongodb", {
  expect_error(export_scores_to_mongodb(
    formatted_data = cbind(siret_df, score = scores_1),
    batch = 123,
    algo = "test_algo",
    f_scores = c(F2 = 0.2, F1 = 0.4),
    database = test_database,
    collection = test_collection
  ))

  expect_error(export_scores(
    formatted_data = cbind(siret_df, score = scores_1),
    batch = "1234",
    algo = "test_algo",
    f_scores = c(F2 = 0.2, F1 = 0.4),
    database = c(test_database, test_database),
    collection = test_collection
  ))

  # Wrong dataframe input columns
  expect_error(export_scores(
    formatted_data = data.frame(),
    batch = "1234",
    algo = "test_algo",
    f_scores = c(F2 = 0.2, F1 = 0.4),
    database = test_database,
    collection = test_collection
  ))
})

test_that("Scores are well exported with export_scores_to_mongodb function", {
  dbconnection$remove(query = "{}")
  expect_error(
    export_scores_to_mongodb(
      formatted_data = cbind(siret_df, score = scores_1),
      batch = "1901",
      database = test_database,
      collection = test_collection,
      algo = "test_algo",
      f_scores = c(F2 = 0.2, F1 = 0.4)
    ),
    NA
  )
  expect_equal(dbconnection$count(), 30)
  expect_true("timestamp" %in% names(
    dbconnection$find('{"siret":"1", "periode":"2019-01-01"}')
  ))
  expect_equal(
    dbconnection$find('{"siret":"1", "periode":"2019-01-01"}') %>%
      dplyr::select(-timestamp),
    structure(list(
      siret = "1", periode = "2019-01-01", score = 0.265508663, score_diff = 1,
      algo = "test_algo", batch = "1901", alert = "Alerte seuil F2"), class = "data.frame", row.names = 1L)
  )

  export_scores_to_mongodb(
    formatted_data = cbind(siret_df, score = scores_2),
    batch = "1902",
    database = test_database,
    collection = test_collection,
    algo = "test_algo",
    f_scores = c(F2 = 0.2, F1 = 0.4)
  )
  export_scores_to_mongodb(
    formatted_data = cbind(siret_df, score = scores_3),
    batch = "1903",
    database = test_database,
    collection = test_collection,
    algo = "test_algo",
    f_scores = c(F2 = 0.2, F1 = 0.4)
  )
  expect_equal(dbconnection$count(), 90)
})

# ------------------------------------------------------------------------------
# -------- Import back to R ----------------------------------------------------
# ------------------------------------------------------------------------------
test_that("Input checking for get_scores", {
  # wrong inputs
  expect_error(get_scores(
    database = 123,
    collection = test_collection,
    method = "historic"
  ))

  expect_error(get_scores(
    database = test_database,
    collection = c(test_database, test_database),
    method = "historic"
  ))

  # Method mispelled
  expect_error(get_scores(
    database = test_database,
    collection = test_collection,
    method = "historci"
  ))

  # Database / collection does not exist
  expect_error(get_scores(
    database = "wrong_name",
    collection = test_collection,
    method = "historic"
  ))

  expect_error(get_scores(
    database = test_database,
    collection = "wrong_name",
    method = "historic"
  ))
})

test_that("Get_scores works with 'historical' method", {
  timestamp <- Sys.time()
  Sys.sleep(1)
  export_scores_to_mongodb(
    formatted_data = cbind(siret_df, score = scores_2),
    batch = "1901",
    database = test_database,
    collection = test_collection,
    f_scores = c(F2 = 0.2, F1 = 0.4),
    algo = "test_algo"
  )
  export_scores_to_mongodb(
    formatted_data = cbind(siret_df, score = scores_3),
    batch = "1903",
    database = test_database,
    collection = test_collection,
    f_scores = c(F2 = 0.2, F1 = 0.4),
    algo = "test_algo"
  )

  res_last <- get_scores(
    database = test_database,
    collection = test_collection,
    algo = "test_algo",
    method = "last",
    sirets = c("1", "2", "3")
  )

  res_first <- get_scores(
    database = test_database,
    collection = test_collection,
    algo = "test_algo",
    method = "first",
    sirets = c("1", "2", "3")
  )

  # Last available value is taken
  expect_gt(res_last$timestamp[1], timestamp)
  expect_gt(res_first$timestamp[1], timestamp)
  # Right batch is taken
  expect_equal(res_last$batch[1], "1903")
  expect_equal(res_first$batch[1], "1901")
})

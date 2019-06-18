context("Check database connection")
# Tester à chaque fois pour un frame spark et pour un dataframe

test_db <- "unittest_signauxfaibles"
test_col <- "Features_for_tests"

test_that("Les requêtes sont bien formées", {
  test_grid <- expand.grid(
    batch = "1812",
    date_inf = c("2014-02-01", NA),
    date_sup = c("2015-03-01", NA),
    min_effectif = c(10, NA),
    fields = c("siren", NA),
    siren = c("0123456789", NA),
    code_ape = c("A", NA),
    subsample = c(10, NULL),
    stringsAsFactors = FALSE
    )

  aux_null <- function(x) ifelse(is.na(x), return(NULL), return(x))
  aux_test_function <- function(batch, date_inf, date_sup, min_effectif,
    fields, siren, code_ape, subsample) {
    req <- factor_request(
      batch = batch,
      date_inf = aux_null(date_inf),
      date_sup = aux_null(date_sup),
      min_effectif = aux_null(min_effectif),
      fields = aux_null(fields),
      siren = aux_null(siren),
      code_ape = aux_null(code_ape),
      subsample = subsample
      )
    expect_true(jsonlite::validate(req),
      info = req
      )
  }

  mapply(
    aux_test_function,
    test_grid$batch,
    test_grid$date_inf,
    test_grid$date_sup,
    test_grid$min_effectif,
    test_grid$fields,
    test_grid$siren,
    test_grid$code_ape,
    test_grid$subsample
    )
})


test_that("Une requête vide renvoie un dataframe vide", {
  empty_query <- connect_to_database(
    test_db,
    test_col,
    min_effectif = 10,
    "1901_interim",
    date_sup = "2001-01-01",
    fields = c("siret", "periode"),
    verbose = FALSE
    )
  expect_true(any(dim(empty_query %>% collect()) == 0))
})

test_that(": Un champs vide présent et complété avec des NA", {
  missing_field <- connect_to_database(
    test_db,
    test_col,
    min_effectif = 10,
    "1901_interim",
    fields = c("siret", "periode", "missing_field"),
    verbose = FALSE
    )

  expect_true(all(is.na(
        missing_field %>% select("missing_field") %>% collect()
        )))
})

test_frame_1 <- connect_to_database(
  test_db,
  test_col,
  batch = "1901_interim",
  siren = c("012345678", "876543210"),
  date_inf = "2014-06-01",
  min_effectif = 10,
  fields = NULL,
  code_ape = NULL,
  subsample = NULL,
  verbose = FALSE
  )

test_frame_2 <- connect_to_database(
  test_db,
  test_col,
  batch = "1901_interim",
  siren = NULL,
  date_inf = NULL,
  min_effectif = NULL,
  fields = c("siret", "siren", "periode", "effectif", "code_naf"),
  code_ape = c("H"),
  subsample = NULL,
  verbose = FALSE
  )

test_that("Les filtres fonctionnent comme espéré", {
  test_summary <- test_frame_1 %>%
    summarize(
      date_min = min(periode),
      effectif_min = min(effectif)
      ) %>%
  collect()

expect_gte(
  as.numeric(as.Date(test_summary$date_min)),
  as.numeric(as.Date("2014-01-01"))
  )
expect_gte(test_summary$effectif_min, 10)

sirens_summary <- test_frame_1 %>%
  select("siren") %>%
  distinct() %>%
  collect()
# test siren
expect_setequal(sirens_summary$siren, c("012345678", "876543210"))

# Test fields and code_ape
expect_setequal(
  tbl_vars(test_frame_2),
  c("siret", "siren", "periode", "effectif", "code_naf")
  )
expect_equal(
  test_frame_2 %>%
    select(code_naf) %>%
    distinct() %>%
    collect() %>%
    .$code_naf,
  "H"
  )
  })

# TODO: wrong database or collection should throw an error, not returning empty dataframe.

test_that("No problem in get_fields", {
  expect_false(all(c(
        "outcome",
        "siret",
        "periode",
        "siren",
        "code_ape"
        ) %in% get_fields(training = TRUE)))
  expect_true(all(c(
        "TargetEncode_code_ape_niveau2",
        "TargetEncode_code_ape_niveau3"
        ) %in% get_fields(training = TRUE, target_encode = 2)))
  expect_false(all(c(
        "TargetEncode_code_ape_niveau2",
        "TargetEncode_code_ape_niveau3"
        ) %in% get_fields(training = FALSE, target_encode = 2)))
  expect_true(all(c(
        "outcome",
        "siret",
        "periode",
        "siren",
        "code_ape"
        ) %in% get_fields(training = FALSE)))
  })

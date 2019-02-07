context("Check database connection")
# Tester à chaque fois pour un frame spark et pour un dataframe

test_that("Les requêtes sont bien formées", {
  test_grid <- expand.grid(
    batch = "1812", date_inf = c("2014-02-01", NA),
    date_sup = c("2015-03-01", NA), min_effectif = c(10, NA),
    fields = c("siren", NA), siren = c("0123456789", NA), code_ape = c("A", NA),
    limit = NULL,
    stringsAsFactors = FALSE)

  aux_null <- function(x) ifelse(is.na(x), return(NULL), return(x))
  aux_test_function <- function(batch, date_inf, date_sup, min_effectif,
                                fields, siren, code_ape, limit){
    req <- factor_request(
      batch = batch,
      date_inf = aux_null(date_inf),
      date_sup = aux_null(date_sup),
      min_effectif = aux_null(min_effectif),
      fields = aux_null(fields),
      siren = aux_null(siren),
      code_ape = aux_null(code_ape),
      limit = limit
    )
    expect_true( req %>% jsonlite::validate(),
                 info = req)
  }

  mapply(aux_test_function,
         test_grid$batch,
         test_grid$date_inf,
         test_grid$date_sup,
         test_grid$min_effectif,
         test_grid$fields,
         test_grid$siren,
         test_grid$code_ape,
         test_grid$limit
  )
})


test_procedure <- function(type){


  test_that(paste0(type, ": Une requête vide renvoie un dataframe vide"), {
    empty_query <- connect_to_database(
      "test_signauxfaibles", "test1",
      "1812", date_sup = "2001-01-01", fields = c("siret", "periode"),
      type = type)
    expect_true(any(dim(empty_query %>% collect()) == 0))
  })

  test_that(paste0(type, ": Un champs vide est présent et complété avec des NA"), {
    missing_field <- connect_to_database(
      "test_signauxfaibles", "test1",
      "1812", fields = c("siret", "periode", "missing_field"),
      type = type)

    expect_true(all(is.na(
      missing_field %>% select("missing_field") %>% collect()
    )))
  })

  test_frame_1 <- connect_to_database(
    "test_signauxfaibles",
    "test1",
    "1812",
    siren = c("006280234", "007080260"),
    date_inf = "2016-01-01",
    min_effectif = 10,
    fields = NULL,
    code_ape = NULL,
    type = type,
    limit = NULL
  )

  test_frame_2 <- connect_to_database(
    "test_signauxfaibles",
    "test1",
    "1812",
    siren = NULL,
    date_inf = NULL,
    min_effectif = NULL,
    fields = c("siret", "siren", "periode", "effectif", "code_naf"),
    code_ape = c("H"),
    type = type,
    limit = NULL
  )

  test_that(paste0(type, ": Les filtres fonctionnent comme espéré"), {

    test_summary <- test_frame_1 %>% summarize(
      date_min = min(periode),
      effectif_min = min(effectif)) %>%
      collect()

    expect_gte(as.numeric(as.Date(test_summary$date_min)),
               as.numeric(as.Date("2016-01-01")))
    expect_gte(test_summary$effectif_min, 10)

    sirens_summary <- test_frame_1 %>%
      select("siren") %>%
      distinct() %>%
      collect()
    # test siren
    expect_setequal(sirens_summary$siren, c("006280234","007080260"))

    #Test fields and code_ape
    expect_setequal(tbl_vars(test_frame_2),
                    c("siret", "siren", "periode", "effectif", "code_naf"))
    expect_equal(test_frame_2 %>% select(code_naf) %>% distinct() %>% collect() %>% .$code_naf,
                 "H")

  })

}

test_procedure("dataframe")
test_procedure("spark")

# TODO: wrong database or collection should throw an error, not returning empty dataframe.

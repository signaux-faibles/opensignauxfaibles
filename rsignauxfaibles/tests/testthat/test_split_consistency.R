context("Check the test/train split function")

# Tester à chaque fois pour un frame spark et pour un dataframe

my_test_frame <- expand.grid(
  siret = unlist(tidyr::unite(expand.grid(
    siren = 100000001:100000500,
    et_id = 20001:20004
  ), col = "siret", sep = "")),
  periode = seq.Date(
    from = as.Date("2014-01-01"),
    to = as.Date("2014-12-01"), "month"
  ),
  stringsAsFactors = FALSE
) %>%
  dplyr::mutate(siren = substr(siret, 1, 9)) %>%
  tibble::as.tibble()


sc <- sparklyr::spark_connect(master = "local[*]")
my_test_frame_spark <- copy_to(sc, my_test_frame, overwrite = TRUE)

test_procedure <- function(frame_to_test, prefix) {
  res <- split_snapshot_rdm_month(
    frame_to_test,
    frac_train = 0.60,
    frac_val = 0.25
  )

  train <- res[["train"]] %>%
    mutate(siren = substr(siret, 1, 9))
  validation <- res[["validation"]] %>%
    mutate(siren = substr(siret, 1, 9))
  test <- res[["test"]] %>%
    mutate(siren = substr(siret, 1, 9))

  combined <- rbind(train, validation, test)

  add_prefix <- function(text) {
    paste(prefix, text, sep = "_")
  }

  test_that(add_prefix("Le format du resultat est respecte"), {
    expect_false(is_grouped_df(train))
    expect_true(all(tbl_vars(res[["train"]]) == c("siret", "periode")))
  })



  test_that(add_prefix("Les échantillons ont les bonnes proportions"), {
    expect_ratio <- function(sample, total, frac) {
      # sparklyr compatible function
      n_distinct <- function(myframe) {
        myframe %>%
          distinct() %>%
          count() %>%
          collect() %>%
          unlist() %>%
          as.vector()
      }
      expect_lt(abs(
        n_distinct(sample %>% select(siren)) /
          n_distinct(total %>% select(siren)) - frac
      ), 0.1)
    }

    expect_ratio(train, frame_to_test, 0.60)
    expect_ratio(validation, frame_to_test, 0.25)
    expect_ratio(test, frame_to_test, 0.15)
  })

  test_that(add_prefix("Il n'y a pas de fuite de données entre échantillons"), {
    if (train %>% inherits("tbl_spark")) {
      expect_equal(sdf_nrow(train %>% semi_join(validation, by = "siren")), 0)
      expect_equal(sdf_nrow(train %>% semi_join(test, by = "siren")), 0)
      expect_equal(sdf_nrow(test %>% semi_join(validation, by = "siren")), 0)
    } else {
      expect_equal(nrow(train %>% semi_join(validation, by = "siren")), 0)
      expect_equal(nrow(train %>% semi_join(test, by = "siren")), 0)
      expect_equal(nrow(test %>% semi_join(validation, by = "siren")), 0)
    }
  })

  test_that(add_prefix("Les échantillons ne dépendent pas de l'ordre des données d'entrée
          et restent identiques d'une fois sur l'autre"), {
    folder <- rprojroot::find_rstudio_root_file("tests", ".known_outputs")

    if (!dir.exists(folder)) skip("known values only on local repository")

    if (frame_to_test %>% inherits("tbl_spark")) {
      expect_known_output(frame_to_test %>% mutate(.aux = rand()) %>% arrange(.aux) %>% select(-.aux),
        file.path(folder, add_prefix("test_split")),
        update = TRUE
      )
    } else {
      expect_known_output(sample_n(frame_to_test, size = nrow(frame_to_test), replace = FALSE),
        file.path(folder, add_prefix("test_split")),
        update = TRUE
      )
    }
  })

  test_that(add_prefix("Chaque entreprise appartient au moins à un échantillon"), {
    expect_true(all(unique(frame_to_test$siret) %in% combined$siret))
  })
}

test_procedure(my_test_frame, "R_dataframe")
test_procedure(my_test_frame_spark, "spark_dataframe")

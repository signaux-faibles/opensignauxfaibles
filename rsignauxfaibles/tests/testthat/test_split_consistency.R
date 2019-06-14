context("Check the test/train split function")

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
tibble::as_tibble()


res <- split_snapshot_rdm_month(
  my_test_frame,
  fracs = c(0.60, 0.25, 0.15),
  names = c("train", "validation", "test")
  )

train <- res[["train"]] %>%
  mutate(siren = substr(siret, 1, 9))
validation <- res[["validation"]] %>%
  mutate(siren = substr(siret, 1, 9))
test <- res[["test"]] %>%
  mutate(siren = substr(siret, 1, 9))

combined <- rbind(train, validation, test)

test_that("Le format du resultat est respecte", {
  expect_false(is_grouped_df(train))
  expect_true(all(tbl_vars(res[["train"]]) == c("siret", "periode")))
    })

test_that("Les échantillons ont les bonnes proportions", {
  expect_ratio <- function(sample, total, frac) {
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

  expect_ratio(train, my_test_frame, 0.60)
  expect_ratio(validation, my_test_frame, 0.25)
  expect_ratio(test, my_test_frame, 0.15)
    })


test_that("Il n'y a pas de fuite de données entre échantillons", {
  expect_equal(nrow(train %>% semi_join(validation, by = "siren")), 0)
  expect_equal(nrow(train %>% semi_join(test, by = "siren")), 0)
  expect_equal(nrow(test %>% semi_join(validation, by = "siren")), 0)
    })

test_that(
  "Les échantillons ne dépendent pas de l'ordre des données d'entrée et
  restent identiques d'une fois sur l'autre", {
    folder <- rprojroot::find_rstudio_root_file("tests", ".known_outputs")

    if (!dir.exists(folder)) skip("known values only on local repository")

    scrambled_test_frame  <- sample_n(
      my_test_frame,
      size = nrow(my_test_frame),
      replace = FALSE
      )
    expect_equal(
      split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.60, 0.25, 0.15),
        names = c("train", "validation", "test")
        ),
      split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.60, 0.25, 0.15),
        names = c("train", "validation", "test")
        )
      )
    expect_known_output(
      split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.60, 0.25, 0.15),
        names = c("train", "validation", "test")
        ),
      file.path(folder, "test_split"),
      print = TRUE,
      update = FALSE
      )
  })

test_that(
  "Chaque entreprise appartient au moins à un échantillon", {
    expect_true(all(unique(my_test_frame$siret) %in% combined$siret))
  }
  )


test_that(
  "split_snapshot_rdm_month gère un nom unique", {
      split <- split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.25, 0.25, 0.25, 0.25),
        names = c("cv")
        )
  actual_names <- names(split)
  expected_names  <- c("cv_1", "cv_2", "cv_3", "cv_4")

  expect_equal(actual_names, expected_names)
  }
  )


test_that(
  "La fusion de deux catégories ne change pas la reproductibilité de la
  troisième", {
  split_1  <- split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.25, 0.25, 0.25, 0.25),
        names = c("no1", "no2", "no3", "test")
        )
  split_2  <- split_snapshot_rdm_month(
        my_test_frame,
        fracs = c(0.75, 0.25),
        names = c("no1", "test")
        )

  expect_equal(split_1[["test"]], split_2[["test"]])
}
)

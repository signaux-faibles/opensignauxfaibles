context("Test des utilitaires pour charger ou enregistrer les objets H2O")


connect_to_h2o()

h2o::h2o.no_progress()
h2o_test_frame <- h2o::as.h2o(
  datasets::mtcars %>%
    dplyr::mutate(
      vehicle_name = rownames(mtcars),
      outcome = round(runif(
        nrow(datasets::mtcars)
      ))
    )
)

test_frame_list_1 <- list(h2o_test_frame, h2o_test_frame)
test_frame_list_2 <- list(h2o_test_frame)
test_model_1 <- h2o::h2o.glm(
  "model_1",
  x = names(h2o_test_frame)[1:11],
  y = "outcome",
  family = "binomial",
  training_frame = h2o_test_frame
)

test_model_2 <- h2o::h2o.glm(
  "model_2",
  x = names(h2o_test_frame)[1:11],
  y = "outcome",
  family = "binomial",
  training_frame = h2o_test_frame
  )

test_that("L'enregistrement et le chargement de fichiers
  de classe model fonctionne", {
  full_path <- file.path(
    rprojroot::find_rstudio_root_file(),
    "tests",
    "testthat"
    )
  rel_path <- file.path(".", "tests", "testthat")
  obj_path_1  <- save_h2o_object(test_model_1, "test_model_obj", rel_path)
  obj_path_2  <- save_h2o_object(test_model_2, "test_model_obj", rel_path)
  obj <- load_h2o_object(
    "test_model_obj",
    "model",
    rel_path,
    last = TRUE
  )
  expect_equal(obj@model_id, "model_2")
  obj <- load_h2o_object(
    last = FALSE,
    file_name = paste0(Sys.Date(), "_v1_test_model_obj.model"),
    relative_path = rel_path
    )
  expect_equal(obj@model_id, "model_1")
  file.remove(c(obj_path_1, obj_path_2))
})

test_that("Un mauvais chemin relatif, un fichier inexistant,
  déclenchent une erreur", {
  expect_error(save_h2o_object(
    test_model_1,
    "test_model_obj",
    file.path(".", "tests", "wrong_directory")
  ))
  expect_error(load_h2o_object(
    "wrong_filename",
    "model",
    file.path(".", "tests", "testthat"),
    last = TRUE
  ))
})

test_that("L'enregistrement et le chargement de fichiers
  de classe temap fonctionne", {
  full_path <- file.path(
    rprojroot::find_rstudio_root_file(),
    "tests"
    )
  rel_path  <- file.path(".", "tests", "testthat")
  obj_path_1 <- save_h2o_object(test_frame_list_1, "test_temap_obj", rel_path)
  obj_path_2 <- save_h2o_object(test_frame_list_2, "test_temap_obj", rel_path)
  obj <- load_h2o_object(
    "test_temap_obj",
    "temap",
    file.path(".", "tests", "testthat"),
    last = TRUE
  )
  expect_length(obj, 1)
  obj <- load_h2o_object(
    last = FALSE,
    file_name = paste0(Sys.Date(), "_v1_test_temap_obj.temap"),
    relative_path = rel_path
    )
  expect_length(obj, 2)
  file.remove(c(obj_path_1, obj_path_2))
})

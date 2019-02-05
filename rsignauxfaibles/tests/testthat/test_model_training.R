context("Check H2O model training")

connect_to_h2o()

h2o_test_frame <-
  h2o::as.h2o(datasets::mtcars %>% dplyr::mutate(
    vehicle_name = rownames(mtcars),
    outcome = round(runif(
      nrow(datasets::mtcars)
    ))))

test_that("Typing correction works as expected", {
  h2o_typed_frame <- set_h2o_types(h2o_test_frame)
  expect_equal(h2o::h2o.getTypes(h2o_typed_frame)[[12]], "enum")
  expect_equal(h2o::h2o.getTypes(h2o_typed_frame)[[13]], "enum")
})

aux <- function(my_frame){
  prepare_frame_light_gradient_boosting(my_frame, save_results = FALSE)
}

test_that("convert_to_h2o does not throw an error on NullType conversion",
          {
            # TODO
          })


test_that("Data preparation runs without error", {
  aux_frame <-  h2o_test_frame %>%
    as.data.frame() %>%
    mutate(code_naf = "A", code_ape_niveau2 = "32", code_ape_niveau3 = "123")

  expect_error(aux(aux_frame),
               NA)
  sc <- connect_to_spark(database = "test_signauxfaibles", collection = "test2")
  expect_error(aux(sparklyr::sdf_copy_to(sc, aux_frame)),
               NA)
})


test_that("Model training runs without error", {

  expect_error(train_light_gradient_boosting(
    h2o_train_data = h2o_test_frame,
    save_result = FALSE,
    x_fields_model = names(h2o_test_frame)[
      names(h2o_test_frame) != "outcome"
      ]
  ), NA)
})



## FIX ME
# Test that convert_to_h2o works as expected

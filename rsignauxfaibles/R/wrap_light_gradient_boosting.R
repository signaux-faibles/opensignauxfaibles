
#' Run the model on data
#'
#' @param database Mongodb database name
#' @param collection Mongodb collection name
#' @param periods
#' @param last_batch
#' @param min_effectif
#' @param retrain_model Should the model be retrained, or last available
#' saved model be used ?
#' @param training_date_inf
#' @param training_date_sup Not included
#' @param type Currently only "dataframe" works. Was used to build spark
#' models.
#' @param export_type Should the results be exported ? Can be "none" or NULL
#' for no export, any or both of "csv" and "mongodb" for csv or mongodb
#' export.
#' @param verbose logical. Give info on what is going on ?
#'
#' @return
#' @export
#'
#' @examples
full_light_gradient_boosting <- function(
  database,
  collection,
  periods,
  last_batch,
  min_effectif = 10,
  retrain_model = FALSE,
  training_date_inf = as.Date("2015-01-01"),
  training_date_sup = as.Date("2017-01-01"),
  type = "dataframe",
  export_type = c("csv", "mongodb"),
  verbose = TRUE) {

  require(logger)
  if (verbose) {
    log_threshold(TRACE)
  } else {
    log_threshold(WARN)
  }

  fields <- get_fields(training = FALSE)
  x_fields_model <- get_fields(training = TRUE)


  if (retrain_model) {
    log_info("Loading training data")
    data_frame <- connect_to_database(
      database,
      collection,
      last_batch,
      siren = NULL,
      date_inf = training_date_inf,
      date_sup = training_date_sup,
      min_effectif = min_effectif,
      fields = fields,
      code_ape = NULL,
      type = type,
      subsample = 200000,
      verbose = verbose
    )

    log_info("Model is prepared, in particular target encoding is computed")
    out <- prepare_frame(
      data_to_prepare = data_frame,
      test_or_train = "train",
      save_or_load_map = TRUE, # save
      outcome = "outcome"
    )


    train_data <- out[["data"]]
    te_map <- out[["te_map"]]

    log_info("Model is retrained")
    model <- train_light_gradient_boosting(
      h2o_train_data = train_data,
      x_fields_model = x_fields_model,
      save_results = TRUE
    )
  } else {
    log_info("Last model is loaded")
    model <- load_h2o_object("lgb", "model", last = TRUE)
    te_map <- load_h2o_object("te_map", "temap", last = TRUE)
  }


  log_info("Model is being applied on last batch ...")
  pred_data <- predict_on_last_batch(
    model,
    database,
    collection,
    te_map,
    last_batch,
    periods,
    min_effectif = min_effectif,
    fields = fields
  )
  log_info("Done.")
  export_fields <- c(
    "siret",
    "periode",
    "raison_sociale",
    "departement",
    "region",
    "prob",
    "diff",
    "connu",
    "date_ccsf",
    "etat_proc_collective",
    "date_proc_collective",
    "interessante_urssaf",
    # "default_urssaf",
    "effectif",
    "libelle_naf",
    "libelle_ape5",
    "code_ape",
    "montant_part_ouvriere",
    "montant_part_patronale",
    "ca",
    "ca_past_1",
    "benefice_ou_perte",
    "benefice_ou_perte_past_1",
    "resultat_expl",
    "resultat_expl_past_1",
    "poids_frng",
    "taux_marge",
    "frais_financier",
    "financier_court_terme",
    "delai_fournisseur",
    "dette_fiscale",
    "apart_heures_consommees",
    "apart_heures_autorisees",
    # "cotisation_moy12m",
    "compte_urssaf",
    "montant_majorations",
    "exercice_bdf",
    "exercice_diane",
    "delai",
    "apparait",
    "disparait"
  )

  # TODO compute F-scores !!
  F_scores <- c(F1 = 0.31, F2 = 0.13)
  # Export
  # export(res, batch = last_batch, database = "test_signauxfaibles",
  # destination = "mongodb")


  if (!is.null(export_type) && export_type != "none") {
    assertthat::assert_that(all(export_type %in% c("csv", "mongodb")))

    log_info("Adding additional fields for export")
    res <- pred_data %>%
      prepare_for_export(
        export_fields = export_fields,
        database = database,
        collection = collection,
        last_batch = last_batch
      )
    log_info("Data is exported to {paste(export_type, collapse = ' and ')}")
    purrr::map(
      .x = export_type, function(x, ...) export_scores(destination = x, ...),
      donnees = res,
      batch = last_batch,
      F_scores = F_scores
    )
  }
  return(TRUE)
}

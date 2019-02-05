predict_model <- function(
    model,
    database,
    collection,
    te_map,
    last_batch,
    actual_period,
    algo,
    min_effectif,
    fields
  ){

  current_data <- connect_to_database(
    database,
    collection,
    last_batch,
    date_inf = actual_period %m-% months(1),
    date_sup = actual_period %m+% months(1),
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields[fields != "outcome"]
    )

  current_data <- convert_to_h2o(current_data)

  current <- h2o_target_encode(
    te_map,
    current_data,
    "test")

  prediction <- h2o::h2o.cbind(current,
                          h2o::h2o.predict(model, current) %>%
                          setNames(list("predicted", ".discard", "prob"))) %>%
      tibble::as.tibble()
  prediction <- prediction %>%
    select(-".discard") %>%
    mutate(
    # H2O bug ??
    periode =  as.Date(structure(periode / 1000,
        class = c("POSIXct", "POSIXt")))
    ) %>%
  select(predicted, prob, siret, periode)

  pred_data <- prediction %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(last_prob = lag(prob)) %>%
    ungroup() %>%
    mutate(diff = prob - last_prob) %>%
    filter(periode == actual_period)

  return(pred_data)
}

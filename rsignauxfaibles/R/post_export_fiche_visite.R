export_fiche_visite <- function(
  sirets,
  database = "test_signauxfaibles",
  collection = "Features",
  batch,
  with_urssaf = FALSE,
  folder = batch){




  for (i in seq_along(sirets)) {

    elementary_info  <- connect_to_database(
      database = database,
      collection = collection,
      batch = batch,
      siren = substr(sirets[i], 1, 9),
      fields = c("siret", "raison_sociale", "periode", "code_ape_niveau3"),
      min_effectif = 0
    ) %>%
      filter(siret == sirets[i]) %>%
      head(n = 1)

    raison_sociale <- elementary_info %>%
      .$raison_sociale

    rmarkdown::render("/home/pierre/Documents/opensignauxfaibles/rsignauxfaibles/R/post_fiche_visite.Rmd",
      params  = list(
        siret = sirets[i],
        batch = batch,
        database = database,
        collection = collection,
        with_urssaf = with_urssaf
      ),
      output_file = paste0("/home/pierre/Documents/opensignauxfaibles/output/Fiches/", folder, "/Fiche_visite_", stringr::str_replace_all(raison_sociale, "[^[:alnum:]]", "_"), ".pdf"),
      clean = TRUE
    )
    }
}
